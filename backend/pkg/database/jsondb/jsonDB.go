package jsondb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
)

type JSON struct {
	products      []product.Product
	recipes       []dbtypes.Recipe
	menus         []dbtypes.Menu
	pantries      []dbtypes.Pantry
	shoppingLists []dbtypes.ShoppingList

	productsPath      string
	recipesPath       string
	menusPath         string
	pantriesPath      string
	shoppingListsPath string

	log logger.Logger
	mu  sync.RWMutex
}

type Settings struct {
	Products      string
	Recipes       string
	Menus         string
	Pantries      string
	ShoppingLists string
}

func DefaultSettings() Settings {
	return DefaultSettingsPath("/mnt/grocery-price-fetcher")
}

func DefaultSettingsPath(root string) Settings {
	return Settings{
		Products:      filepath.Join(root, "products.json"),
		Recipes:       filepath.Join(root, "recipes.json"),
		Menus:         filepath.Join(root, "menus.json"),
		Pantries:      filepath.Join(root, "pantries.json"),
		ShoppingLists: filepath.Join(root, "shoppingLists.json"),
	}
}

func New(ctx context.Context, log logger.Logger, s Settings) (*JSON, error) {
	db := &JSON{
		log:               log,
		productsPath:      s.Products,
		recipesPath:       s.Recipes,
		menusPath:         s.Menus,
		pantriesPath:      s.Pantries,
		shoppingListsPath: s.ShoppingLists,
	}

	log = log.WithField("database", "json")
	log.Tracef("Loading database")

	return db, errors.Join(
		load(db.productsPath, &db.products),
		load(db.recipesPath, &db.recipes),
		load(db.menusPath, &db.menus),
		load(db.pantriesPath, &db.pantries),
		load(db.shoppingListsPath, &db.shoppingLists),
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

func (db *JSON) Products() ([]product.Product, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	out := make([]product.Product, len(db.products))
	copy(out, db.products)
	return out, nil
}

func (db *JSON) LookupProduct(ID product.ID) (product.Product, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	i := slices.IndexFunc(db.products, func(p product.Product) bool {
		return p.ID == ID
	})

	if i == -1 {
		return product.Product{}, fs.ErrNotExist
	}

	return db.products[i], nil
}

func (db *JSON) SetProduct(p product.Product) (product.ID, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if p.ID == 0 {
		p.ID = product.NewRandomID()
	}

	i := slices.IndexFunc(db.products, func(entry product.Product) bool {
		return entry.ID == p.ID
	})

	if i == -1 {
		db.products = append(db.products, p)
	} else {
		db.products[i] = p
	}

	if err := db.save(); err != nil {
		return 0, err
	}

	return p.ID, nil
}

func (db *JSON) DeleteProduct(ID product.ID) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	i := slices.IndexFunc(db.products, func(p product.Product) bool {
		return p.ID == ID
	})

	if i == -1 {
		return fmt.Errorf("product with ID %d not found", ID)
	}

	db.products = append(db.products[:i], db.products[i+1:]...)

	if err := db.save(); err != nil {
		return err
	}

	return nil
}

func (db *JSON) Recipes() ([]dbtypes.Recipe, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	out := make([]dbtypes.Recipe, len(db.recipes))
	copy(out, db.recipes)
	return out, nil
}

func (db *JSON) LookupRecipe(name string) (dbtypes.Recipe, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	i := slices.IndexFunc(db.recipes, func(p dbtypes.Recipe) bool {
		return p.Name == name
	})

	if i == -1 {
		return dbtypes.Recipe{}, false
	}

	return db.recipes[i], true
}

func (db *JSON) SetRecipe(r dbtypes.Recipe) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	i := slices.IndexFunc(db.recipes, func(entry dbtypes.Recipe) bool {
		return entry.Name == r.Name
	})

	if i == -1 {
		db.recipes = append(db.recipes, r)
	} else {
		db.recipes[i] = r
	}

	if err := db.save(); err != nil {
		return err
	}
	return nil
}

func (db *JSON) DeleteRecipe(name string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	i := slices.IndexFunc(db.recipes, func(p dbtypes.Recipe) bool {
		return p.Name == name
	})

	if i == -1 {
		return fmt.Errorf("recipe %q not found", name)
	}

	db.recipes = append(db.recipes[:i], db.recipes[i+1:]...)

	if err := db.save(); err != nil {
		return err
	}

	return nil
}

