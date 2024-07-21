package session

import (
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/auth"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
)

type LogoutService struct {
	Enable   bool
	sessions *auth.Manager
}

func NewLogout(settings Settings, manager *auth.Manager) LogoutService {
	return LogoutService{
		Enable:   settings.Enable,
		sessions: manager,
	}
}

func (s LogoutService) Name() string {
	return "logout"
}

func (s LogoutService) Path() string {
	return "/api/auth/logout"
}

func (s LogoutService) Enabled() bool {
	return s.Enable
}

func (s LogoutService) Handle(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}

	token, err := s.sessions.GetToken(r)
	if err != nil {
		return httputils.Errorf(http.StatusUnauthorized, "could not get token: %v", err)
	}

	if token == "" {
		return nil
	}

	if err := s.sessions.Remove(token); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not remove session: %v", err)
	}

	return nil
}
