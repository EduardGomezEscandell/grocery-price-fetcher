package session

import (
	"fmt"
	"net/http"
	"strings"

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
	return "/api/logout"
}

func (s LogoutService) Enabled() bool {
	return s.Enable
}

func (s LogoutService) Handle(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}

	auth := r.Header.Get("Authorization")
	if auth == "" {
		return httputils.Errorf(http.StatusUnauthorized, "authorization header not found")
	}

	var token string
	_, err := fmt.Fscanf(strings.NewReader(auth), "Bearer %s", &token)
	if err != nil {
		return httputils.Errorf(http.StatusUnauthorized, "could not parse authorization header: %v", err)
	}

	if err := s.sessions.Remove(token); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not remove session: %v", err)
	}

	return nil
}
