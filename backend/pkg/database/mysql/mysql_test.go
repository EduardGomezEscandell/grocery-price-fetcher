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
