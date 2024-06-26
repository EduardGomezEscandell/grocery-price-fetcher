package pantry

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
)

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

func (s Service) Name() string {
	return "pantry"
}

func (s Service) Path() string {
	return "/api/pantry/{pantry}"
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
	case http.MethodDelete:
		return s.handleDelete(log, w, r)
	default:
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}
}

func (s *Service) handleGet(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if err := httputils.ValidateAccepts(r, httputils.MediaTypeJSON); err != nil {
		return err
	}

	p := r.PathValue("pantry")
	if p == "" {
		return httputils.Error(http.StatusBadRequest, "missing pantry")
	}

	pantry, ok := s.db.LookupPantry(p)
	if !ok {
		return httputils.Error(http.StatusNotFound, "pantry not found")
	}

	type Item struct {
		recipe.Ingredient
		Name string `json:"name"`
	}

	items := make([]Item, 0, len(pantry.Contents))
	for _, ing := range pantry.Contents {
		p, err := s.db.LookupProduct(ing.ProductID)
		if err != nil {
			log.Warnf("Product %d not found: %v", ing.ProductID, err)
			continue
		}

		items = append(items, Item{
			Ingredient: ing,
			Name:       p.Name,
		})
	}

	if err := json.NewEncoder(w).Encode(map[string]any{
		"name":     pantry.Name,
		"contents": items,
	}); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not write menus to output: %w", err)
	}

	return nil
}

func (s *Service) handlePut(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if err := httputils.ValidateContentType(r, httputils.MediaTypeJSON); err != nil {
		return err
	}

	out, err := io.ReadAll(r.Body)
	if err != nil {
		return httputils.Error(http.StatusBadRequest, "failed to read request")
	}
	r.Body.Close()

	p := r.PathValue("pantry")
	if p == "" {
		return httputils.Error(http.StatusBadRequest, "missing pantry")
	}

	pantry := dbtypes.Pantry{
		Name: p,
	}

	if err := json.Unmarshal(out, &pantry); err != nil {
		return httputils.Errorf(http.StatusBadRequest, "could not unmarshal pantry: %w", err)
	}

	log.Debugf("Received pantry with %d items", len(pantry.Contents))

	if err := s.db.SetPantry(pantry); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not set pantry: %w", err)
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

func (s *Service) handleDelete(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
	p := r.PathValue("pantry")
	if p == "" {
		return httputils.Error(http.StatusBadRequest, "missing pantry")
	}

	if err := s.db.DeletePantry(p); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not delete pantry: %w", err)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
