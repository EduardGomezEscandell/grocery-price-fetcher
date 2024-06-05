package menu

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/types"
)

// ProductData represents a the need for a product and its unit cost.
type ProductData struct {
	product.Product

	Need float32 `json:",omitempty"`
	Have float32 `json:",omitempty"`
}

type Service struct {
	settings Settings
	db       database.DB
}

type Settings struct {
	Enable bool
}

func (Settings) Defaults() Settings {
	return Settings{
		Enable: true,
	}
}

func New(s Settings, db database.DB) *Service {
	if !s.Enable {
		return nil
	}

	return &Service{
		settings: s,
		db:       db,
	}
}

func (s Service) Enabled() bool {
	return s.settings.Enable
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

func (s *Service) handleGet(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Accept") != "application/json" {
		return httputils.Errorf(http.StatusBadRequest, "unsupported format: %s", r.Header.Get("Accept"))
	}

	menus, err := s.db.Menus()
	if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not get menus: %v", err)
	}

	if err := json.NewEncoder(w).Encode(menus); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not write menus to output: %w", err)
	}

	return nil
}

func (s *Service) handlePost(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return httputils.Errorf(http.StatusBadRequest, "unsupported content type: %s", r.Header.Get("Content-Type"))
	}

	if r.Header.Get("Accept") != "application/json" {
		return httputils.Errorf(http.StatusBadRequest, "unsupported format: %s", r.Header.Get("Accept"))
	}

	out, err := io.ReadAll(r.Body)
	if err != nil {
		return httputils.Error(http.StatusBadRequest, "failed to read request")
	}
	r.Body.Close()

	menu := types.Menu{
		Name: "default",
	}

	if err := json.Unmarshal(out, &menu); err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to unmarshal request: %v:\n%s", err, string(out))
	}

	log.Debugf("Received request with %d days", len(menu.Days))

	if err := s.UpdateMenu(log, menu); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to update menu: %v", err)
	}

	w.WriteHeader(http.StatusCreated)

	shoppingList, err := s.ComputeShoppingList(log, menu.Days, nil)
	if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to compute shopping list: %v", err)
	}

	type responseItem struct {
		Product   string  `json:"product"`
		Need      float32 `json:"need"`
		Have      float32 `json:"have"`
		BatchSize float32 `json:"batch_size"`
		Price     float32 `json:"price"`
	}

	response := make([]responseItem, 0)
	for _, p := range shoppingList {
		if p.Need == 0 {
			continue
		}

		response = append(response, responseItem{
			Product:   p.Name,
			Need:      p.Need,
			Have:      p.Have,
			BatchSize: p.BatchSize,
			Price:     p.Price,
		})
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not write menus to output: %w", err)
	}

	log.Debugf("Responded with %d items", len(response))
	return nil
}

func (s Service) UpdateMenu(log logger.Logger, menu types.Menu) error {
	if menu.Name == "" {
		menu.Name = "default"
	}

	if err := s.db.SetMenu(menu); err != nil {
		return err
	}

	return nil
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
			Product: product,
			Need:    amount,
			Have:    have[product.Name],
		})
	}

	// Make output consistent
	slices.SortFunc(table, func(a, b ProductData) int {
		return strings.Compare(a.Name, b.Name)
	})

	return table, nil
}
