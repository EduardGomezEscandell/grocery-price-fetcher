package jsonDB

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"sync"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/recipe"
)

type JSON struct {
	Prods []product.Product `json:"products"`
	Recs  []recipe.Recipe   `json:"recipes"`

	path string
	log  logger.Logger
	mu   sync.RWMutex
}

func New(ctx context.Context, log logger.Logger, options map[string]interface{}) (*JSON, error) {
	var path string

	if p, ok := options["path"]; !ok {
		return nil, errors.New("missing path option")
	} else if path, ok = p.(string); !ok {
		return nil, errors.New("path option must be a string")
	}

	out, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	db := &JSON{
		path: path,
		log:  log,
	}

	if err := json.Unmarshal(out, db); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *JSON) Close() error {
	if err := db.save(); err != nil {
		return err
	}

	return nil
}

func (db *JSON) save() error {
	if err := os.Rename(db.path, db.path+".bak"); err != nil {
		db.log.Warnf("could not create backup of database: %v", err)
	}

	out, err := json.Marshal(db)
	if err != nil {
		return err
	}

	if err := os.WriteFile(db.path, out, 0600); err != nil {
		return err
	}

	return nil
}

func (db *JSON) Products() []product.Product {
	db.mu.RLock()
	defer db.mu.RUnlock()

	out := make([]product.Product, len(db.Prods))
	copy(out, db.Prods)
	return out
}

func (db *JSON) LookupProduct(name string) (product.Product, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	i := slices.IndexFunc(db.Prods, func(p product.Product) bool {
		return p.Name == name
	})

	if i == -1 {
		return product.Product{}, false
	}

	return db.Prods[i], true
}

func (db *JSON) Validate() error {
	for _, r := range db.Recs {
		for _, i := range r.Ingredients {
			if _, ok := db.LookupProduct(i.Name); !ok {
				return fmt.Errorf("invalid database: recipe %s: ingredient %q is not registered", r.Name, i.Name)
			}
		}
	}

	return nil
}

func (db *JSON) SetProduct(p product.Product) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	i := slices.IndexFunc(db.Prods, func(entry product.Product) bool {
		return entry.Name == p.Name
	})

	if i == -1 {
		db.Prods = append(db.Prods, p)
	} else {
		db.Prods[i] = p
	}

	if err := db.save(); err != nil {
		return err
	}

	return nil
}

func (db *JSON) Recipes() []recipe.Recipe {
	db.mu.RLock()
	defer db.mu.RUnlock()

	out := make([]recipe.Recipe, len(db.Recs))
	copy(out, db.Recs)
	return out
}

func (db *JSON) LookupRecipe(name string) (recipe.Recipe, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	i := slices.IndexFunc(db.Recs, func(p recipe.Recipe) bool {
		return p.Name == name
	})

	if i == -1 {
		return recipe.Recipe{}, false
	}

	return db.Recs[i], true
}
