package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/blank"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/bonpreu"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/mercadona"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/helloworld"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/ingredientuse"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/menu"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/pantry"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/pricing"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/recipes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/shoppinglist"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/shoppingneeds"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/version"
)

type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc

	log          logger.Logger
	db           database.DB
	pricing      *pricing.Service
	httpServices map[string]HTTPService
}

type Settings struct {
	Database      database.Settings
	HelloWorld    helloworld.Settings
	IngredientUse ingredientuse.Settings
	Menu          menu.Settings
	Pantry        pantry.Settings
	Pricing       pricing.Settings
	Recipes       recipes.Settings
	ShoppingList  shoppinglist.Settings
	ShoppingNeeds shoppingneeds.Settings
	Version       version.Settings
}

func (Settings) Defaults() Settings {
	return Settings{
		Database:      database.Settings{}.Defaults(),
		HelloWorld:    helloworld.Settings{}.Defaults(),
		IngredientUse: ingredientuse.Settings{}.Defaults(),
		Menu:          menu.Settings{}.Defaults(),
		Pantry:        pantry.Settings{}.Defaults(),
		Pricing:       pricing.Settings{}.Defaults(),
		Recipes:       recipes.Settings{}.Defaults(),
		ShoppingList:  shoppinglist.Settings{}.Defaults(),
		ShoppingNeeds: shoppingneeds.Settings{}.Defaults(),
		Version:       version.Settings{}.Defaults(),
	}
}

type HTTPService interface {
	Name() string
	Path() string
	Handle(logger.Logger, http.ResponseWriter, *http.Request) error
	Enabled() bool
}

func New(ctx context.Context, logger logger.Logger, settings Settings) (*Manager, error) {
	ctx, cancel := context.WithCancel(ctx)

	providers.Register(blank.Provider{})
	providers.Register(bonpreu.New(logger))
	providers.Register(mercadona.New(logger))

	db, err := database.New(ctx, logger, settings.Database)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("could not load database: %v", err)
	}

	m := &Manager{
		ctx:    ctx,
		cancel: cancel,

		db:      db,
		log:     logger,
		pricing: pricing.New(ctx, settings.Pricing, logger, db),

		httpServices: map[string]HTTPService{},
	}

	for _, s := range []HTTPService{
		helloworld.New(settings.HelloWorld),
		ingredientuse.New(settings.IngredientUse, db),
		menu.New(settings.Menu, db),
		pantry.New(settings.Pantry, db),
		recipes.New(settings.Recipes, db),
		shoppinglist.New(settings.ShoppingList, db),
		shoppingneeds.New(settings.ShoppingNeeds, db),
		version.New(settings.Version),
	} {
		m.httpServices[s.Name()] = s
	}

	return m, nil
}

func (s *Manager) Register(registerer func(endpoint string, handler httputils.Handler)) {
	for _, p := range s.httpServices {
		if !p.Enabled() {
			s.log.Infof("Skipping dynamic endpoint %s", p.Name())
			continue
		}

		registerer(p.Path(), p.Handle)
	}
}

func (s *Manager) Run() {
	s.pricing.Run()
}

func (s *Manager) Stop() {
	s.log.Info("Stopping services")

	s.pricing.Stop()
	s.cancel()
}
