package recipes

import (
	"encoding/json"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/server/httputils"
	"github.com/sirupsen/logrus"
)

func Handler(db *database.DB) httputils.Handler {
	return func(log *logrus.Entry, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
		}

		var names []string
		for _, r := range db.Recipes {
			names = append(names, r.Name)
		}

		if err := json.NewEncoder(w).Encode(names); err != nil {
			return httputils.Errorf(http.StatusInternalServerError, "failed to encode recipes: %v", err)
		}

		log.Debugf("Responded with %d items", len(names))

		return nil
	}
}
