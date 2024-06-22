package provider

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
)

type Service struct {
	settings Settings
}

type Settings struct {
	Enable bool
}

func (s Settings) Defaults() Settings {
	return Settings{
		Enable: true,
	}
}

func New(s Settings) Service {
	return Service{
		settings: s,
	}
}

func (s Service) Name() string {
	return "provider"
}

func (s Service) Path() string {
	return "/api/provider"
}

func (s Service) Enabled() bool {
	return s.settings.Enable
}

func (s Service) Handle(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGet(log, w, r)
	default:
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}
}

func (s Service) handleGet(_ logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if err := httputils.ValidateAccepts(r, httputils.MediaTypeJSON); err != nil {
		return err
	}

	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to parse query: %v", err)
	}

	prov := q.Get("provider")
	if prov == "" {
		return httputils.Errorf(http.StatusBadRequest, "missing query parameter: provider")
	}

	code := q.Get("id")
	if code == "" {
		return httputils.Errorf(http.StatusBadRequest, "missing query parameter: id")
	}

	provider, ok := providers.Lookup(prov)
	if !ok {
		return httputils.Errorf(http.StatusNotFound, "unknown provider %q", prov)
	}

	pid, err := parseProductID(provider, code)
	if err != nil {
		return err
	}

	price, err := provider.FetchPrice(r.Context(), pid)
	if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to fetch price: %v", err)
	}

	if err := json.NewEncoder(w).Encode(map[string]float32{"price": price}); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to write response: %v", err)
	}

	return nil
}

func parseProductID(provider providers.Provider, code string) (providers.ProductID, error) {
	pid := providers.ProductID{code}

	if provider.Name() == "Mercadona" {
		// Small cheat to avoid asking users for their postal code
		pid[1] = "bcn1"
	}

	if err := provider.ValidateID(pid); err != nil {
		return pid, httputils.Errorf(http.StatusBadRequest, "invalid code %q: %v", code, err)
	}

	return pid, nil
}
