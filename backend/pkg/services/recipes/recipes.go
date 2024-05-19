package recipes

import (
	"encoding/json"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
)

type Service struct {
	db database.DB
}

func New(db database.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Handle(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}

	var names []string
	for _, r := range s.db.Recipes() {
		names = append(names, r.Name)
	}

	if err := json.NewEncoder(w).Encode(names); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to encode recipes: %v", err)
	}

	log.Debugf("Responded with %d items", len(names))

	return nil
}
