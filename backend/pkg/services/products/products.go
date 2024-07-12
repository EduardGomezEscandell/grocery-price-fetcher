package products

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strconv"

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
	return "/api/products/{id}"
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

	nm := r.PathValue("id")
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

	id, err := strconv.ParseUint(nm, 10, product.IDSize)
	if err != nil {
		return httputils.Error(http.StatusBadRequest, "could not parse product ID")
	}

	// Return a single product
	p, err := s.db.LookupProduct(product.ID(id))
	if errors.Is(err, fs.ErrNotExist) {
		return httputils.Errorf(http.StatusNotFound, "product %s not found", nm)
	} else if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to fetch product: %v", err)
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to write response: %v", err)
	}

	return nil
}

func (s Service) handlePost(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if err := httputils.ValidateContentType(r, httputils.MediaTypeJSON); err != nil {
		return err
	}

	urlID, err := parseEndpoint(r)
	if err != nil {
		return err
	}

	var body product.Product
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to decode request: %v", err)
	}

	if body.ID != urlID {
		return httputils.Errorf(http.StatusBadRequest, "product ID in URL does not match product ID in body")
	}

	newID, err := s.db.SetProduct(body)
	if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to set product: %v", err)
	}

	if urlID != 0 {
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Location", path.Join("/api/recipe/", fmt.Sprint(newID)))
	}

	if err := json.NewEncoder(w).Encode(map[string]interface{}{"id": newID}); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to write response: %v", err)
	}

	return nil
}

func (s Service) handleDelete(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
	id, err := parseEndpoint(r)
	if err != nil {
		return err
	}

	if err := s.db.DeleteProduct(id); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to delete product: %v", err)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func parseEndpoint(r *http.Request) (id product.ID, err error) {
	sid := r.PathValue("id")
	if sid == "" {
		return 0, httputils.Error(http.StatusBadRequest, "missing id")
	}

	idURL, err := strconv.ParseUint(sid, 10, product.IDSize)
	if err != nil {
		return 0, httputils.Errorf(http.StatusBadRequest, "invalid id: %v", err)
	}

	return product.ID(idURL), nil
}
