package recipes

import (
	"encoding/json"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/auth"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
)

type Service struct {
	settings Settings
	db       database.DB
	auth     auth.Getter
}

type Settings struct {
	Enable bool
}

func (s Settings) Defaults() Settings {
	return Settings{
		Enable: true,
	}
}

func New(s Settings, db database.DB, auth auth.Getter) *Service {
	if !s.Enable {
		return nil
	}

	return &Service{
		settings: s,
		db:       db,
		auth:     auth,
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

	user, err := s.auth.GetUserID(r)
	if err != nil {
		return httputils.Errorf(http.StatusUnauthorized, "could not get user: %v", err)
	}

	recs, err := s.db.Recipes(user)
	if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not get recipes: %v", err)
	}

	type item struct {
		ID   recipe.ID `json:"id"`
		Name string    `json:"name"`
	}

	items := make([]item, 0)
	for _, r := range recs {
		items = append(items, item{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	if err := json.NewEncoder(w).Encode(items); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not encode recipes: %v", err)
	}

	log.Debugf("Responded with %d items", len(items))

	return nil
}
