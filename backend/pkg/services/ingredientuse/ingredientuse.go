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

func (s Service) Name() string {
	return "ingredientuse"
}

func (s Service) Path() string {
	return "/api/ingredient-use/{menu}/{ingredient}"
}

func (s Service) Enabled() bool {
	return s.settings.Enable
}

type respBodyItem struct {
	Day    string  `json:"day"`
	Meal   string  `json:"meal"`
	Dish   string  `json:"dish"`
	Amount float32 `json:"amount"`
}

func (s Service) Handle(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}

	if err := httputils.ValidateAccepts(r, "application/json"); err != nil {
		return err
	}

	menu := r.PathValue("menu")
	ingredient := r.PathValue("ingredient")

	resp, err := s.compute(menu, ingredient)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not write response: %v", err)
	}

	return nil
}

func (s *Service) compute(menuName, ingredientName string) ([]respBodyItem, error) {
	menu, ok := s.db.LookupMenu(menuName)
	if !ok {
		return nil, httputils.Errorf(http.StatusNotFound, "menu %q not found", menuName)
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
					if ingredient.Name == ingredientName {
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
