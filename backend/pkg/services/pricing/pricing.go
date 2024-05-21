package pricing

import (
	"context"
	"sync"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
)

type Service struct {
	ctx    context.Context
	cancel context.CancelFunc

	settings Settings
	db       database.DB
	log      logger.Logger
}

type Settings struct {
	Enable      bool
	RefreshRate time.Duration
}

func (Settings) Defaults() Settings {
	return Settings{
		Enable:      true,
		RefreshRate: 6 * time.Hour,
	}
}

func New(ctx context.Context, s Settings, log logger.Logger, db database.DB) *Service {
	if !s.Enable {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)

	return &Service{
		settings: s,
		ctx:      ctx,
		cancel:   cancel,
		db:       db,
		log:      log,
	}
}

func (s Service) Enabled() bool {
	return s.settings.Enable
}

func OneShot(ctx context.Context, log logger.Logger, db database.DB) {
	s := New(ctx, Settings{}.Defaults(), log, db)
	s.update()
	s.Stop()
}

func (s *Service) Stop() {
	if s == nil {
		return
	}

	s.cancel()
}

func (s *Service) Run() {
	if !s.settings.Enable {
		return
	}

	s.update()

	if s.settings.RefreshRate == 0 {
		return
	}

	go func() {
		ticker := time.NewTicker(s.settings.RefreshRate)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.update()
			case <-s.ctx.Done():
				return
			}
		}
	}()
}

func (s *Service) update() {
	s.log.Debug("Pricing service: fetching prices")
	defer s.log.Debug("Pricing service: prices fetch complete")

	var wg sync.WaitGroup
	for _, prod := range s.db.Products() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := prod.FetchPrice(s.ctx); err != nil {
				s.log.Warningf("Database price update: %v", err)
				return
			}

			if err := s.db.SetProduct(prod); err != nil {
				s.log.Warningf("Database price update: %v", err)
			}
		}()
	}

	wg.Wait()
}
