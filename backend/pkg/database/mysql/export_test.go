package mysql

import (
	"context"
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/stretchr/testify/require"
)

//nolint:revive // Testing T goes first
func ClearDB(t *testing.T, ctx context.Context, log logger.Logger, options Settings) {
	t.Helper()

	db, err := New(ctx, log, options)
	require.NoError(t, err)
	defer db.Close()

	tx, err := db.db.BeginTx(db.ctx, nil)
	require.NoError(t, err, "could not begin transaction: %v", err)
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	require.NoError(t, db.clearShoppingLists(tx), "could not clear ingredients")
	require.NoError(t, db.clearPantries(tx), "could not clear pantries")
	require.NoError(t, db.clearMenus(tx), "could not clear menus")
	require.NoError(t, db.clearRecipes(tx), "could not clear recipes")
	require.NoError(t, db.clearProducts(tx), "could not clear products")
	require.NoError(t, db.clearUsers(tx), "could not clear users")

	require.NoError(t, tx.Commit(), "could not commit transaction")
}
