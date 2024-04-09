package mercadona_test

import (
	"context"
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/providers/mercadona"
	"github.com/stretchr/testify/require"
)

func TestMercadonaBadTSV(t *testing.T) {
	t.Parallel()

	p := mercadona.New()

	err := p.UnmarshalTSV()
	require.Error(t, err, "expected an error when unmarshalling")
}

func TestMercadonaBadID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	p := mercadona.New()

	err := p.UnmarshalTSV("1", "0", "bcn1")
	require.NoError(t, err, "expected no error when unmarshalling")

	_, err = p.FetchPrice(ctx)
	require.Error(t, err, "Product with ID 0 should not be found")
}

func TestMercadonaGoodID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	p := mercadona.New()

	err := p.UnmarshalTSV("1", "8713", "bcn1")
	require.NoError(t, err, "expected no error when unmarshalling")

	price, err := p.FetchPrice(ctx)
	require.NoError(t, err, "Product with ID 8713 should be found")
	require.Greater(t, price, float32(0), "expected price to be greater than 0")
}

func TestMercadonaMap(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	p := mercadona.New()

	err := p.UnmarshalMap(map[string]string{
		"batch_size": "7",
		"zone_code":  "bcn1",
		"id":         "8713",
	})
	require.NoError(t, err, "expected no error when unmarshalling")

	price, err := p.FetchPrice(ctx)
	require.NoError(t, err, "Product with ID 8713 should be found")
	require.Greater(t, price, float32(0), "expected price to be greater than 0")
}
