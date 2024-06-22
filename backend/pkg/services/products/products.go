package products

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
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
	return "/api/products/{namespace}/{name}"
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

	nm := r.PathValue("name")
	if nm == "*" {
		// Return all products
		rec, err := s.db.Products()
		if err != nil {
			return httputils.Errorf(http.StatusInternalServerError, "failed to fetch products: %v", err)
		}

		if err := json.NewEncoder(w).Encode(rec); err != nil {
			return httputils.Errorf(http.StatusInternalServerError, "failed to write response: %v", err)
		}

		return nil
	}

	// Return a single product
	p, ok := s.db.LookupProduct(nm)
	if !ok {
		return httputils.Errorf(http.StatusNotFound, "product %s not found", nm)
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to write response: %v", err)
	}

	return nil
}

func (s Service) handlePost(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if err := httputils.ValidateContentType(r, httputils.MediaTypeJSON); err != nil {
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
		return httputils.Error(http.StatusBadRequest, "missing product name")
	}

	var body product.Product
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to decode request: %v", err)
	}

	// Easy case: just update the product
	if body.Name == name {
		if err := s.db.SetProduct(body); err != nil {
			return httputils.Errorf(http.StatusInternalServerError, "failed to update product: %v", err)
		}

		return nil
	}

	// Hard case: update the product name
	if _, ok := s.db.LookupProduct(body.Name); ok {
		return httputils.Errorf(http.StatusConflict, "product %s already exists", body.Name)
	}

	if err := s.db.SetProduct(body); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to update product: %v", err)
	}
	w.WriteHeader(http.StatusCreated)

	if err := s.db.DeleteProduct(name); err != nil {
		log.Errorf("failed to delete old product during re-naming from %s to %s: %v", name, body.Name, err)
	}

	w.Header().Set("Location", path.Join("/api/products/%s/%s", ns, body.Name))
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
		return httputils.Error(http.StatusBadRequest, "missing product name")
	}

	if err := s.db.DeleteProduct(name); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to delete product: %v", err)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
