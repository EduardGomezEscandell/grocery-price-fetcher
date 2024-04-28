package menu

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/formatter"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/menu"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/server/httputils"
	"github.com/sirupsen/logrus"
)

// RequestData is the data structure that the API expects to receive.
type RequestData struct {
	Menu   menu.Menu          `json:",omitempty"`
	Pantry []menu.ProductData `json:",omitempty"`
	Format string             `json:",omitempty"`
}

func Handler(db *database.DB) httputils.Handler {
	return func(log *logrus.Entry, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPost {
			return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
		}

		out, err := io.ReadAll(r.Body)
		if err != nil {
			return httputils.Error(http.StatusBadRequest, "failed to read request")
		}
		r.Body.Close()

		var data RequestData
		if err := json.Unmarshal(out, &data); err != nil {
			return httputils.Errorf(http.StatusBadRequest, "failed to unmarshal request: %v:\n%s", err, string(out))
		}

		log.Debugf("Received request with %d days and %d items in the pantry: ", len(data.Menu.Days), len(data.Pantry))

		if data.Format == "" {
			data.Format = "table"
		}
		f, err := formatter.New(data.Format)
		if err != nil {
			return httputils.Errorf(http.StatusBadRequest, "failed to create formatter: %v", err)
		}

		log.Debug("Selected formatter: ", data.Format)

		shoppingList, err := data.Menu.Compute(db, data.Pantry)
		if err != nil {
			return httputils.Errorf(http.StatusBadRequest, "failed to compute shopping list: %v", err)
		}

		log.Debug("Computed shopping list")

		if err := f.PrintHead(w, "Product", "Amount", "Unit Cost"); err != nil {
			return httputils.Errorf(http.StatusInternalServerError, "could not write header to output: %w", err)
		}

		i := 0
		for _, p := range shoppingList {
			if p.Amount == 0 {
				continue
			}

			if err := f.PrintRow(w, map[string]any{
				"Product":   p.Name,
				"Amount":    p.Amount,
				"Unit Cost": formatter.Euro(p.UnitCost),
			}); err != nil {
				return httputils.Errorf(http.StatusInternalServerError, "could not write results to output: %w", err)
			}
			i++
		}

		log.Debugf("Responded with %d items", i)

		if err := f.PrintTail(w); err != nil {
			return httputils.Errorf(http.StatusInternalServerError, "could not write footer to output: %w", err)
		}

		return nil
	}
}
