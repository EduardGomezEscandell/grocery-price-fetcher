package helloworld

import (
	"fmt"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
)

type Service struct {
	settings Settings
}

type Settings struct {
	Enable bool
}

func (Settings) Defaults() Settings {
	return Settings{
		Enable: true,
	}
}

func New(settings Settings) *Service {
	return &Service{
		settings: settings,
	}
}

func (s Service) Name() string {
	return "helloworld"
}

func (s Service) Path() string {
	return "/api/helloworld"
}

func (s Service) Enabled() bool {
	return s.settings.Enable
}

func (s Service) Handle(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}

	if err := httputils.ValidateAccepts(r, "text/plain"); err != nil {
		return err
	}

	fmt.Fprintln(w, "Hello, world!")
	return nil
}
