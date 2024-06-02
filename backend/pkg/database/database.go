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
	Products() ([]product.Product, error)
	LookupProduct(name string) (product.Product, bool)
	SetProduct(p product.Product) error
	DeleteProduct(name string) error

	Recipes() ([]types.Recipe, error)
	LookupRecipe(name string) (types.Recipe, bool)
	SetRecipe(r types.Recipe) error
	DeleteRecipe(name string) error

	Menus() ([]types.Menu, error)
	LookupMenu(name string) (types.Menu, bool)
	SetMenu(m types.Menu) error
	DeleteMenu(name string) error

	Pantries() ([]types.Pantry, error)
	LookupPantry(name string) (types.Pantry, bool)
	SetPantry(p types.Pantry) error
	DeletePantry(name string) error

	ShoppingLists() ([]types.ShoppingList, error)
	LookupShoppingList(name string) (types.ShoppingList, bool)
	SetShoppingList(m types.ShoppingList) error
	DeleteShoppingList(name string) error

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
