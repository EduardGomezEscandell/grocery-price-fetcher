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
  password-file: /etc/secret
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

func TestMySQL(t *testing.T) {
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
			log.SetLevel(int(logrus.TraceLevel))

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

func TestMySQLProducts(t *testing.T) {
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

func TestMySQLRecipes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := testutils.NewLogger(t)
	log.SetLevel(int(logrus.TraceLevel))

	options := mysql.DefaultSettings()
	options.PasswordFile = "./testdata/db_root_password.txt"
	mysql.ClearDB(t, ctx, log, options)

	db, err := mysql.New(ctx, log, options)
	require.NoError(t, err)
	defer db.Close()

	hydrogen := product.Product{
		Name:      "Hydrogen",
		BatchSize: 1,
		Price:     1,
		Provider:  blank.Provider{},
	}

	oxygen := product.Product{
		Name:      "Oxygen",
		BatchSize: 16,
		Price:     14,
		Provider:  blank.Provider{},
	}

	pID, err := db.SetProduct(hydrogen)
	require.NoError(t, err)
	hydrogen.ID = pID

	pID, err = db.SetProduct(oxygen)
	require.NoError(t, err)
	oxygen.ID = pID

	rec := recipe.Recipe{
		Name: "Water",
		Ingredients: []recipe.Ingredient{
			{ProductID: hydrogen.ID, Amount: 2},
			{ProductID: oxygen.ID, Amount: 1},
		},
	}

	recs, err := db.Recipes()
	require.NoError(t, err)
	require.Empty(t, recs)

	_, err = db.LookupRecipe(555)
	require.ErrorIs(t, err, fs.ErrNotExist)

	rID, err := db.SetRecipe(rec)
	require.NoError(t, err)
	require.NotZero(t, rID, "expected non-zero ID")
	rec.ID = rID

	got, err := db.LookupRecipe(rec.ID)
	require.NoError(t, err)
	requireSameRecipe(t, rec, got)

	recs, err = db.Recipes()
	require.NoError(t, err)
	require.Len(t, recs, 1)
	requireSameRecipe(t, rec, recs[0])

	t.Logf("Updating recipe to %+#v", rec)

	rec.Ingredients[0].Amount = 3
	_, err = db.SetRecipe(rec)
	require.NoError(t, err)

	recs, err = db.Recipes()
	require.NoError(t, err)
	require.Len(t, recs, 1)
	requireSameRecipe(t, rec, recs[0])

	err = db.DeleteRecipe(rec.ID)
	require.NoError(t, err)

	recs, err = db.Recipes()
	require.NoError(t, err)
	require.Empty(t, recs)
}

func requireSameRecipe(t *testing.T, want, got recipe.Recipe) {
	t.Helper()
	require.Equal(t, want.Name, got.Name)
	require.Equal(t, want.ID, got.ID)
	require.ElementsMatch(t, want.Ingredients, got.Ingredients)
}

func TestMySQLMenus(t *testing.T) {
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

	user := "user_id_123"
	
	H := product.Product{
		Name:      "Hydrogen",
		BatchSize: 1,
		Price:     1,
		Provider:  blank.Provider{},
	}
	O := product.Product{
		Name:      "Oxygen",
		BatchSize: 16,
		Price:     14,
		Provider:  blank.Provider{},
	}

	id, err := db.SetProduct(H)
	require.NoErrorf(t, err, "could not set product %s", H.Name)
	H.ID = id

	id, err = db.SetProduct(O)
	require.NoErrorf(t, err, "could not set product %s", O.Name)
	O.ID = id

	H2O := recipe.Recipe{
		Name: "Water",
		Ingredients: []recipe.Ingredient{
			{ProductID: H.ID, Amount: 2},
			{ProductID: O.ID, Amount: 1},
		},
	}
	H2O2 := recipe.Recipe{
		Name: "Hydrogen Peroxide",
		Ingredients: []recipe.Ingredient{
			{ProductID: H.ID, Amount: 2},
			{ProductID: O.ID, Amount: 2},
		},
	}
	O2 := recipe.Recipe{
		Name: "Oxygen Gas",
		Ingredients: []recipe.Ingredient{
			{ProductID: O.ID, Amount: 2},
		},
	}

	rID, err := db.SetRecipe(H2O)
	require.NoErrorf(t, err, "could not set recipe %s", H2O.Name)
	H2O.ID = rID

	rID, err = db.SetRecipe(H2O2)
	require.NoErrorf(t, err, "could not set recipe %s", H2O2.Name)
	H2O2.ID = rID

	rID, err = db.SetRecipe(O2)
	require.NoErrorf(t, err, "could not set recipe %s", O2.Name)
	O2.ID = rID

	menus, err := db.Menus(user)
	require.NoError(t, err)
	require.Empty(t, menus)

	m := dbtypes.Menu{
		User: user,
		Name: "Test Menu",
		Days: []dbtypes.Day{
			{
				Name: "Monday",
				Meals: []dbtypes.Meal{
					{Name: "Lunch", Dishes: []dbtypes.Dish{
						{ID: H2O.ID, Amount: 1.12},
					}},
					{Name: "Dinner", Dishes: []dbtypes.Dish{
						{ID: H2O2.ID, Amount: 3},
						{ID: O2.ID, Amount: 4}}},
				},
			},
			{
				Name: "Saturday",
				Meals: []dbtypes.Meal{
					{Name: "Lunch", Dishes: []dbtypes.Dish{
						{ID: H2O.ID, Amount: 1}}},
					{Name: "Dinner", Dishes: []dbtypes.Dish{
						{ID: H2O2.ID, Amount: 3}}},
				},
			},
		},
	}

	_, err = db.LookupMenu(user, m.Name)
	require.Error(t, err)

	err = db.SetMenu(m)
	require.NoError(t, err)

	got, err := db.LookupMenu(user, m.Name)
	require.NoError(t, err)
	require.Equal(t, m, got)

	menus, err = db.Menus(user)
	require.NoError(t, err)
	require.ElementsMatch(t, []dbtypes.Menu{m}, menus)

	m.Days[0].Meals[0].Dishes[0].Amount = 2.34
	err = db.SetMenu(m)
	require.NoError(t, err)

	menus, err = db.Menus(user)
	require.NoError(t, err)
	require.ElementsMatch(t, []dbtypes.Menu{m}, menus)

	err = db.DeleteMenu(user, m.Name)
	require.NoError(t, err)

	menus, err = db.Menus(user)
	require.NoError(t, err)
	require.Empty(t, menus)
}

func TestMySQLPantries(t *testing.T) {
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

	pantry := dbtypes.Pantry{
		Name: "Test Pantry",
		Contents: []recipe.Ingredient{
			{ProductID: p[0].ID, Amount: 2165},
			{ProductID: p[1].ID, Amount: 100},
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
	requireSamePantry(t, pantry, got)

	pantries, err = db.Pantries()
	require.NoError(t, err)
	require.Len(t, pantries, 1)
	requireSamePantry(t, pantry, pantries[0])

	pantry.Contents[0].Amount = 1
	err = db.SetPantry(pantry)
	require.NoError(t, err)

	pantries, err = db.Pantries()
	require.NoError(t, err)
	require.Len(t, pantries, 1)
	requireSamePantry(t, pantry, pantries[0])

	err = db.DeletePantry(pantry.Name)
	require.NoError(t, err)

	pantries, err = db.Pantries()
	require.NoError(t, err)
	require.Empty(t, pantries)
}

func requireSamePantry(t *testing.T, want, got dbtypes.Pantry) {
	t.Helper()
	require.Equal(t, want.Name, got.Name)
	require.ElementsMatch(t, want.Contents, got.Contents)
}

func TestMySQLShoopingLists(t *testing.T) {
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

	err = db.SetMenu(dbtypes.Menu{Name: "My test menu"})
	require.NoError(t, err)

	err = db.SetPantry(dbtypes.Pantry{Name: "My test pantry"})
	require.NoError(t, err)

	sl := dbtypes.ShoppingList{
		Menu:     "My test menu",
		Pantry:   "My test pantry",
		Contents: []product.ID{p[0].ID},
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
	requireSameShoppingList(t, sl, got)

	pantries, err = db.ShoppingLists()
	require.NoError(t, err)
	require.Len(t, pantries, 1)
	requireSameShoppingList(t, sl, pantries[0])

	sl.Contents = append(sl.Contents, p[1].ID)
	err = db.SetShoppingList(sl)
	require.NoError(t, err)

	pantries, err = db.ShoppingLists()
	require.NoError(t, err)
	require.Len(t, pantries, 1)
	requireSameShoppingList(t, sl, pantries[0])

	err = db.DeleteShoppingList(sl.Menu, sl.Pantry)
	require.NoError(t, err)

	pantries, err = db.ShoppingLists()
	require.NoError(t, err)
	require.Empty(t, pantries)
}

func requireSameShoppingList(t *testing.T, want, got dbtypes.ShoppingList) {
	t.Helper()
	require.Equal(t, want.Menu, got.Menu)
	require.Equal(t, want.Pantry, got.Pantry)
	require.ElementsMatch(t, want.Contents, got.Contents)
}
