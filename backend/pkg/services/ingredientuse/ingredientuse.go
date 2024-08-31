package ingredientuse

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"strconv"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/auth"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/utils"
)

type Service struct {
	settings Settings

	db   database.DB
	auth *auth.Manager
}

type Settings struct {
	Enable bool
}

func (Settings) Defaults() Settings {
	return Settings{
		Enable: true,
	}
}

func New(settings Settings, db database.DB, auth *auth.Manager) *Service {
	return &Service{
		settings: settings,
		db:       db,
		auth:     auth,
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

	if err := httputils.ValidateAccepts(r, httputils.MediaTypeJSON); err != nil {
		return err
	}

	user, err := s.auth.GetUserID(r)
	if err != nil {
		return httputils.Errorf(http.StatusUnauthorized, "could not get user: %v", err)
	}

	menu := r.PathValue("menu")
	ingredientRaw := r.PathValue("ingredient")

	tmp, err := strconv.ParseUint(ingredientRaw, 10, product.IDSize)
	if err != nil {
		return httputils.Errorf(http.StatusBadRequest, "invalid ingredient ID %q: %v", ingredientRaw, err)
	}

	ingredient, err := utils.SafeIntConvert[product.ID](tmp)
	if err != nil {
		return httputils.Errorf(http.StatusBadRequest, "invalid ingredient ID %q: %v", ingredientRaw, err)
	}

	resp, err := s.compute(user, menu, ingredient)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not write response: %v", err)
	}

	return nil
}

func (s *Service) compute(user, menuName string, ingredientID product.ID) ([]respBodyItem, error) {
	menu, err := s.db.LookupMenu(user, menuName)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, httputils.Errorf(http.StatusNotFound, "menu %q not found", menuName)
	} else if err != nil {
		return nil, httputils.Errorf(http.StatusInternalServerError, "could not get menu %q: %v", menuName, err)
	}

	resp := make([]respBodyItem, 0)
	cached := database.NewCachedUserLookup(user, s.db.LookupRecipe)

	for _, day := range menu.Days {
		for _, meal := range day.Meals {
			for _, dish := range meal.Dishes {
				recipe, err := cached.Lookup(dish.ID)
				if err != nil {
					continue
				}

				for _, ingredient := range recipe.Ingredients {
					if ingredient.ProductID == ingredientID {
						resp = append(resp, respBodyItem{
							Day:    day.Name,
							Meal:   meal.Name,
							Dish:   recipe.Name,
							Amount: ingredient.Amount * dish.Amount,
						})
					}
				}
			}
		}
	}

	return resp, nil
}
