package menu

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/formatter"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/types"
)

// ProductData represents a the need for a product and its unit cost.
type ProductData struct {
	Name     string
	Need     float32 `json:",omitempty"`
	Have     float32 `json:",omitempty"`
	UnitCost float32 `json:"unit_cost,omitempty"`
}

type Service struct {
	db database.DB
}

func OneShot(log logger.Logger, db database.DB, menu types.Menu, pantry []ProductData) ([]ProductData, error) {
	s := New(db)
	return s.ComputeShoppingList(log, menu.Days, pantry)
}

func New(db database.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Handle(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGet(log, w, r)
	case http.MethodPost:
		return s.handlePost(log, w, r)
	default:
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}
}

func (s *Service) handleGet(_ logger.Logger, w http.ResponseWriter, _ *http.Request) error {
	if err := json.NewEncoder(w).Encode(s.db.Menus()); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not write menus to output: %w", err)
	}

	return nil
}

func (s *Service) handlePost(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	out, err := io.ReadAll(r.Body)
	if err != nil {
		return httputils.Error(http.StatusBadRequest, "failed to read request")
	}
	r.Body.Close()

	var menu types.Menu
	if err := json.Unmarshal(out, &menu); err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to unmarshal request: %v:\n%s", err, string(out))
	}

	log.Debugf("Received request with %d days", len(menu.Days))

	format := "table"
	switch r.Header.Get("Accept") {
	case "application/json":
		format = "json"
	case "text/csv":
		format = "csv"
	}

	f, err := formatter.New(format)
	if err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to create formatter: %v", err)
	}

	log.Debug("Selected formatter: ", format)

	s.UpdateMenu(log, menu)

	shoppingList, err := s.ComputeShoppingList(log, menu.Days, nil)
	if err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to compute shopping list: %v", err)
	}

	log.Debug("Computed shopping list")

	if err := f.PrintHead(w, "Product", "Need", "Have", "Price"); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not write header to output: %w", err)
	}

	i := 0
	for _, p := range shoppingList {
		if p.Need == 0 {
			continue
		}

		if err := f.PrintRow(w, map[string]any{
			"Product": p.Name,
			"Need":    p.Need,
			"Have":    p.Have,
			"Price":   formatter.Euro(p.UnitCost),
		}); err != nil {
			return httputils.Errorf(http.StatusInternalServerError, "could not write results to output: %w", err)
		}
		i++
	}

	log.Debugf("Responded with %d items", i)

	if err := f.PrintTail(w); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not write footer to output: %w", err)
	}

	return nil
}

func (s Service) UpdateMenu(log logger.Logger, menu types.Menu) {
	if menu.Name == "" {
		menu.Name = "default"
	}

	if err := s.db.SetMenu(menu); err != nil {
		log.Warnf("Could not update menu: %v", err)
	}
}

func (s Service) ComputeShoppingList(log logger.Logger, menu []types.Day, pantry []ProductData) ([]ProductData, error) {
	type recipeAmount struct {
		recipe *types.Recipe
		amount float32
	}

	// Compute the amount of each recipe needed
	recipes := make(map[string]recipeAmount)

	for _, day := range menu {
		for _, meal := range day.Meals {
			for _, dish := range meal.Dishes {
				rpe, ok := s.db.LookupRecipe(dish.Name)
				if !ok {
					log.Warningf("%s: %s: Recipe %q is not registered", day.Name, meal.Name, dish.Name)
					continue
				}
				recipes[rpe.Name] = recipeAmount{
					recipe: &rpe,
					amount: recipes[rpe.Name].amount + dish.Amount,
				}
			}
		}
	}

	// Compute the amount of each product needed
	need := make(map[string]float32)
	for _, rec := range recipes {
		for _, i := range rec.recipe.Ingredients {
			_, ok := need[i.Name]
			if !ok {
				need[i.Name] = 0
			}
			need[i.Name] += rec.amount * i.Amount
		}
	}

	have := make(map[string]float32)
	for _, p := range pantry {
		have[p.Name] = p.Have
	}

	// Assemble the output
	table := make([]ProductData, 0, len(need))
	for name, amount := range need {
		product, ok := s.db.LookupProduct(name)
		if !ok {
			log.Warningf("Product %q is not registered", name)
			continue
		}

		table = append(table, ProductData{
			Name:     product.Name,
			Need:     amount,
			Have:     have[product.Name],
			UnitCost: product.Price,
		})
	}

	// Make output consistent
	slices.SortFunc(table, func(a, b ProductData) int {
		return strings.Compare(a.Name, b.Name)
	})

	return table, nil
}
