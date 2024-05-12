package jsondb

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
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/types"
)

type JSON struct {
	products []product.Product
	recipes  []types.Recipe
	menus    []types.Menu
	pantries []types.Pantry

	productsPath string
	recipesPath  string
	menusPath    string
	pantriesPath string

	log logger.Logger
	mu  sync.RWMutex
}

func New(ctx context.Context, log logger.Logger, options map[string]interface{}) (*JSON, error) {
	prods, errP := getStringOption(options, "products")
	recs, errR := getStringOption(options, "recipes")
	menus, errM := getStringOption(options, "menus")
	pants, errX := getStringOption(options, "pantries")

	if err := errors.Join(errP, errR, errM, errX); err != nil {
		return nil, fmt.Errorf("JSON database: %v", err)
	}

	db := &JSON{
		log: log,

		productsPath: prods,
		recipesPath:  recs,
		menusPath:    menus,
		pantriesPath: pants,
	}

	return db, errors.Join(
		load(db.productsPath, &db.products),
		load(db.recipesPath, &db.recipes),
		load(db.menusPath, &db.menus),
		load(db.pantriesPath, &db.pantries),
	)
}

func (db *JSON) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.save(); err != nil {
		return err
	}

	return nil
}

func (db *JSON) Products() []product.Product {
	db.mu.RLock()
	defer db.mu.RUnlock()

	out := make([]product.Product, len(db.products))
	copy(out, db.products)
	return out
}

func (db *JSON) LookupProduct(name string) (product.Product, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	i := slices.IndexFunc(db.products, func(p product.Product) bool {
		return p.Name == name
	})

	if i == -1 {
		return product.Product{}, false
	}

	return db.products[i], true
}

func (db *JSON) SetProduct(p product.Product) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	i := slices.IndexFunc(db.products, func(entry product.Product) bool {
		return entry.Name == p.Name
	})

	if i == -1 {
		db.products = append(db.products, p)
	} else {
		db.products[i] = p
	}

	if err := db.save(); err != nil {
		return err
	}

	return nil
}

func (db *JSON) Recipes() []types.Recipe {
	db.mu.RLock()
	defer db.mu.RUnlock()

	out := make([]types.Recipe, len(db.recipes))
	copy(out, db.recipes)
	return out
}

func (db *JSON) LookupRecipe(name string) (types.Recipe, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	i := slices.IndexFunc(db.recipes, func(p types.Recipe) bool {
		return p.Name == name
	})

	if i == -1 {
		return types.Recipe{}, false
	}

	return db.recipes[i], true
}

func (db *JSON) Menus() []types.Menu {
	db.mu.RLock()
	defer db.mu.RUnlock()

	out := make([]types.Menu, len(db.menus))
	copy(out, db.menus)
	return out
}

func (db *JSON) LookupMenu(name string) (types.Menu, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	i := slices.IndexFunc(db.menus, func(p types.Menu) bool {
		return p.Name == name
	})

	if i == -1 {
		return types.Menu{}, false
	}

	return db.menus[i], true
}

func (db *JSON) SetMenu(m types.Menu) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	i := slices.IndexFunc(db.menus, func(entry types.Menu) bool {
		return entry.Name == m.Name
	})

	if i == -1 {
		db.menus = append(db.menus, m)
	} else {
		db.menus[i] = m
	}

	if err := db.save(); err != nil {
		return err
	}
	return nil
}

func (db *JSON) Pantries() []types.Pantry {
	db.mu.RLock()
	defer db.mu.RUnlock()

	out := make([]types.Pantry, len(db.pantries))
	copy(out, db.pantries)
	return out
}

func (db *JSON) LookupPantry(name string) (types.Pantry, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	i := slices.IndexFunc(db.pantries, func(p types.Pantry) bool {
		return p.Name == name
	})

	if i == -1 {
		return types.Pantry{}, false
	}

	return db.pantries[i], true
}

func (db *JSON) SetPantry(p types.Pantry) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if p.Name == "" {
		p.Name = "default"
	}

	i := slices.IndexFunc(db.pantries, func(entry types.Pantry) bool {
		return entry.Name == p.Name
	})

	if i == -1 {
		db.pantries = append(db.pantries, p)
	} else {
		db.pantries[i] = p
	}

	if err := db.save(); err != nil {
		return err
	}
	return nil
}

func getStringOption(options map[string]any, key string) (string, error) {
	p, ok := options[key]
	if !ok {
		return "", fmt.Errorf("missing option %q", key)
	}

	path, ok := p.(string)
	if !ok {
		return "", fmt.Errorf("option %q is not a string", key)
	}

	return path, nil
}

func load(path string, ptr interface{}) error {
	out, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return fmt.Errorf("JSON database: %v", err)
	}

	if err := json.Unmarshal(out, ptr); err != nil {
		return fmt.Errorf("JSON database: could not unmarshal file %q: %v", path, err)
	}
	return nil
}

func (db *JSON) save() error {
	return errors.Join(
		save(db.log, db.productsPath, db.products),
		save(db.log, db.recipesPath, db.recipes),
		save(db.log, db.menusPath, db.menus),
		save(db.log, db.pantriesPath, db.pantries),
	)
}

func save(log logger.Logger, path string, structure interface{}) (err error) {
	b, err := newBackup(path)
	if err != nil {
		return err
	}
	defer b.onExit(&err, log)

	out, err := json.Marshal(structure)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, out, 0600); err != nil {
		return err
	}

	return nil
}

type backup struct {
	path string
	tmp  string
}

func newBackup(path string) (backup, error) {
	b := backup{
		path: path,
		tmp:  path + ".bak",
	}

	if err := os.Rename(b.path, b.tmp); err != nil {
		return b, fmt.Errorf("could not create backup: %v", err)
	}

	return b, nil
}

func (b backup) onExit(err *error, log logger.Logger) {
	if *err != nil {
		b.restore(log)
	} else {
		b.remove(log)
	}
}

func (b backup) restore(log logger.Logger) {
	if err := os.Rename(b.tmp, b.path); err != nil {
		log.Warnf("Could not restore backup: %v", b.path, err)
	}
}

func (b backup) remove(log logger.Logger) {
	if err := os.RemoveAll(b.tmp); err != nil {
		log.Warnf("Could not remove backup: %v", err)
	}
}
