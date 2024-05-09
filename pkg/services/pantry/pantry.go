package pantry

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/types"
)

type Service struct {
	db database.DB
}

func New(db database.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Handle(log logger.Logger, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGet(log, w, r)
	case http.MethodPost:
		return s.handlePost(log, w, r)
	default:
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}
}

func (s *Service) handleGet(_ logger.Logger, w http.ResponseWriter, _ *http.Request) error {
	if err := json.NewEncoder(w).Encode(s.db.Pantries()); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not write menus to output: %w", err)
	}

	return nil
}

func (s *Service) handlePost(_ logger.Logger, _ http.ResponseWriter, r *http.Request) error {
	out, err := io.ReadAll(r.Body)
	if err != nil {
		return httputils.Error(http.StatusBadRequest, "failed to read request")
	}
	r.Body.Close()

	var pantry types.Pantry
	if err := json.Unmarshal(out, &pantry); err != nil {
		return httputils.Errorf(http.StatusBadRequest, "could not unmarshal pantry: %w", err)
	}

	if err := s.db.SetPantry(pantry); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not set pantry: %w", err)
	}

	return nil
}
