package menu

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/auth"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
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
	auth     auth.Getter
}

type Settings struct {
	Enable bool
}

func (Settings) Defaults() Settings {
	return Settings{
		Enable: true,
	}
}

func New(s Settings, db database.DB, auth auth.Getter) *Service {
	if !s.Enable {
		return nil
	}

	return &Service{
		settings: s,
		db:       db,
		auth:     auth,
	}
}

func (s Service) Name() string {
	return "menu"
}

func (s Service) Path() string {
	return "/api/menu/{menu}"
}

func (s Service) Enabled() bool {
	return s.settings.Enable
}

func (s *Service) Handle(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGet(log, w, r)
	case http.MethodPut:
		return s.handlePut(log, w, r)
	default:
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}
}

func (s *Service) handleGet(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if err := httputils.ValidateAccepts(r, httputils.MediaTypeJSON); err != nil {
		return err
	}

	user, err := s.auth.GetUserID(r)
	if err != nil {
		return httputils.Errorf(http.StatusUnauthorized, "could not get user: %v", err)
	}

	m := r.PathValue("menu")
	if m == "" {
		return httputils.Error(http.StatusBadRequest, "missing menu")
	}

	menu, err := s.db.LookupMenu(user, m)
	if errors.Is(err, fs.ErrNotExist) {
		return httputils.Errorf(http.StatusNotFound, "menu %s not found", m)
	} else if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not get menu: %v", err)
	}

	if err := s.writeMenu(w, menu); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not write menus to output: %w", err)
	}

	return nil
}

func (s *Service) handlePut(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if err := httputils.ValidateContentType(r, httputils.MediaTypeJSON); err != nil {
		return err
	}

	user, err := s.auth.GetUserID(r)
	if err != nil {
		return httputils.Errorf(http.StatusUnauthorized, "could not get user: %v", err)
	}

	name := r.PathValue("menu")
	if name == "" {
		return httputils.Error(http.StatusBadRequest, "missing menu")
	}

	out, err := io.ReadAll(r.Body)
	if err != nil {
		return httputils.Error(http.StatusBadRequest, "failed to read request")
	}
	r.Body.Close()

	menu := dbtypes.Menu{
		Name: name,
	}

	if err := json.Unmarshal(out, &menu); err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to unmarshal request: %v:\n%s", err, string(out))
	}

	menu.User = user

	log.Debugf("Received request with %d days", len(menu.Days))

	if err := s.UpdateMenu(log, menu); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to update menu: %v", err)
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

func (s Service) UpdateMenu(log logger.Logger, menu dbtypes.Menu) error {
	if menu.Name == "" {
		menu.Name = "default"
	}

	if err := s.db.SetMenu(menu); err != nil {
		return err
	}

	return nil
}

// writeMenu writes the menu to the output stream, adding the names to the dishes.
func (s *Service) writeMenu(w io.Writer, m dbtypes.Menu) error {
	type msgDish struct {
		RecipeID recipe.ID `json:"recipe_id"`
		Name     string    `json:"name"`
		Amount   float32   `json:"amount"`
	}

	type msgMeal struct {
		Name   string    `json:"name"`
		Dishes []msgDish `json:"dishes"`
	}

	type msgDay struct {
		Name  string    `json:"name"`
		Meals []msgMeal `json:"meals"`
	}

	type msgMenu struct {
		Name string   `json:"name"`
		Days []msgDay `json:"days"`
	}

	var days []msgDay
	for _, d := range m.Days {
		var meals []msgMeal
		for _, m := range d.Meals {
			var dishes []msgDish
			for _, dish := range m.Dishes {
				recipe, err := s.db.LookupRecipe(dish.ID)
				if err != nil {
					continue
				}

				dishes = append(dishes, msgDish{
					RecipeID: recipe.ID,
					Name:     recipe.Name,
					Amount:   dish.Amount,
				})
			}

			meals = append(meals, msgMeal{
				Name:   m.Name,
				Dishes: dishes,
			})
		}

		days = append(days, msgDay{
			Name:  d.Name,
			Meals: meals,
		})
	}

	return json.NewEncoder(w).Encode(msgMenu{
		Name: m.Name,
		Days: days,
	})
}
