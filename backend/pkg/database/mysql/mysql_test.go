package mysql_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/mysql"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/blank"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/testutils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
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

func TestDBProducts(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := testutils.NewLogger(t)
	log.SetLevel(int(logrus.DebugLevel))

	options := mysql.DefaultSettings()
	mysql.ClearDB(t, ctx, log, options)

	db, err := mysql.New(ctx, log, options)
	require.NoError(t, err)
	defer db.Close()

	products, err := db.Products()
	require.NoError(t, err)
	require.Empty(t, products)

	p := product.Product{
		Name:      "test",
		BatchSize: 1.163,
		Price:     111.84,
		Provider:  blank.Provider{},
		ProductID: [3]string{"1"},
	}

	_, ok := db.LookupProduct(p.Name)
	require.False(t, ok)

	err = db.SetProduct(p)
	require.NoError(t, err)

	got, ok := db.LookupProduct(p.Name)
	require.True(t, ok)
	require.Equal(t, p, got)

	products, err = db.Products()
	require.NoError(t, err)
	require.ElementsMatch(t, []product.Product{p}, products)

	err = db.DeleteProduct(p.Name)
	require.NoError(t, err)

	products, err = db.Products()
	require.NoError(t, err)
	require.Empty(t, products)
}

//nolint:dupl // This is a test file, duplication is expected
func TestDBRecipes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := testutils.NewLogger(t)
	log.SetLevel(int(logrus.DebugLevel))

	options := mysql.DefaultSettings()
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

	err = db.SetProduct(p[0])
	require.NoError(t, err)

	err = db.SetProduct(p[1])
	require.NoError(t, err)

	rec := types.Recipe{
		Name: "Water",
		Ingredients: []types.Ingredient{
			{Name: "Hydrogen", Amount: 2},
			{Name: "Oxygen", Amount: 1},
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
	require.ElementsMatch(t, []types.Recipe{rec}, recs)

	rec.Ingredients[0].Amount = 3
	err = db.SetRecipe(rec)
	require.NoError(t, err)

	recs, err = db.Recipes()
	require.NoError(t, err)
	require.ElementsMatch(t, []types.Recipe{rec}, recs)

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

	r := []types.Recipe{
		{
			Name: "Water",
			Ingredients: []types.Ingredient{
				{Name: "Hydrogen", Amount: 2},
				{Name: "Oxygen", Amount: 1},
			},
		},
		{
			Name: "Hydrogen Peroxide",
			Ingredients: []types.Ingredient{
				{Name: "Hydrogen", Amount: 2},
				{Name: "Oxygen", Amount: 2},
			},
		},
		{
			Name: "Oxygen Gas",
			Ingredients: []types.Ingredient{
				{Name: "Oxygen", Amount: 2},
			},
		},
	}

	for _, product := range p {
		err := db.SetProduct(product)
		require.NoErrorf(t, err, "could not set product %s", product.Name)
	}

	for _, recipe := range r {
		err := db.SetRecipe(recipe)
		require.NoErrorf(t, err, "could not set recipe %s", recipe.Name)
	}

	menus, err := db.Menus()
	require.NoError(t, err)
	require.Empty(t, menus)

	m := types.Menu{
		Name: "Test Menu",
		Days: []types.Day{
			{
				Name: "Monday",
				Meals: []types.Meal{
					{Name: "Lunch", Dishes: []types.Dish{
						{Name: "Water", Amount: 1.12},
					}},
					{Name: "Dinner", Dishes: []types.Dish{
						{Name: "Hydrogen Peroxide", Amount: 3},
						{Name: "Oxygen Gas", Amount: 4}}},
				},
			},
			{
				Name: "Saturday",
				Meals: []types.Meal{
					{Name: "Lunch", Dishes: []types.Dish{
						{Name: "Water", Amount: 1}}},
					{Name: "Dinner", Dishes: []types.Dish{
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
	require.ElementsMatch(t, []types.Menu{m}, menus)

	m.Days[0].Meals[0].Dishes[0].Amount = 2.34
	err = db.SetMenu(m)
	require.NoError(t, err)

	menus, err = db.Menus()
	require.NoError(t, err)
	require.ElementsMatch(t, []types.Menu{m}, menus)

	err = db.DeleteMenu(m.Name)
	require.NoError(t, err)

	menus, err = db.Menus()
	require.NoError(t, err)
	require.Empty(t, menus)
}

//nolint:dupl // This is a test file, duplication is expected
func TestDBPantries(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := testutils.NewLogger(t)
	log.SetLevel(int(logrus.DebugLevel))

	options := mysql.DefaultSettings()
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

	err = db.SetProduct(p[0])
	require.NoError(t, err)

	err = db.SetProduct(p[1])
	require.NoError(t, err)

	pantry := types.Pantry{
		Name: "Test Pantry",
		Contents: []types.Ingredient{
			{Name: "Hydrogen", Amount: 2165},
			{Name: "Oxygen", Amount: 100},
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
	require.ElementsMatch(t, []types.Pantry{pantry}, pantries)

	pantry.Contents[0].Amount = 1
	err = db.SetPantry(pantry)
	require.NoError(t, err)

	pantries, err = db.Pantries()
	require.NoError(t, err)
	require.ElementsMatch(t, []types.Pantry{pantry}, pantries)

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

	err = db.SetProduct(p[0])
	require.NoError(t, err)

	err = db.SetProduct(p[1])
	require.NoError(t, err)

	sl := types.ShoppingList{
		Name:      "Test Shopping List",
		TimeStamp: "2021-01-01T00:00:00Z",
		Items: []string{
			"Hydrogen",
		},
	}

	pantries, err := db.ShoppingLists()
	require.NoError(t, err)
	require.Empty(t, pantries)

	_, ok := db.LookupShoppingList(sl.Name)
	require.False(t, ok)

	err = db.SetShoppingList(sl)
	require.NoError(t, err)

	got, ok := db.LookupShoppingList(sl.Name)
	require.True(t, ok)
	require.Equal(t, sl, got)

	pantries, err = db.ShoppingLists()
	require.NoError(t, err)
	require.ElementsMatch(t, []types.ShoppingList{sl}, pantries)

	sl.Items = append(sl.Items, "Oxygen")
	err = db.SetShoppingList(sl)
	require.NoError(t, err)

	pantries, err = db.ShoppingLists()
	require.NoError(t, err)
	require.ElementsMatch(t, []types.ShoppingList{sl}, pantries)

	err = db.DeleteShoppingList(sl.Name)
	require.NoError(t, err)

	pantries, err = db.ShoppingLists()
	require.NoError(t, err)
	require.Empty(t, pantries)
}
