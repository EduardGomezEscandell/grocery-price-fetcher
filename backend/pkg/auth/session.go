package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/auth/google"
)

type Manager struct {
	app      google.Application
	sessions map[string]google.JWTDecoded
	mu       sync.Mutex
}

type Settings struct {
	ClientID         string `yaml:"client-id"`
	ClientSecretFile string `yaml:"client-secret-file"`
	RedirectURI      string `yaml:"redirect-uri"`
}

func (Settings) Defaults() Settings {
	return Settings{
		RedirectURI: "https://localhost",
	}
}

func NewManager(s Settings) (*Manager, error) {
	out, err := os.ReadFile(s.ClientSecretFile)
	if err != nil {
		return nil, fmt.Errorf("could not read client secret file: %w", err)
	} else if len(out) == 0 {
		return nil, errors.New("empty client secret file")
	}

	return &Manager{
		sessions: make(map[string]google.JWTDecoded),
		app: google.Application{
			ClientID:     s.ClientID,
			ClientSecret: strings.TrimSpace(string(out)),
			RedirectURI:  s.RedirectURI,
		},
	}, nil
}

func (m *Manager) NewSession(code string) (string, error) {
	d, err := m.app.Login(code)
	if err != nil {
		return "", fmt.Errorf("could not login to Google: %v", err)
	}

	if d.TokenType != "Bearer" {
		return "", fmt.Errorf("unexpected token type: %s", d.TokenType)
	}

	id, err := google.DecodeJWT(d.ID)
	if err != nil {
		return "", fmt.Errorf("could not decode JWT: %v", err)
	}

	if err := m.Add(d.AccessToken, id); err != nil {
		return "", fmt.Errorf("could not add session: %v", err)
	}

	fmt.Println("New session created for", id.Subject)

	return d.AccessToken, nil
}

func (m *Manager) Add(token string, info google.JWTDecoded) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.app.Validate(info); err != nil {
		return err
	}

	m.sessions[token] = info

	exp := time.Unix(info.NotAfter, 0)
	time.AfterFunc(time.Until(exp), func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		delete(m.sessions, token)
	})

	return nil
}

func (m *Manager) Remove(token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, token)
	return nil
}

func (m *Manager) Get(token string) (string, error) {
	m.mu.Lock()
	info, ok := m.sessions[token]
	m.mu.Unlock()

	if !ok {
		return "", errors.New("session not found")
	}

	if info.Expired() {
		return "", errors.New("session expired")
	}

	return info.Subject, nil
}

type Getter interface {
	GetUserID(r *http.Request) (string, error)
}

func (m *Manager) GetUserID(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", errors.New("missing token")
	}

	var token string
	if _, err := fmt.Sscanf(auth, "Bearer %s", &token); err != nil {
		return "", errors.New("could not parse token")
	}

	if token == "" {
		return "", errors.New("empty token")
	}

	return m.Get(token)
}
