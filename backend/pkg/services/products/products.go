package products

import (
	"encoding/json"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
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
	return "products"
}

func (s Service) Path() string {
	return "/api/products/{namespace}"
}

func (s Service) Enabled() bool {
	return s.settings.Enable
}

func (s Service) Handle(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGet(log, w, r)
	default:
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}
}

func (s Service) handleGet(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if err := httputils.ValidateAccepts(r, "application/json"); err != nil {
		return err
	}

	ns := r.PathValue("namespace")
	if ns == "" {
		return httputils.Error(http.StatusBadRequest, "missing namespace")
	} else if ns != "default" {
		// Only the default namespace is supported for now
		return httputils.Errorf(http.StatusNotFound, "namespace %s not found", ns)
	}

	rec, err := s.db.Products()
	if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to fetch products: %v", err)
	}

	if err := json.NewEncoder(w).Encode(rec); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to write response: %v", err)
	}

	return nil
}
