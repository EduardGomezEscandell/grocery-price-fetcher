package pricing

import (
	"context"
	"sync"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/product"
)

type Service struct {
	ctx    context.Context
	cancel context.CancelFunc

	db  database.DB
	log logger.Logger
}

func New(ctx context.Context, log logger.Logger, db database.DB) *Service {
	ctx, cancel := context.WithCancel(ctx)

	return &Service{
		ctx:    ctx,
		cancel: cancel,
		db:     db,
		log:    log,
	}
}

func OneShot(ctx context.Context, log logger.Logger, db database.DB) {
	s := New(ctx, log, db)
	s.update()
	s.Stop()
}

func (s *Service) Stop() {
	s.cancel()
}

func (s *Service) Run() {
	s.update()

	go func() {
		ticker := time.NewTicker(time.Hour)
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
	s.log.Debug("Database: fetching prices")
	defer s.log.Debug("Database: prices fetch complete")

	var wg sync.WaitGroup
	for _, prod := range s.db.Products() {
		wg.Add(1)
		go func(i product.Product) {
			defer wg.Done()
			if err := prod.FetchPrice(s.ctx, s.log); err != nil {
				s.log.Warningf("Database price update: %v", err)
				return
			}

			if err := s.db.SetProduct(prod); err != nil {
				s.log.Warningf("Database price update: %v", err)
			}
		}(prod)
	}

	wg.Wait()
}
