package menu

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
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
	return "menu"
}

func (s Service) Path() string {
	return "/api/menu"
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
	if r.Header.Get("Accept") != "application/json" {
		return httputils.Errorf(http.StatusBadRequest, "unsupported format: %s", r.Header.Get("Accept"))
	}

	menus, err := s.db.Menus()
	if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not get menus: %v", err)
	}

	if err := json.NewEncoder(w).Encode(menus); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not write menus to output: %w", err)
	}

	return nil
}

func (s *Service) handlePut(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return httputils.Errorf(http.StatusBadRequest, "unsupported content type: %s", r.Header.Get("Content-Type"))
	}

	out, err := io.ReadAll(r.Body)
	if err != nil {
		return httputils.Error(http.StatusBadRequest, "failed to read request")
	}
	r.Body.Close()

	menu := dbtypes.Menu{
		Name: "default",
	}

	if err := json.Unmarshal(out, &menu); err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to unmarshal request: %v:\n%s", err, string(out))
	}

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
