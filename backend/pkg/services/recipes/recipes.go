package recipes

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
	return "recipes"
}

func (s Service) Path() string {
	return "/api/recipes"
}

func (s Service) Enabled() bool {
	return s.settings.Enable
}

func (s *Service) Handle(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}

	if err := httputils.ValidateAccepts(r, httputils.MediaTypeJSON); err != nil {
		return err
	}

	recs, err := s.db.Recipes()
	if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to get recipes: %v", err)
	}

	names := make([]string, 0)
	for _, r := range recs {
		names = append(names, r.Name)
	}

	if err := json.NewEncoder(w).Encode(names); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to encode recipes: %v", err)
	}

	log.Debugf("Responded with %d items", len(names))

	return nil
}
