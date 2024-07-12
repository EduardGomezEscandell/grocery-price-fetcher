package shoppingneeds

import (
	"cmp"
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"slices"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/auth"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/menuneeds"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
)

type Service struct {
	settings Settings

	db   database.DB
	auth auth.Getter
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

func New(settings Settings, db database.DB, auth auth.Getter) *Service {
	return &Service{
		settings: settings,
		db:       db,
		auth:     auth,
	}
}

func (s Service) Name() string {
	return "shopping-needs"
}

func (s Service) Path() string {
	return "/api/shopping-needs/{menu}"
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
	if err := httputils.ValidateAccepts(r, httputils.MediaTypeJSON); err != nil {
		return err
	}

	user, err := s.auth.GetUserID(r)
	if err != nil {
		return httputils.Errorf(http.StatusUnauthorized, "could not get user id: %v", err)
	}

	menu := r.PathValue("menu")

	m, err := s.db.LookupMenu(user, menu)
	if errors.Is(err, fs.ErrNotExist) {
		return httputils.Errorf(http.StatusNotFound, "menu not found")
	} else if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not get menu: %v", err)
	}

	// Compute needs for the menu
	need := menuneeds.ComputeNeeds(log, s.db, m)
	log.Debugf("Responding menu-needs with %d items", len(need))

	// Build response
	type Item struct {
		ProductID product.ID `json:"product_id"`
		Name      string     `json:"name"`
		Amount    float32    `json:"amount"`
	}

	var items []Item
	for _, i := range need {
		p, err := s.db.LookupProduct(i.ProductID)
		if err != nil {
			log.Warningf("Product %d not found: %v", i.ProductID, err)
			continue
		}

		items = append(items, Item{
			ProductID: i.ProductID,
			Name:      p.Name,
			Amount:    i.Amount,
		})
	}

	slices.SortFunc(items, func(a, b Item) int { return cmp.Compare(a.ProductID, b.ProductID) })

	if err := json.NewEncoder(w).Encode(map[string]any{
		"menu":  m.Name,
		"items": items,
	}); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not encode response: %v", err)
	}

	return nil
}
