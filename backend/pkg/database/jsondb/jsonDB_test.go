package jsondb_test

import (
	"context"
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtestutils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/jsondb"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/blank"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	providers.Register(blank.Provider{})
	m.Run()
}

func TestJSONDB(t *testing.T) {
	t.Parallel()

	testCases := map[string]func(*testing.T, func() database.DB){
		"Products": dbtestutils.ProductsTest,
		"Recipes":  dbtestutils.RecipesTest,
		"Menus":    dbtestutils.MenuTest,
		"Pantries": dbtestutils.PantriesTest,
		"Shopping": dbtestutils.ShoppingListsTest,
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			log := testutils.NewLogger(t)
			dir := t.TempDir()

			test(t, func() database.DB {
				db, err := jsondb.New(context.Background(), log, jsondb.DefaultSettingsPath(dir))
				require.NoError(t, err)
				return db
			})
		})
	}
}

func TestMenu(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dir := t.TempDir()

	dbtestutils.MenuTest(t, func() database.DB {
		db, err := jsondb.New(ctx, testutils.NewLogger(t), jsondb.DefaultSettingsPath(dir))
		require.NoError(t, err)
		return db
	})
}
