package database

import (
	"context"
	"errors"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/jsondb"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/mysql"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
	"gopkg.in/yaml.v3"
)

type DB interface {
	Products() ([]product.Product, error)
	LookupProduct(ID product.ID) (product.Product, error)
	SetProduct(p product.Product) (product.ID, error)
	DeleteProduct(ID product.ID) error

	Recipes() ([]recipe.Recipe, error)
	LookupRecipe(name recipe.ID) (recipe.Recipe, error)
	SetRecipe(r recipe.Recipe) (recipe.ID, error)
	DeleteRecipe(name recipe.ID) error

	Menus(user string) ([]dbtypes.Menu, error)
	LookupMenu(user, name string) (dbtypes.Menu, error)
	SetMenu(m dbtypes.Menu) error
	DeleteMenu(user, name string) error

	Pantries(user string) ([]dbtypes.Pantry, error)
	LookupPantry(user, name string) (dbtypes.Pantry, error)
	SetPantry(p dbtypes.Pantry) error
	DeletePantry(user, name string) error

	ShoppingLists(user string) ([]dbtypes.ShoppingList, error)
	LookupShoppingList(user, menu, pantry string) (dbtypes.ShoppingList, error)
	SetShoppingList(m dbtypes.ShoppingList) error
	DeleteShoppingList(user, menu, pantry string) error

	Close() error
}

type Settings struct {
	Type    string
	Options any
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
		db, err := jsondb.New(ctx, logger, s.Options.(jsondb.Settings))
		if err != nil {
			return nil, err
		}
		return db, nil
	case "mysql":
		db, err := mysql.New(ctx, logger, s.Options.(mysql.Settings))
		if err != nil {
			return nil, err
		}
		return db, nil
	default:
		return nil, errors.New("unknown database type")
	}
}

func (s *Settings) UnmarshalYAML(node *yaml.Node) error {
	var raw struct {
		Type    string
		Options yaml.Node
	}

	if err := node.Decode(&raw); err != nil {
		return err
	}

	switch raw.Type {
	case "":
		*s = s.Defaults()
	case "json":
		s.Type = raw.Type
		opt := jsondb.DefaultSettings()
		if err := raw.Options.Decode(&opt); err != nil {
			return err
		}
		s.Options = opt
	case "mysql":
		s.Type = raw.Type
		opt := mysql.DefaultSettings()
		if err := raw.Options.Decode(&opt); err != nil {
			return err
		}
		s.Options = opt
	default:
		return errors.New("unknown database type")
	}

	return nil
}
