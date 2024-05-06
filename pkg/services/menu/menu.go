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
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/recipe"
)

// RequestData is the data structure that the API expects to receive.
type RequestData struct {
	Menu   Menu          `json:",omitempty"`
	Pantry []ProductData `json:",omitempty"`
	Format string        `json:",omitempty"`
}

type Service struct {
	db database.DB
}

func OneShot(log logger.Logger, db database.DB, menu Menu, pantry []ProductData) ([]ProductData, error) {
	s := New(db)
	return s.ComputeShoppingList(log, menu, pantry)
}

func New(db database.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Handle(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}

	out, err := io.ReadAll(r.Body)
	if err != nil {
		return httputils.Error(http.StatusBadRequest, "failed to read request")
	}
	r.Body.Close()

	var data RequestData
	if err := json.Unmarshal(out, &data); err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to unmarshal request: %v:\n%s", err, string(out))
	}

	log.Debugf("Received request with %d days and %d items in the pantry: ", len(data.Menu.Days), len(data.Pantry))

	if data.Format == "" {
		data.Format = "table"
	}
	f, err := formatter.New(data.Format)
	if err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to create formatter: %v", err)
	}

	log.Debug("Selected formatter: ", data.Format)

	shoppingList, err := s.ComputeShoppingList(log, data.Menu, data.Pantry)
	if err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to compute shopping list: %v", err)
	}

	log.Debug("Computed shopping list")

	if err := f.PrintHead(w, "Product", "Amount", "Unit Cost"); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not write header to output: %w", err)
	}

	i := 0
	for _, p := range shoppingList {
		if p.Amount == 0 {
			continue
		}

		if err := f.PrintRow(w, map[string]any{
			"Product":   p.Name,
			"Amount":    p.Amount,
			"Unit Cost": formatter.Euro(p.UnitCost),
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

func (s Service) ComputeShoppingList(log logger.Logger, menu Menu, pantry []ProductData) ([]ProductData, error) {
	type recipeAmount struct {
		recipe *recipe.Recipe
		amount float32
	}

	// Compute the amount of each recipe needed
	recipes := make(map[string]recipeAmount)

	for _, day := range menu.Days {
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
	products := make(map[string]float32)
	for _, rec := range recipes {
		for _, i := range rec.recipe.Ingredients {
			_, ok := products[i.Name]
			if !ok {
				products[i.Name] = 0
			}
			products[i.Name] += rec.amount * i.Amount
		}
	}

	// Subtract the amount of products in the pantry
	for _, p := range pantry {
		_, ok := products[p.Name]
		if !ok {
			continue
		}
		products[p.Name] = max(0, products[p.Name]-p.Amount)
	}

	// Assemble the output
	table := make([]ProductData, 0, len(products))
	for name, amount := range products {
		product, ok := s.db.LookupProduct(name)
		if !ok {
			log.Warningf("Product %q is not registered", name)
			continue
		}

		table = append(table, ProductData{
			Name:     product.Name,
			Amount:   amount,
			UnitCost: product.Price,
		})
	}

	// Make output consistent
	slices.SortFunc(table, func(a, b ProductData) int {
		return strings.Compare(a.Name, b.Name)
	})

	return table, nil
}
