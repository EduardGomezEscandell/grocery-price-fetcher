package mysql

import (
	"context"
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/stretchr/testify/require"
)

//nolint:revive // Testing T goes first
func ClearDB(t *testing.T, ctx context.Context, log logger.Logger, options map[string]any) {
	t.Helper()

	db, err := New(ctx, log, options)
	require.NoError(t, err)
	defer db.Close()

	tx, err := db.db.BeginTx(db.ctx, nil)
	require.NoError(t, err, "could not begin transaction: %v", err)
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	require.NoError(t, db.clearProducts(tx), "could not clear products")

	require.NoError(t, tx.Commit(), "could not commit transaction")
}
