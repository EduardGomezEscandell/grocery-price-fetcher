package shoppinglist

import (
	"encoding/json"
	"io"
	"math"
	"net/http"
	"slices"
	"strings"

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

func New(settings Settings, db database.DB) *Service {
	return &Service{
		settings: settings,
		db:       db,
	}
}

func (s Service) Name() string {
	return "shopping-list"
}

func (s Service) Path() string {
	return "/api/shopping-list/{menu}/{pantry}"
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
	if err := httputils.ValidateAccepts(r, "application/json"); err != nil {
		return err
	}

	menu := r.PathValue("menu")
	pantry := r.PathValue("pantry")

	m, ok := s.db.LookupMenu(menu)
	if !ok {
		return httputils.Errorf(http.StatusNotFound, "menu not found")
	}

	p, ok := s.db.LookupPantry(pantry)
	if !ok {
		return httputils.Errorf(http.StatusNotFound, "pantry not found")
	}

	var done []string
	if D, ok := s.db.LookupShoppingList(menu, pantry); ok {
		done = D.Contents
	}

	sl := s.computeShoppingList(log, m, p, done)
	log.Debugf("Responding with shopping list with %d items", len(sl))

	if err := json.NewEncoder(w).Encode(map[string]any{
		"menu":   menu,
		"pantry": pantry,
		"items":  sl,
	}); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not encode response: %v", err)
	}

	return nil
}

func (s *Service) handlePut(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if err := httputils.ValidateContentType(r, "application/json"); err != nil {
		return err
	}

	menu := r.PathValue("menu")
	pantry := r.PathValue("pantry")

	if _, ok := s.db.LookupMenu(menu); !ok {
		return httputils.Errorf(http.StatusNotFound, "menu not found")
	}

	if _, ok := s.db.LookupPantry(pantry); !ok {
		return httputils.Errorf(http.StatusNotFound, "pantry not found")
	}

	out, err := io.ReadAll(r.Body)
	if err != nil {
		return httputils.Error(http.StatusBadRequest, "failed to read request")
	}
	r.Body.Close()

	sl := dbtypes.ShoppingList{
		Menu:   menu,
		Pantry: pantry,
	}

	if err := json.Unmarshal(out, &sl.Contents); err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to unmarshal request: %v:\n%s", err, string(out))
	}

	if err := s.db.SetShoppingList(sl); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to store shopping list: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

func (s *Service) handleDelete(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
	menu := r.PathValue("menu")
	pantry := r.PathValue("pantry")

	if _, ok := s.db.LookupMenu(menu); !ok {
		return httputils.Errorf(http.StatusNotFound, "menu not found")
	}

	if _, ok := s.db.LookupPantry(pantry); !ok {
		return httputils.Errorf(http.StatusNotFound, "pantry not found")
	}

	if err := s.db.DeleteShoppingList(menu, pantry); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to delete shopping list: %v", err)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

type shoppingListItem struct {
	Name  string  `json:"name"`
	Done  bool    `json:"done"`
	Units float32 `json:"units"`
	Packs int     `json:"packs"`
	Cost  float32 `json:"cost"`
}

func (s *Service) computeShoppingList(log logger.Logger, menu dbtypes.Menu, pantry dbtypes.Pantry, done []string) []shoppingListItem {
	need := menuneeds.ComputeNeeds(log, s.db, &menu)
	need.Subtract(&pantry)
	slices.SortFunc(need.Items, func(a, b menuneeds.RecipeItem) int { return strings.Compare(a.Product.Name, b.Product.Name) })
	slices.SortFunc(done, strings.Compare)

	list := make([]shoppingListItem, 0, len(need.Items))
	var i, j int
	for i < len(need.Items) && j < len(done) {
		switch strings.Compare(need.Items[i].Product.Name, done[j]) {
		case 0:
			// This product is needed and marked done in the DB
			list = append(list, newItem(need.Items[i], true))
			i++
			j++
		case -1:
			// This product is needed
			list = append(list, newItem(need.Items[i], false))
			i++
		case 1:
			// This product is marked done but not needed
			j++
		}
	}

	for ; i < len(need.Items); i++ {
		list = append(list, newItem(need.Items[i], false))
	}

	return list
}

func newItem(item menuneeds.RecipeItem, isDone bool) shoppingListItem {
	prod := item.Product
	units := item.Amount
	packs := int(math.Ceil(float64(units / prod.BatchSize)))

	return shoppingListItem{
		Name:  prod.Name,
		Units: units,
		Packs: packs,
		Cost:  float32(packs) * prod.Price,
		Done:  isDone,
	}
}
