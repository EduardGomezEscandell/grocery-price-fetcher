package database

import (
	"context"
	"slices"
	"sync"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/recipe"
	log "github.com/sirupsen/logrus"
)

type DB struct {
	Products []product.Product
	Recipes  []recipe.Recipe
}

func (db *DB) LookupProduct(name string) (*product.Product, bool) {
	i := slices.IndexFunc(db.Products, func(p product.Product) bool {
		return p.Name == name
	})

	if i == -1 {
		return nil, false
	}

	return &db.Products[i], true
}

func (db *DB) LookupRecipe(name string) (*recipe.Recipe, bool) {
	i := slices.IndexFunc(db.Recipes, func(p recipe.Recipe) bool {
		return p.Name == name
	})

	if i == -1 {
		return nil, false
	}

	return &db.Recipes[i], true
}

func (db *DB) UpdatePrices(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	for i := range db.Products {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := db.Products[i].FetchPrice(ctx)
			if err != nil {
				log.Warningf("Database price update: %v", err)
			}
		}(i)
	}

	wg.Wait()
}
