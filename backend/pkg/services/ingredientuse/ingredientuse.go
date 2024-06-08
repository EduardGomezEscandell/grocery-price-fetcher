package ingredientuse

import (
	"encoding/json"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
)

type Service struct {
	settings Settings

	db database.DB
}

type Settings struct {
	Enable bool
}

func (Settings) Defaults() Settings {
	return Settings{
		Enable: true,
	}
}

func New(settings Settings, db database.DB) *Service {
	return &Service{
		settings: settings,
		db:       db,
	}
}

func (s Service) Enabled() bool {
	return s.settings.Enable
}

type reqBody struct {
	MenuName       string `json:"menu_name"`
	IngredientName string `json:"ingredient_name"`
}

type respBodyItem struct {
	Day    string  `json:"day"`
	Meal   string  `json:"meal"`
	Dish   string  `json:"dish"`
	Amount float32 `json:"amount"`
}

func (s Service) Handle(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}

	var b reqBody
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to decode request body: %v", err)
	}

	if b.MenuName == "" {
		return httputils.Errorf(http.StatusBadRequest, "menu_name is required")
	}

	if b.IngredientName == "" {
		return httputils.Errorf(http.StatusBadRequest, "ingredient is required")
	}

	resp, err := s.compute(b)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not write response: %v", err)
	}

	return nil
}

func (s *Service) compute(b reqBody) ([]respBodyItem, error) {
	menu, ok := s.db.LookupMenu(b.MenuName)
	if !ok {
		return nil, httputils.Errorf(http.StatusNotFound, "menu %q not found", b.MenuName)
	}

	resp := make([]respBodyItem, 0)
	cached := database.NewCachedLookup(s.db.LookupRecipe)

	for _, day := range menu.Days {
		for _, meal := range day.Meals {
			for _, dish := range meal.Dishes {
				recipe, ok := cached.Lookup(dish.Name)
				if !ok {
					continue
				}

				for _, ingredient := range recipe.Ingredients {
					if ingredient.Name == b.IngredientName {
						resp = append(resp, respBodyItem{
							Day:    day.Name,
							Meal:   meal.Name,
							Dish:   dish.Name,
							Amount: ingredient.Amount * dish.Amount,
						})
					}
				}
			}
		}
	}

	return resp, nil
}
