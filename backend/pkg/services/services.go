package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/auth"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/blank"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/bonpreu"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/mercadona"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/frontend"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/helloworld"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/ingredientuse"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/menu"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/pantry"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/pricing"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/products"
	providersservice "github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/recipe"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/recipes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/session"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/shoppinglist"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/shoppingneeds"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/version"
)

type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc

	log          logger.Logger
	db           database.DB
	auth         *auth.Manager
	pricing      *pricing.Service
	httpServices map[string]HTTPService
	frontEnd     frontend.Service
}

type Settings struct {
	Database      database.Settings
	Auth          auth.Settings
	FrontEnd      frontend.Settings
	AuthLogin     session.Settings
	AuthLogout    session.Settings
	AuthRefresh   session.Settings
	HelloWorld    helloworld.Settings
	IngredientUse ingredientuse.Settings
	Menu          menu.Settings
	Pantry        pantry.Settings
	Pricing       pricing.Settings
	Products      products.Settings
	Providers     providersservice.Settings
	Recipe        recipe.Settings
	Recipes       recipes.Settings
	ShoppingList  shoppinglist.Settings
	ShoppingNeeds shoppingneeds.Settings
	Version       version.Settings
}

func (Settings) Defaults() Settings {
	return Settings{
		Auth:          auth.Settings{}.Defaults(),
		Database:      database.Settings{}.Defaults(),
		FrontEnd:      frontend.Settings{}.Defaults(),
		AuthLogin:     session.Settings{}.Defaults(),
		AuthLogout:    session.Settings{}.Defaults(),
		AuthRefresh:   session.Settings{}.Defaults(),
		HelloWorld:    helloworld.Settings{}.Defaults(),
		IngredientUse: ingredientuse.Settings{}.Defaults(),
		Menu:          menu.Settings{}.Defaults(),
		Pantry:        pantry.Settings{}.Defaults(),
		Pricing:       pricing.Settings{}.Defaults(),
		Products:      products.Settings{}.Defaults(),
		Providers:     providersservice.Settings{}.Defaults(),
		Recipe:        recipe.Settings{}.Defaults(),
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

	auth, err := auth.NewManager(ctx, settings.Auth, logger, db)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("could not create auth manager: %v", err)
	}

	m := &Manager{
		ctx:    ctx,
		cancel: cancel,

		db:      db,
		auth:    auth,
		log:     logger,
		pricing: pricing.New(ctx, settings.Pricing, logger, db),

		httpServices: map[string]HTTPService{},
		frontEnd:     frontend.New(settings.FrontEnd),
	}

	for _, s := range []HTTPService{
		session.NewLogin(settings.AuthLogin, auth),
		session.NewRefresh(settings.AuthRefresh, auth),
		session.NewLogout(settings.AuthLogout, auth),
		helloworld.New(settings.HelloWorld),
		ingredientuse.New(settings.IngredientUse, db, auth),
		menu.New(settings.Menu, db, auth),
		pantry.New(settings.Pantry, db, auth),
		products.New(settings.Products, db),
		providersservice.New(settings.Providers),
		recipe.New(settings.Recipe, db, auth),
		recipes.New(settings.Recipes, db, auth),
		shoppinglist.New(settings.ShoppingList, db, auth),
		shoppingneeds.New(settings.ShoppingNeeds, db, auth),
		version.New(settings.Version),
	} {
		m.httpServices[s.Name()] = s
	}

	if err := m.auth.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("could not start auth manager: %v", err)
	}

	m.pricing.Start()

	return m, nil
}

func (m *Manager) Register(registerer func(endpoint string, handler func(http.ResponseWriter, *http.Request))) {
	registerer(m.frontEnd.Path(), m.frontEnd.HandleHTTP)

	for _, p := range m.httpServices {
		if !p.Enabled() {
			m.log.Infof("Skipping dynamic endpoint %s", p.Name())
			continue
		}
		registerer(p.Path(), httputils.HandleRequest(m.log.WithField("service", p.Name()), p.Handle))
	}
}

func (m *Manager) Stop() error {
	m.log.Info("Stopping services")
	defer m.cancel()

	m.pricing.Stop()
	m.auth.Stop()

	if err := m.db.Close(); err != nil {
		return fmt.Errorf("could not close database: %v", err)
	}

	return nil
}
