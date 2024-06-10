package shoppingneeds

import (
	"encoding/json"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/menuneeds"
)

type Service struct {
	settings Settings

	db database.DB
}

type Settings struct {
	Enable bool
}

func (Settings) Defaults() Settings {
	return Settings{
		Enable: true,
	}
}

var Version = "dev"

func New(settings Settings, db database.DB) *Service {
	return &Service{
		settings: settings,
		db:       db,
	}
}

func (s Service) Path() string {
	return "shopping-needs/{menu}"
}

func (s Service) Enabled() bool {
	return s.settings.Enable
}

func (s *Service) Handle(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGet(log, w, r)
	default:
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}
}

func (s *Service) handleGet(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Accept") != "application/json" {
		return httputils.Errorf(http.StatusBadRequest, "unsupported format: %s", r.Header.Get("Accept"))
	}

	menu := r.PathValue("menu")

	m, ok := s.db.LookupMenu(menu)
	if !ok {
		return httputils.Errorf(http.StatusNotFound, "menu not found")
	}

	// Compute needs for the menu
	need := menuneeds.ComputeNeeds(log, s.db, &m)
	log.Debugf("Responding menu-needs with %d items", len(need.Items))

	// Build response
	var items []dbtypes.Ingredient
	for _, i := range need.Items {
		items = append(items, dbtypes.Ingredient{
			Name:   i.Product.Name,
			Amount: i.Amount,
		})
	}

	if err := json.NewEncoder(w).Encode(map[string]any{
		"menu":  need.Menu.Name,
		"items": items,
	}); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not encode response: %v", err)
	}

	return nil
}
