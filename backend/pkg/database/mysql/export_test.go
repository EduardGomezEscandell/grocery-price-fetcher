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

	require.NoError(t, db.dropAllTables(), "could not commit transaction")
}
