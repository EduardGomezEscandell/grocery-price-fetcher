package recipe

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strconv"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/auth"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
)

type Service struct {
	settings Settings
	db       database.DB
	auth     auth.Getter
}

type Settings struct {
	Enable bool
}

func (s Settings) Defaults() Settings {
	return Settings{
		Enable: true,
	}
}

func New(s Settings, db database.DB, auth auth.Getter) Service {
	return Service{
		settings: s,
		db:       db,
		auth:     auth,
	}
}

func (s Service) Name() string {
	return "recipe"
}

func (s Service) Path() string {
	return "/api/recipe/{namespace}/{id}"
}

func (s Service) Enabled() bool {
	return s.settings.Enable
}

func (s Service) Handle(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGet(log, w, r)
	case http.MethodPost:
		return s.handlePost(log, w, r)
	case http.MethodDelete:
		return s.handleDelete(log, w, r)
	default:
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}
}

func (s Service) handleGet(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if err := httputils.ValidateAccepts(r, httputils.MediaTypeJSON); err != nil {
		return err
	}

	_, id, err := parseEndpoint(r)
	if err != nil {
		return err
	}

	user, err := s.auth.GetUserID(r)
	if err != nil {
		return httputils.Errorf(http.StatusUnauthorized, "could not get user: %v", err)
	}

	rec, err := s.db.LookupRecipe(user, id)
	if errors.Is(err, fs.ErrNotExist) {
		return httputils.Errorf(http.StatusNotFound, "recipe %d not found", id)
	} else if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to lookup recipe: %v", err)
	}

	body := recipeMsg{
		ID:          rec.ID,
		Name:        rec.Name,
		Ingredients: make([]ingredient, 0, len(rec.Ingredients)),
	}

	for _, ing := range rec.Ingredients {
		prod, err := s.db.LookupProduct(ing.ProductID)
		if err != nil {
			log.Warningf("Product %d not found: %v", ing.ProductID, err)
			continue
		}

		body.Ingredients = append(body.Ingredients, ingredient{
			ID:        prod.ID,
			Name:      prod.Name,
			Amount:    ing.Amount,
			UnitPrice: prod.Price / prod.BatchSize,
		})
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to write response: %v", err)
	}

	return nil
}

func (s Service) handlePost(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if err := httputils.ValidateAccepts(r, httputils.MediaTypeJSON); err != nil {
		return err
	}

	namespace, urlID, err := parseEndpoint(r)
	if err != nil {
		return err
	}

	user, err := s.auth.GetUserID(r)
	if err != nil {
		return httputils.Errorf(http.StatusUnauthorized, "could not get user: %v", err)
	}

	var body recipeMsg
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to read request: %v", err)
	}

	if body.ID != urlID {
		return httputils.Errorf(http.StatusBadRequest, "recipe ID mismatch: %d != %d", body.ID, urlID)
	}

	dbRecipe := recipe.Recipe{
		User:        user,
		ID:          body.ID,
		Name:        body.Name,
		Ingredients: make([]recipe.Ingredient, 0, len(body.Ingredients)),
	}

	if len(body.Ingredients) > 1000 {
		return httputils.Error(http.StatusBadRequest, "recipe cannot have more than 1000 ingredients")
	}

	if body.Name == "" {
		return httputils.Error(http.StatusBadRequest, "recipe name cannot be empty")
	}

	for _, ing := range body.Ingredients {
		dbRecipe.Ingredients = append(dbRecipe.Ingredients, recipe.Ingredient{
			ProductID: ing.ID,
			Amount:    ing.Amount,
		})
	}

	newID, err := s.db.SetRecipe(dbRecipe)
	if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to save recipe: %v", err)
	}

	if urlID != 0 {
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Location", path.Join("/api/recipe/", namespace, fmt.Sprint(newID)))
	}

	if err := json.NewEncoder(w).Encode(map[string]interface{}{"id": newID}); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to write response: %v", err)
	}

	return nil
}

func (s Service) handleDelete(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
	_, id, err := parseEndpoint(r)
	if err != nil {
		return err
	}

	user, err := s.auth.GetUserID(r)
	if err != nil {
		return httputils.Errorf(http.StatusUnauthorized, "could not get user: %v", err)
	}

	if err := s.db.DeleteRecipe(user, id); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to delete recipe: %v", err)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

type ingredient struct {
	ID        product.ID `json:"id"`
	Name      string     `json:"name"`
	Amount    float32    `json:"amount"`
	UnitPrice float32    `json:"unit_price"`
}

type recipeMsg struct {
	ID          recipe.ID    `json:"id"`
	Name        string       `json:"name"`
	Ingredients []ingredient `json:"ingredients"`
}

func parseEndpoint(r *http.Request) (namespace string, id recipe.ID, err error) {
	n := r.PathValue("namespace")
	if n == "" {
		return "", 0, httputils.Error(http.StatusBadRequest, "missing namespace")
	} else if n != "default" {
		// Only the default namespace is supported for now
		return "", 0, httputils.Errorf(http.StatusNotFound, "namespace %s not found", n)
	}

	sid := r.PathValue("id")
	if sid == "" {
		return "", 0, httputils.Error(http.StatusBadRequest, "missing id")
	}

	idURL, err := strconv.ParseUint(sid, 10, recipe.IDSize)
	if err != nil {
		return "", 0, httputils.Errorf(http.StatusBadRequest, "invalid id: %v", err)
	}

	return n, recipe.ID(idURL), nil
}
