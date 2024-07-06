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
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
)

type JSON struct {
	products      []product.Product
	recipes       []recipe.Recipe
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

func (db *JSON) Recipes() ([]recipe.Recipe, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	out := make([]recipe.Recipe, len(db.recipes))
	copy(out, db.recipes)
	return out, nil
}

func (db *JSON) LookupRecipe(id recipe.ID) (recipe.Recipe, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	i := slices.IndexFunc(db.recipes, func(p recipe.Recipe) bool {
		return p.ID == id
	})

	if i == -1 {
		return recipe.Recipe{}, fs.ErrNotExist
	}

	return db.recipes[i], nil
}

func (db *JSON) SetRecipe(r recipe.Recipe) (recipe.ID, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for r.ID == 0 {
		newID := recipe.NewRandomID()
		idx := slices.IndexFunc(db.recipes, func(entry recipe.Recipe) bool { return entry.ID == newID })
		if idx == -1 {
			r.ID = newID
		}
	}

	i := slices.IndexFunc(db.recipes, func(entry recipe.Recipe) bool { return entry.ID == r.ID })
	if i == -1 {
		db.recipes = append(db.recipes, r)
	} else {
		db.recipes[i] = r
	}

	if err := db.save(); err != nil {
		return 0, err
	}

	return r.ID, nil
}

func (db *JSON) DeleteRecipe(id recipe.ID) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	i := slices.IndexFunc(db.recipes, func(p recipe.Recipe) bool {
		return p.ID == id
	})

	if i == -1 {
		return fmt.Errorf("recipe %d not found", id)
	}

	db.recipes = append(db.recipes[:i], db.recipes[i+1:]...)

	if err := db.save(); err != nil {
		return err
	}

	return nil
}

func (db *JSON) Menus(user string) ([]dbtypes.Menu, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	out := make([]dbtypes.Menu, 0, 1)
	for _, m := range db.menus {
		if m.User == user {
			out = append(out, m)
		}
	}

	return out, nil
}

func (db *JSON) LookupMenu(user, name string) (dbtypes.Menu, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	i := slices.IndexFunc(db.menus, func(p dbtypes.Menu) bool {
		return p.User == user && p.Name == name
	})

	if i == -1 {
		return dbtypes.Menu{}, fs.ErrNotExist
	}

	return db.menus[i], nil
}

func (db *JSON) SetMenu(m dbtypes.Menu) error {
	if m.User == "" {
		return errors.New("user cannot be empty")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	i := slices.IndexFunc(db.menus, func(entry dbtypes.Menu) bool {
		return entry.User == m.User && entry.Name == m.Name
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

func (db *JSON) DeleteMenu(user, name string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	i := slices.IndexFunc(db.menus, func(p dbtypes.Menu) bool {
		return p.User == user && p.Name == name
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

func (db *JSON) Pantries(user string) ([]dbtypes.Pantry, error) {
	if user == "" {
		return nil, errors.New("user cannot be empty")
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	out := make([]dbtypes.Pantry, 0, 1)
	for _, p := range db.pantries {
		if p.User == user {
			out = append(out, p)
		}
	}

	return out, nil
}

func (db *JSON) LookupPantry(user, name string) (dbtypes.Pantry, error) {
	if user == "" {
		return dbtypes.Pantry{}, errors.New("user cannot be empty")
	} else if name == "" {
		return dbtypes.Pantry{}, errors.New("name cannot be empty")
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	i := slices.IndexFunc(db.pantries, func(p dbtypes.Pantry) bool {
		return p.User == user && p.Name == name
	})

	if i == -1 {
		return dbtypes.Pantry{}, fs.ErrNotExist
	}

	return db.pantries[i], nil
}

func (db *JSON) SetPantry(p dbtypes.Pantry) error {
	if p.User == "" {
		return errors.New("user cannot be empty")
	} else if p.Name == "" {
		return errors.New("name cannot be empty")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if p.Name == "" {
		p.Name = "default"
	}

	i := slices.IndexFunc(db.pantries, func(entry dbtypes.Pantry) bool {
		return entry.User == p.User && entry.Name == p.Name
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

func (db *JSON) DeletePantry(user, name string) error {
	if user == "" {
		return errors.New("user cannot be empty")
	} else if name == "" {
		return errors.New("name cannot be empty")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	i := slices.IndexFunc(db.pantries, func(p dbtypes.Pantry) bool {
		return p.User == user && p.Name == name
	})

	if i == -1 {
		return nil
	}

	db.pantries = append(db.pantries[:i], db.pantries[i+1:]...)

	if err := db.save(); err != nil {
		return err
	}

	return nil
}

func (db *JSON) ShoppingLists(user string) ([]dbtypes.ShoppingList, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	out := make([]dbtypes.ShoppingList, 0, 1)
	for _, p := range db.shoppingLists {
		if p.User == user {
			out = append(out, p)
		}
	}

	return out, nil
}

func (db *JSON) LookupShoppingList(user, menu, pantry string) (dbtypes.ShoppingList, error) {
	if user == "" {
		return dbtypes.ShoppingList{}, errors.New("user cannot be empty")
	} else if menu == "" {
		return dbtypes.ShoppingList{}, errors.New("menu cannot be empty")
	} else if pantry == "" {
		return dbtypes.ShoppingList{}, errors.New("pantry cannot be empty")
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	i := slices.IndexFunc(db.shoppingLists, func(p dbtypes.ShoppingList) bool {
		return p.User == user && p.Menu == menu && p.Pantry == pantry
	})

	if i == -1 {
		return dbtypes.ShoppingList{}, fs.ErrNotExist
	}

	return db.shoppingLists[i], nil
}

func (db *JSON) SetShoppingList(p dbtypes.ShoppingList) error {
	if p.User == "" {
		return errors.New("user cannot be empty")
	} else if p.Menu == "" {
		return errors.New("menu cannot be empty")
	} else if p.Pantry == "" {
		return errors.New("pantry cannot be empty")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	i := slices.IndexFunc(db.shoppingLists, func(entry dbtypes.ShoppingList) bool {
		return entry.User == p.User && entry.Menu == p.Menu && entry.Pantry == p.Pantry
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

func (db *JSON) DeleteShoppingList(user, menu, pantry string) error {
	if user == "" {
		return errors.New("user cannot be empty")
	} else if menu == "" {
		return errors.New("menu cannot be empty")
	} else if pantry == "" {
		return errors.New("pantry cannot be empty")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	i := slices.IndexFunc(db.shoppingLists, func(p dbtypes.ShoppingList) bool {
		return p.User == user && p.Menu == menu && p.Pantry == pantry
	})

	if i == -1 {
		return nil
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
