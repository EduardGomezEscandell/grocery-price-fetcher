package shoppinglist

import (
	"cmp"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"math"
	"net/http"
	"slices"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/auth"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/menuneeds"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/utils"
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

func New(settings Settings, db database.DB, auth auth.Getter) *Service {
	return &Service{
		settings: settings,
		db:       db,
		auth:     auth,
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
	if err := httputils.ValidateAccepts(r, httputils.MediaTypeJSON); err != nil {
		return err
	}

	user, err := s.auth.GetUserID(r)
	if err != nil {
		return httputils.Errorf(http.StatusUnauthorized, "failed to get user ID: %v", err)
	}

	menu := r.PathValue("menu")
	pantry := r.PathValue("pantry")

	m, err := s.db.LookupMenu(user, menu)
	if errors.Is(err, fs.ErrNotExist) {
		return httputils.Errorf(http.StatusNotFound, "menu not found")
	} else if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to lookup menu: %v", err)
	}

	p, err := s.db.LookupPantry(user, pantry)
	if errors.Is(err, fs.ErrNotExist) {
		return httputils.Errorf(http.StatusNotFound, "pantry not found")
	} else if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to lookup pantry: %v", err)
	}

	var done []product.ID
	if D, err := s.db.LookupShoppingList(user, menu, pantry); errors.Is(err, fs.ErrNotExist) {
		done = make([]product.ID, 0)
	} else if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to lookup shopping list: %v", err)
	} else {
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
	if err := httputils.ValidateContentType(r, httputils.MediaTypeJSON); err != nil {
		return err
	}

	user, err := s.auth.GetUserID(r)
	if err != nil {
		return httputils.Errorf(http.StatusUnauthorized, "failed to get user ID: %v", err)
	}

	menu := r.PathValue("menu")
	pantry := r.PathValue("pantry")

	if _, err = s.db.LookupMenu(user, menu); errors.Is(err, fs.ErrNotExist) {
		return httputils.Errorf(http.StatusNotFound, "menu not found")
	} else if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to lookup menu: %v", err)
	}

	if _, err = s.db.LookupPantry(user, pantry); errors.Is(err, fs.ErrNotExist) {
		return httputils.Errorf(http.StatusNotFound, "pantry not found")
	} else if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to lookup pantry: %v", err)
	}

	out, err := io.ReadAll(r.Body)
	if err != nil {
		return httputils.Error(http.StatusBadRequest, "failed to read request")
	}
	r.Body.Close()

	sl := dbtypes.ShoppingList{
		User:   user,
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
	user, err := s.auth.GetUserID(r)
	if err != nil {
		return httputils.Errorf(http.StatusUnauthorized, "failed to get user ID: %v", err)
	}

	menu := r.PathValue("menu")
	pantry := r.PathValue("pantry")

	if err := s.db.DeleteShoppingList(user, menu, pantry); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to delete shopping list: %v", err)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

type shoppingListItem struct {
	ProductID product.ID `json:"product_id"`
	Name      string     `json:"name"`
	Done      bool       `json:"done"`
	Units     float32    `json:"units"`
	Packs     int        `json:"packs"`
	Cost      float32    `json:"cost"`
}

func (s *Service) computeShoppingList(log logger.Logger, menu dbtypes.Menu, pantry dbtypes.Pantry, done []product.ID) []shoppingListItem {
	need := menuneeds.ComputeNeeds(log, s.db, menu)

	slices.SortFunc(need, func(i, j recipe.Ingredient) int { return cmp.Compare(i.ProductID, j.ProductID) })
	slices.SortFunc(pantry.Contents, func(i, j recipe.Ingredient) int { return cmp.Compare(i.ProductID, j.ProductID) })
	slices.Sort(done)

	tmpList := menuneeds.Subtract(need, pantry.Contents)

	list := make([]shoppingListItem, 0, len(tmpList))
	utils.Zipper(tmpList, done,
		func(a recipe.Ingredient, id product.ID) int { return cmp.Compare(a.ProductID, id) },
		func(a recipe.Ingredient) {
			// This product is needed but not marked done
			p, ok := getProduct(log, s.db, a.ProductID)
			if ok {
				list = append(list, newItem(p, a.Amount, false))
			}
		},
		func(a recipe.Ingredient, id product.ID) {
			// This product is needed and marked done in the DB
			p, ok := getProduct(log, s.db, a.ProductID)
			if ok {
				list = append(list, newItem(p, a.Amount, true))
			}
		},
		func(id product.ID) {
			// This product is marked done but not needed
		})

	return list
}

func newItem(prod product.Product, units float32, isDone bool) shoppingListItem {
	packs := int(math.Ceil(float64(units / prod.BatchSize)))

	return shoppingListItem{
		ProductID: prod.ID,
		Name:      prod.Name,
		Units:     units,
		Packs:     packs,
		Cost:      float32(packs) * prod.Price,
		Done:      isDone,
	}
}

func getProduct(log logger.Logger, db database.DB, ID product.ID) (product.Product, bool) {
	p, err := db.LookupProduct(ID)
	if err != nil {
		log.Warningf("Failed to lookup product %d: %v", ID, err)
		return product.Product{}, false
	}
	return p, true
}
