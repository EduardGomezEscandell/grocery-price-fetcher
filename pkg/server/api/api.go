package api

import (
	"fmt"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/server/api/menu"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/server/api/recipes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/server/httputils"
	"github.com/sirupsen/logrus"
)

func RegisterEndpoints(db *database.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/hello-world", httputils.HandleRequest(helloWorldHandler))
	mux.HandleFunc("/api/recipes", httputils.HandleRequest(recipes.Handler(db)))
	mux.HandleFunc("/api/menu", httputils.HandleRequest(menu.Handler(db)))
	return mux
}

func helloWorldHandler(_ *logrus.Entry, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}

	fmt.Fprintln(w, "Hello, world!")
	return nil
}
