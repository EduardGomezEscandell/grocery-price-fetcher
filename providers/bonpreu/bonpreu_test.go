package bonpreu_test

import (
	"context"
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/providers/bonpreu"
	"github.com/stretchr/testify/require"
)

func TestBonpreu(t *testing.T) {
	t.Parallel()

	_, err := bonpreu.Get(context.Background(), "Example", "0")
	require.Error(t, err, "Product with ID 0 should not be found")

	price, err := bonpreu.Get(context.Background(), "Example", "90041")
	require.NoError(t, err, "Product with ID 8713 should be found")
	require.Greater(t, price, float32(0), "expected price to be greater than 0")
}
