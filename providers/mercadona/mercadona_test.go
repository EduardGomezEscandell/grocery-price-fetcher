package mercadona_test

import (
	"context"
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/providers/mercadona"
	"github.com/stretchr/testify/require"
)

func TestMercadona(t *testing.T) {
	t.Parallel()

	_, err := mercadona.Get(context.Background(), "Example", "0")
	require.Error(t, err, "Product with a single argument should not be found")

	_, err = mercadona.Get(context.Background(), "Example", "0", "bcn1")
	require.Error(t, err, "Product with ID 0 should not be found")

	price, err := mercadona.Get(context.Background(), "Example", "8713", "bcn1")
	require.NoError(t, err, "Product with ID 8713 should be found")
	require.Greater(t, price, float32(0), "expected price to be greater than 0")
}
