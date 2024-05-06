package database

import (
	"context"
	"errors"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/database/jsonDB"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/recipe"
)

type DB interface {
	Products() []product.Product
	LookupProduct(name string) (product.Product, bool)
	SetProduct(p product.Product) error

	Recipes() []recipe.Recipe
	LookupRecipe(name string) (recipe.Recipe, bool)

	Close() error
}

type Settings struct {
	Type    string
	Options map[string]interface{}
}

func New(ctx context.Context, logger logger.Logger, s Settings) (DB, error) {
	switch s.Type {
	case "json":
		db, err := jsonDB.New(ctx, logger, s.Options)
		if err != nil {
			return nil, err
		}
		return db, nil
	default:
		return nil, errors.New("unknown database type")
	}
}
