package session

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/auth"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
)

type LoginService struct {
	Enable   bool
	sessions *auth.Manager
}

func NewLogin(settings Settings, manager *auth.Manager) LoginService {
	return LoginService{
		Enable:   settings.Enable,
		sessions: manager,
	}
}

func (s LoginService) Name() string {
	return "login"
}

func (s LoginService) Path() string {
	return "/api/login"
}

func (s LoginService) Enabled() bool {
	return s.Enable
}

func (s LoginService) Handle(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}

	if err := httputils.ValidateContentType(r, httputils.MediaTypeJSON); err != nil {
		return err
	}

	if err := httputils.ValidateAccepts(r, httputils.MediaTypeText); err != nil {
		return err
	}

	var data struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return httputils.Errorf(http.StatusBadRequest, "could not decode request body: %v", err)
	}

	token, err := s.sessions.NewSession(data.Code)
	if err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not create session: %v", err)
	}

	if _, err := fmt.Fprintf(w, "Bearer %s", token); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not write response: %v", err)
	}

	return nil
}
