package database

import (
	"context"
	"errors"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/jsondb"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/types"
)

type DB interface {
	Products() []product.Product
	LookupProduct(name string) (product.Product, bool)
	SetProduct(p product.Product) error

	Recipes() []types.Recipe
	LookupRecipe(name string) (types.Recipe, bool)

	Menus() []types.Menu
	LookupMenu(name string) (types.Menu, bool)
	SetMenu(m types.Menu) error

	Pantries() []types.Pantry
	LookupPantry(name string) (types.Pantry, bool)
	SetPantry(p types.Pantry) error

	Close() error
}

type Settings struct {
	Type    string
	Options map[string]interface{}
}

func (s Settings) Defaults() Settings {
	return Settings{
		Type:    "json",
		Options: jsondb.DefaultSettings(),
	}
}

func New(ctx context.Context, logger logger.Logger, s Settings) (DB, error) {
	switch s.Type {
	case "json":
		db, err := jsondb.New(ctx, logger, s.Options)
		if err != nil {
			return nil, err
		}
		return db, nil
	default:
		return nil, errors.New("unknown database type")
	}
}
