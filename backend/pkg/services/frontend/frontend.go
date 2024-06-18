package frontend

import (
	"errors"
	"net/http"
	"os"
	"path"
	"strings"
)

type Service struct {
	settings Settings
	fs       http.Dir
}

type Settings struct {
	Path    string
	Enable bool
}

func (s Settings) Defaults() Settings {
	return Settings{
		Path:    "/run/secrets/frontend",
		Enable: true,
	}
}

func New(s Settings) Service {
	return Service{
		settings: s,
		fs:       http.Dir(s.Path),
	}
}

func (s Service) Path() string {
	return "/"
}

func (s Service) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	if p, ok := s.resolvePath(w, r.URL.Path); !ok {
		return
	} else {
		r.URL.Path = p
	}

	http.FileServer(s.fs).ServeHTTP(w, r)
}

func (s Service) resolvePath(w http.ResponseWriter, urlPath string) (string, bool) {
	// Prevent path traversal
	if strings.Contains(urlPath, "..") {
		http.Error(w, "path is invalid", http.StatusBadRequest)
		return "", false
	}

	// Avoid serving API endpoints
	if strings.HasPrefix(urlPath, "/api") {
		http.NotFound(w, nil)
		return "", false
	}

	// Check if file exists
	_, err := os.Stat(path.Join(s.settings.Path, urlPath))
	if err == nil {
		return urlPath, true
	}

	// Unexplained error
	if !errors.Is(err, os.ErrNotExist) {
		http.Error(w, "could stat path", http.StatusInternalServerError)
		return "", false
	}

	// Path does not exist: try to serve its parent
	// to allow client-side routing
	if urlPath == "/" {
		return urlPath, true
	}

	return s.resolvePath(w, path.Dir(urlPath))
}
