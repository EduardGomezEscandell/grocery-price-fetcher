package mysql_test

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtestutils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/mysql"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/blank"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/testutils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestMain(m *testing.M) {
	providers.Register(blank.Provider{})

	if os.Getenv("MYSQL_SKIP_TEST_MAIN") == "" {
		fmt.Println("Starting database")
		cmd := exec.Command("make", "stand-up")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal("could not start database")
		}
		fmt.Println("Database started")

		defer func() {
			fmt.Println("Stopping database")
			cmd := exec.Command("make", "stand-down")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Fatal("Could not shut down database")
			}
			fmt.Println("Database stopped")
		}()
	}

	m.Run()
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	input := []byte(`
type: mysql
options:
  user: joe
  passwordfile: /etc/secret
  host: hostboy
  port: 1234
`)

	want := database.Settings{
		Type: "mysql",
		Options: mysql.Settings{
			User:            "joe",
			PasswordFile:    "/etc/secret",
			Host:            "hostboy",
			Port:            "1234",
			ConnectTimeout:  time.Minute,
			ConnectCooldown: 5 * time.Second,
		},
	}

	var got database.Settings
	err := yaml.Unmarshal(input, &got)
	require.NoError(t, err)

	require.Equal(t, want, got)
}

func TestBattery(t *testing.T) {
	testCases := map[string]func(*testing.T, func() database.DB){
		"Products": dbtestutils.ProductsTest,
		"Recipes":  dbtestutils.RecipesTest,
		"Menus":    dbtestutils.MenuTest,
		"Pantries": dbtestutils.PantriesTest,
		"Shopping": dbtestutils.ShoppingListsTest,
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			log := testutils.NewLogger(t)
			log.SetLevel(int(logrus.DebugLevel))

			options := mysql.DefaultSettings()
			options.PasswordFile = "./testdata/db_root_password.txt"
			mysql.ClearDB(t, ctx, log, options)

			test(t, func() database.DB {
				db, err := mysql.New(context.Background(), log, options)
				require.NoError(t, err)
				return db
			})
		})
	}
}

func TestDBProducts(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := testutils.NewLogger(t)
	log.SetLevel(int(logrus.DebugLevel))

	options := mysql.DefaultSettings()
	options.PasswordFile = "./testdata/db_root_password.txt"
	mysql.ClearDB(t, ctx, log, options)

	db, err := mysql.New(ctx, log, options)
	require.NoError(t, err)
	defer db.Close()

	products, err := db.Products()
	require.NoError(t, err)
	require.Empty(t, products)

	p := product.Product{
		Name:        "test",
		BatchSize:   1.163,
		Price:       111.84,
		Provider:    blank.Provider{},
		ProductCode: [3]string{"1"},
	}

	_, err = db.LookupProduct(9999)
	require.ErrorIs(t, err, fs.ErrNotExist)

	id, err := db.SetProduct(p)
	require.NoError(t, err)
	require.NotZero(t, id)
	p.ID = id

	got, err := db.LookupProduct(p.ID)
	require.NoError(t, err)
	require.Equal(t, p, got)

	products, err = db.Products()
	require.NoError(t, err)
	require.ElementsMatch(t, []product.Product{p}, products)

	err = db.DeleteProduct(p.ID)
	require.NoError(t, err)

	products, err = db.Products()
	require.NoError(t, err)
	require.Empty(t, products)
}

