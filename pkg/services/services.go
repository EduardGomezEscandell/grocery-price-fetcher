package services

import (
	"context"
	"fmt"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/providers/bonpreu"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/providers/mercadona"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/services/helloworld"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/services/menu"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/services/pantry"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/services/pricing"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/services/recipes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/services/version"
)

type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc

	db  database.DB
	log logger.Logger

	menu    *menu.Service
	recipes *recipes.Service
	pricing *pricing.Service
	pantry  *pantry.Service
}

func New(ctx context.Context, logger logger.Logger, DBsettings database.Settings) (*Manager, error) {
	ctx, cancel := context.WithCancel(ctx)

	providers.Register(bonpreu.New(logger))
	providers.Register(mercadona.New(logger))

	db, err := database.New(ctx, logger, DBsettings)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("could not load database: %v", err)
	}

	return &Manager{
		ctx:    ctx,
		cancel: cancel,

		db:  db,
		log: logger,

		menu:    menu.New(db),
		pantry:  pantry.New(db),
		pricing: pricing.New(ctx, logger, db),
		recipes: recipes.New(db),
	}, nil
}

func (s *Manager) Register(registerer func(endpoint string, handler httputils.Handler)) {
	registerer("/api/helloworld", helloworld.Service{}.Handle)
	registerer("/api/menu", s.menu.Handle)
	registerer("/api/pantry", s.pantry.Handle)
	registerer("/api/recipes", s.recipes.Handle)
	registerer("/api/version", version.Service{}.Handle)
}

func (s *Manager) Run() {
	s.pricing.Run()
}

func (s *Manager) Stop() {
	s.log.Info("Stopping services")

	s.pricing.Stop()
	s.cancel()
}
