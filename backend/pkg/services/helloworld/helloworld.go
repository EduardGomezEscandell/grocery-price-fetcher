package helloworld

import (
	"fmt"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
)

type Service struct{}

func (s Service) Handle(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}

	fmt.Fprintln(w, "Hello, world!")
	return nil
}