func TestDBRecipes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := testutils.NewLogger(t)
	log.SetLevel(int(logrus.DebugLevel))

	options := mysql.DefaultSettings()
	options.PasswordFile = "./testdata/db_root_password.txt"
	mysql.ClearDB(t, ctx, log, options)

	db, err := mysql.New(ctx, log, options)
	require.NoError(t, err)
	defer db.Close()

	p := []product.Product{
		{
			Name:      "Hydrogen",
			BatchSize: 1,
			Price:     1,
			Provider:  blank.Provider{},
		},
		{
			Name:      "Oxygen",
			BatchSize: 16,
			Price:     14,
			Provider:  blank.Provider{},
		},
	}

	id, err := db.SetProduct(p[0])
	require.NoError(t, err)
	p[0].ID = id

	id, err = db.SetProduct(p[1])
	require.NoError(t, err)
	p[1].ID = id

	rec := recipe.Recipe{
		Name: "Water",
		Ingredients: []recipe.Ingredient{
			{ProductID: 1, Amount: 2},
			{ProductID: 2, Amount: 1},
		},
	}

	recs, err := db.Recipes()
	require.NoError(t, err)
	require.Empty(t, recs)

	_, ok := db.LookupRecipe(rec.Name)
	require.False(t, ok)

	err = db.SetRecipe(rec)
	require.NoError(t, err)

	got, ok := db.LookupRecipe(rec.Name)
	require.True(t, ok)
	require.Equal(t, rec, got)

	recs, err = db.Recipes()
	require.NoError(t, err)
	require.ElementsMatch(t, []recipe.Recipe{rec}, recs)

	rec.Ingredients[0].Amount = 3
	err = db.SetRecipe(rec)
	require.NoError(t, err)

	recs, err = db.Recipes()
	require.NoError(t, err)
	require.ElementsMatch(t, []recipe.Recipe{rec}, recs)

	err = db.DeleteRecipe(rec.Name)
	require.NoError(t, err)

	recs, err = db.Recipes()
	require.NoError(t, err)
	require.Empty(t, recs)
}

func TestDBMenus(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := testutils.NewLogger(t)
	log.SetLevel(int(logrus.DebugLevel))

	options := mysql.DefaultSettings()
	options.PasswordFile = "./testdata/db_root_password.txt"
	mysql.ClearDB(t, ctx, log, options)

	db, err := mysql.New(ctx, log, options)
	require.NoError(t, err)
	defer db.Close()

	p := []product.Product{
		{
			ID:        13,
			Name:      "Hydrogen",
			BatchSize: 1,
			Price:     1,
			Provider:  blank.Provider{},
		},
		{
			ID:        55,
			Name:      "Oxygen",
			BatchSize: 16,
			Price:     14,
			Provider:  blank.Provider{},
		},
	}

	r := []recipe.Recipe{
		{
			Name: "Water",
			Ingredients: []recipe.Ingredient{
				{ProductID: 13, Amount: 2},
				{ProductID: 55, Amount: 1},
			},
		},
		{
			Name: "Hydrogen Peroxide",
			Ingredients: []recipe.Ingredient{
				{ProductID: 13, Amount: 2},
				{ProductID: 55, Amount: 2},
			},
		},
		{
			Name: "Oxygen Gas",
			Ingredients: []recipe.Ingredient{
				{ProductID: 55, Amount: 2},
			},
		},
	}

	for _, product := range p {
		_, err := db.SetProduct(product)
		require.NoErrorf(t, err, "could not set product %s", product.Name)
	}

	for _, recipe := range r {
		err := db.SetRecipe(recipe)
		require.NoErrorf(t, err, "could not set recipe %s", recipe.Name)
	}

	menus, err := db.Menus()
	require.NoError(t, err)
	require.Empty(t, menus)

	m := dbtypes.Menu{
		Name: "Test Menu",
		Days: []dbtypes.Day{
			{
				Name: "Monday",
				Meals: []dbtypes.Meal{
					{Name: "Lunch", Dishes: []dbtypes.Dish{
						{Name: "Water", Amount: 1.12},
					}},
					{Name: "Dinner", Dishes: []dbtypes.Dish{
						{Name: "Hydrogen Peroxide", Amount: 3},
						{Name: "Oxygen Gas", Amount: 4}}},
				},
			},
			{
				Name: "Saturday",
				Meals: []dbtypes.Meal{
					{Name: "Lunch", Dishes: []dbtypes.Dish{
						{Name: "Water", Amount: 1}}},
					{Name: "Dinner", Dishes: []dbtypes.Dish{
						{Name: "Hydrogen Peroxide", Amount: 3}}},
				},
			},
		},
	}

	_, ok := db.LookupMenu(m.Name)
	require.False(t, ok)

	err = db.SetMenu(m)
	require.NoError(t, err)

	got, ok := db.LookupMenu(m.Name)
	require.True(t, ok)
	require.Equal(t, m, got)

	menus, err = db.Menus()
	require.NoError(t, err)
	require.ElementsMatch(t, []dbtypes.Menu{m}, menus)

	m.Days[0].Meals[0].Dishes[0].Amount = 2.34
	err = db.SetMenu(m)
	require.NoError(t, err)

	menus, err = db.Menus()
	require.NoError(t, err)
	require.ElementsMatch(t, []dbtypes.Menu{m}, menus)

	err = db.DeleteMenu(m.Name)
	require.NoError(t, err)

	menus, err = db.Menus()
	require.NoError(t, err)
	require.Empty(t, menus)
}

