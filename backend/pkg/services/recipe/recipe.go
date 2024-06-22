package recipe

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
)

type Service struct {
	settings Settings
	db       database.DB
}

type Settings struct {
	Enable bool
}

func (s Settings) Defaults() Settings {
	return Settings{
		Enable: true,
	}
}

func New(s Settings, db database.DB) Service {
	return Service{
		settings: s,
		db:       db,
	}
}

func (s Service) Name() string {
	return "recipe"
}

func (s Service) Path() string {
	return "/api/recipe/{namespace}/{name}"
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

func (s Service) handleGet(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if err := httputils.ValidateAccepts(r, httputils.MediaTypeJSON); err != nil {
		return err
	}

	ns := r.PathValue("namespace")
	if ns == "" {
		return httputils.Error(http.StatusBadRequest, "missing namespace")
	} else if ns != "default" {
		// Only the default namespace is supported for now
		return httputils.Errorf(http.StatusNotFound, "namespace %s not found", ns)
	}

	name := r.PathValue("name")
	if name == "" {
		return httputils.Error(http.StatusBadRequest, "missing name")
	}

	rec, ok := s.db.LookupRecipe(name)
	if !ok {
		return httputils.Errorf(http.StatusNotFound, "recipe %s not found", name)
	}

	body := recipe{
		Name:        rec.Name,
		Ingredients: make([]ingredient, 0, len(rec.Ingredients)),
	}

	for _, ing := range rec.Ingredients {
		prod, ok := s.db.LookupProduct(ing.Name)
		if !ok {
			continue
		}

		body.Ingredients = append(body.Ingredients, ingredient{
			Name:      ing.Name,
			Amount:    ing.Amount,
			UnitPrice: prod.Price / prod.BatchSize,
		})
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to write response: %v", err)
	}

	return nil
}

func (s Service) handlePost(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if err := httputils.ValidateAccepts(r, httputils.MediaTypeJSON); err != nil {
		return err
	}

	ns := r.PathValue("namespace")
	if ns == "" {
		return httputils.Error(http.StatusBadRequest, "missing namespace")
	} else if ns != "default" {
		// Only the default namespace is supported for now
		return httputils.Errorf(http.StatusNotFound, "namespace %s not found", ns)
	}

	name := r.PathValue("name")
	if name == "" {
		return httputils.Error(http.StatusBadRequest, "missing name")
	}

	var body recipe
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to read request: %v", err)
	}

	dbRecipe := dbtypes.Recipe{
		Name:        body.Name,
		Ingredients: make([]dbtypes.Ingredient, 0, len(body.Ingredients)),
	}

	if len(body.Ingredients) > 1000 {
		return httputils.Error(http.StatusBadRequest, "recipe cannot have more than 1000 ingredients")
	}

	if len(body.Name) == 0 {
		return httputils.Error(http.StatusBadRequest, "recipe name cannot be empty")
	}

	for _, ing := range body.Ingredients {
		dbRecipe.Ingredients = append(dbRecipe.Ingredients, dbtypes.Ingredient{
			Name:   ing.Name,
			Amount: ing.Amount,
		})
	}

	// Simple case: Recipe is being edited
	if name == dbRecipe.Name {
		if err := s.db.SetRecipe(dbRecipe); err != nil {
			return httputils.Errorf(http.StatusInternalServerError, "failed to save recipe: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
		return nil
	}

	// Recipe is being renamed (and possibly edited)
	if _, ok := s.db.LookupRecipe(body.Name); ok {
		// No accidental overwrites
		return httputils.Errorf(http.StatusConflict, "recipe %s already exists", body.Name)
	}

	// Write new recipe, then delete old one
	if err := s.db.SetRecipe(dbRecipe); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to save recipe: %v", err)
	}

	if err := s.db.DeleteRecipe(name); err != nil {
		log.Errorf("failed to delete old recipe during re-naming from %s to %s: %v", name, body.Name, err)
	}

	w.Header().Set("Location", path.Join("/api/recipe/%s/%s", ns, dbRecipe.Name))
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (s Service) handleDelete(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
	ns := r.PathValue("namespace")
	if ns == "" {
		return httputils.Error(http.StatusBadRequest, "missing namespace")
	} else if ns != "default" {
		// Only the default namespace is supported for now
		return httputils.Errorf(http.StatusNotFound, "namespace %s not found", ns)
	}

	name := r.PathValue("name")
	if name == "" {
		return httputils.Error(http.StatusBadRequest, "missing name")
	}

	if err := s.db.DeleteRecipe(name); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to delete recipe: %v", err)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

type ingredient struct {
	Name      string  `json:"name"`
	Amount    float32 `json:"amount"`
	UnitPrice float32 `json:"unit_price"`
}

type recipe struct {
	Name        string       `json:"name"`
	Ingredients []ingredient `json:"ingredients"`
}