func (db *JSON) Menus() ([]dbtypes.Menu, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	out := make([]dbtypes.Menu, len(db.menus))
	copy(out, db.menus)
	return out, nil
}

func (db *JSON) LookupMenu(name string) (dbtypes.Menu, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	i := slices.IndexFunc(db.menus, func(p dbtypes.Menu) bool {
		return p.Name == name
	})

	if i == -1 {
		return dbtypes.Menu{}, false
	}

	return db.menus[i], true
}

func (db *JSON) SetMenu(m dbtypes.Menu) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	i := slices.IndexFunc(db.menus, func(entry dbtypes.Menu) bool {
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

func (db *JSON) DeleteMenu(name string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	i := slices.IndexFunc(db.menus, func(p dbtypes.Menu) bool {
		return p.Name == name
	})

	if i == -1 {
		return fmt.Errorf("menu %q not found", name)
	}

	db.menus = append(db.menus[:i], db.menus[i+1:]...)

	if err := db.save(); err != nil {
		return err
	}

	return nil
}

func (db *JSON) Pantries() ([]dbtypes.Pantry, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	out := make([]dbtypes.Pantry, len(db.pantries))
	copy(out, db.pantries)
	return out, nil
}

func (db *JSON) LookupPantry(name string) (dbtypes.Pantry, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	i := slices.IndexFunc(db.pantries, func(p dbtypes.Pantry) bool {
		return p.Name == name
	})

	if i == -1 {
		return dbtypes.Pantry{}, false
	}

	return db.pantries[i], true
}

func (db *JSON) SetPantry(p dbtypes.Pantry) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if p.Name == "" {
		p.Name = "default"
	}

	i := slices.IndexFunc(db.pantries, func(entry dbtypes.Pantry) bool {
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

func (db *JSON) DeletePantry(name string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	i := slices.IndexFunc(db.pantries, func(p dbtypes.Pantry) bool {
		return p.Name == name
	})

	if i == -1 {
		return fmt.Errorf("pantry %q not found", name)
	}

	db.pantries = append(db.pantries[:i], db.pantries[i+1:]...)

	if err := db.save(); err != nil {
		return err
	}

	return nil
}

func (db *JSON) ShoppingLists() ([]dbtypes.ShoppingList, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	out := make([]dbtypes.ShoppingList, len(db.shoppingLists))
	copy(out, db.shoppingLists)
	return out, nil
}

func (db *JSON) LookupShoppingList(menu, pantry string) (dbtypes.ShoppingList, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	i := slices.IndexFunc(db.shoppingLists, func(p dbtypes.ShoppingList) bool {
		return p.Menu == menu && p.Pantry == pantry
	})

	if i == -1 {
		return dbtypes.ShoppingList{}, false
	}

	return db.shoppingLists[i], true
}

func (db *JSON) SetShoppingList(p dbtypes.ShoppingList) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if p.Menu == "" {
		p.Menu = "default"
	}

	if p.Pantry == "" {
		p.Pantry = "default"
	}

	i := slices.IndexFunc(db.shoppingLists, func(entry dbtypes.ShoppingList) bool {
		return entry.Menu == p.Menu && entry.Pantry == p.Pantry
	})

	if i == -1 {
		db.shoppingLists = append(db.shoppingLists, p)
	} else {
		db.shoppingLists[i] = p
	}

	if err := db.save(); err != nil {
		return err
	}
	return nil
}

func (db *JSON) DeleteShoppingList(menu, pantry string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	i := slices.IndexFunc(db.shoppingLists, func(p dbtypes.ShoppingList) bool {
		return p.Menu == menu && p.Pantry == pantry
	})

	if i == -1 {
		return fmt.Errorf("shopping list (%s, %s) not found", menu, pantry)
	}

	db.shoppingLists = append(db.shoppingLists[:i], db.shoppingLists[i+1:]...)

	if err := db.save(); err != nil {
		return err
	}

	return nil
}

func load(path string, ptr interface{}) error {
	out, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return fmt.Errorf("JSON database: %v", err)
	}

	if len(out) == 0 {
		return nil
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
		save(db.log, db.shoppingListsPath, db.shoppingLists),
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

	err := os.Rename(b.path, b.tmp)
	if err == nil {
		return b, nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return b, err
	}

	if err := os.MkdirAll(filepath.Dir(b.path), 0700); err != nil {
		return b, fmt.Errorf("could not create directory: %v", err)
	}

	if err := os.WriteFile(b.path, nil, 0600); err != nil {
		return b, fmt.Errorf("could not create file: %v", err)
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