func TestDBPantries(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := testutils.NewLogger(t)
	log.SetLevel(int(logrus.DebugLevel))

	options := mysql.DefaultSettings()
	options.PasswordFile = "./testdata/db_root_password.txt"
	mysql.ClearDB(t, ctx, log, options)

	db, err := mysql.New(ctx, log, options)
	require.NoError(t, err)
	defer db.Close()

	p := []product.Product{
		{
			ID:        13,
			Name:      "Hydrogen",
			BatchSize: 1,
			Price:     1,
			Provider:  blank.Provider{},
		},
		{
			ID:        55,
			Name:      "Oxygen",
			BatchSize: 16,
			Price:     14,
			Provider:  blank.Provider{},
		},
	}

	_, err = db.SetProduct(p[0])
	require.NoError(t, err)

	_, err = db.SetProduct(p[1])
	require.NoError(t, err)

	pantry := dbtypes.Pantry{
		Name: "Test Pantry",
		Contents: []recipe.Ingredient{
			{ProductID: 13, Amount: 2165},
			{ProductID: 55, Amount: 100},
		},
	}

	pantries, err := db.Pantries()
	require.NoError(t, err)
	require.Empty(t, pantries)

	_, ok := db.LookupPantry(pantry.Name)
	require.False(t, ok)

	err = db.SetPantry(pantry)
	require.NoError(t, err)

	got, ok := db.LookupPantry(pantry.Name)
	require.True(t, ok)
	require.Equal(t, pantry, got)

	pantries, err = db.Pantries()
	require.NoError(t, err)
	require.ElementsMatch(t, []dbtypes.Pantry{pantry}, pantries)

	pantry.Contents[0].Amount = 1
	err = db.SetPantry(pantry)
	require.NoError(t, err)

	pantries, err = db.Pantries()
	require.NoError(t, err)
	require.ElementsMatch(t, []dbtypes.Pantry{pantry}, pantries)

	err = db.DeletePantry(pantry.Name)
	require.NoError(t, err)

	pantries, err = db.Pantries()
	require.NoError(t, err)
	require.Empty(t, pantries)
}

func TestDBShoopingLists(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := testutils.NewLogger(t)
	log.SetLevel(int(logrus.DebugLevel))

	options := mysql.DefaultSettings()
	options.PasswordFile = "./testdata/db_root_password.txt"
	mysql.ClearDB(t, ctx, log, options)

	db, err := mysql.New(ctx, log, options)
	require.NoError(t, err)
	defer db.Close()

	p := []product.Product{
		{
			ID:        15,
			Name:      "Hydrogen",
			BatchSize: 1,
			Price:     1,
			Provider:  blank.Provider{},
		},
		{
			ID:        99,
			Name:      "Oxygen",
			BatchSize: 16,
			Price:     14,
			Provider:  blank.Provider{},
		},
	}

	_, err = db.SetProduct(p[0])
	require.NoError(t, err)

	_, err = db.SetProduct(p[1])
	require.NoError(t, err)

	sl := dbtypes.ShoppingList{
		Menu:     "My test menu",
		Pantry:   "My test pantry",
		Contents: []product.ID{1},
	}

	pantries, err := db.ShoppingLists()
	require.NoError(t, err)
	require.Empty(t, pantries)

	_, ok := db.LookupShoppingList(sl.Menu, sl.Pantry)
	require.False(t, ok)

	err = db.SetShoppingList(sl)
	require.NoError(t, err)

	got, ok := db.LookupShoppingList(sl.Menu, sl.Pantry)
	require.True(t, ok)
	require.Equal(t, sl, got)

	pantries, err = db.ShoppingLists()
	require.NoError(t, err)
	require.ElementsMatch(t, []dbtypes.ShoppingList{sl}, pantries)

	sl.Contents = append(sl.Contents, 99)
	err = db.SetShoppingList(sl)
	require.NoError(t, err)

	pantries, err = db.ShoppingLists()
	require.NoError(t, err)
	require.ElementsMatch(t, []dbtypes.ShoppingList{sl}, pantries)

	err = db.DeleteShoppingList(sl.Menu, sl.Pantry)
	require.NoError(t, err)

	pantries, err = db.ShoppingLists()
	require.NoError(t, err)
	require.Empty(t, pantries)
}
