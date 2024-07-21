package session

import (
	"fmt"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/auth"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
)

type RefreshService struct {
	Enable   bool
	sessions *auth.Manager
}

func NewRefresh(settings Settings, manager *auth.Manager) RefreshService {
	return RefreshService{
		Enable:   settings.Enable,
		sessions: manager,
	}
}

func (s RefreshService) Name() string {
	return "refresh"
}

func (s RefreshService) Path() string {
	return "/api/auth/refresh"
}

func (s RefreshService) Enabled() bool {
	return s.Enable
}

func (s RefreshService) Handle(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}

	if err := httputils.ValidateAccepts(r, httputils.MediaTypeText); err != nil {
		return err
	}

	token, err := s.sessions.GetToken(r)
	if err != nil {
		return httputils.Errorf(http.StatusBadRequest, "could not get token: %v", err)
	}

	if token == "" {
		return httputils.Error(http.StatusUnauthorized, "token is empty")
	}

	sess, err := s.sessions.Get(token)
	if err != nil {
		return httputils.Errorf(http.StatusUnauthorized, "could not get session: %v", err)
	}

	if _, err := fmt.Fprintf(w, "Bearer %s", sess.ID); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not write response: %v", err)
	}

	return nil
}
