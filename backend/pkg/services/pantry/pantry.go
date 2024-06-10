package pantry

import (
	"encoding/json"
	"io"
	"net/http"

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

	pantries, err := s.db.Pantries()
	if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not get pantries: %w", err)
	}

	if err := json.NewEncoder(w).Encode(pantries); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not write menus to output: %w", err)
	}

	return nil
}

func (s *Service) handlePost(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
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

	pantry := dbtypes.Pantry{
		Name: "default",
	}

	if err := json.Unmarshal(out, &pantry); err != nil {
		return httputils.Errorf(http.StatusBadRequest, "could not unmarshal pantry: %w", err)
	}

	if err := s.db.SetPantry(pantry); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not set pantry: %w", err)
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}
