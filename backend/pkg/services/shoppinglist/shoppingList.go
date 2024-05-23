package shoppinglist

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/types"
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

func (s *Service) handleGet(_ logger.Logger, w http.ResponseWriter, _ *http.Request) error {
	if err := json.NewEncoder(w).Encode(s.db.ShoppingLists()); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not encode response: %v", err)
	}
	return nil
}

func (s *Service) handlePost(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
	out, err := io.ReadAll(r.Body)
	if err != nil {
		return httputils.Error(http.StatusBadRequest, "failed to read request")
	}
	r.Body.Close()

	var sl types.ShoppingList
	if err := json.Unmarshal(out, &sl); err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to unmarshal request: %v:\n%s", err, string(out))
	}

	sl.TimeStamp = time.Now().UTC().Format(time.RFC3339)

	if err := s.db.SetShoppingList(sl); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not add shopping list: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}
