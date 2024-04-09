package bonpreu_test

import (
	"context"
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/providers/bonpreu"
	"github.com/stretchr/testify/require"
)

func TestBonpreuBadTSV(t *testing.T) {
	t.Parallel()

	p := bonpreu.New()

	err := p.UnmarshalTSV()
	require.Error(t, err, "expected no error when unmarshalling")
}

func TestBonpreuBadID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	p := bonpreu.New()

	err := p.UnmarshalTSV("1", "0")
	require.NoError(t, err, "expected no error when unmarshalling")

	_, err = p.FetchPrice(ctx)
	require.Error(t, err, "Product with ID 0 should not be found")
}

func TestBonpreuGoodID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	p := bonpreu.New()

	err := p.UnmarshalTSV("1", "90041")
	require.NoError(t, err, "expected no error when unmarshalling")

	price, err := p.FetchPrice(ctx)
	require.NoError(t, err, "Product with ID 8713 should be found")
	require.Greater(t, price, float32(0), "expected price to be greater than 0")
}

func TestBonpreuMap(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	p := bonpreu.New()

	err := p.UnmarshalMap(map[string]string{
		"batch_size": "7",
		"id":         "90041",
	})
	require.NoError(t, err, "expected no error when unmarshalling")

	price, err := p.FetchPrice(ctx)
	require.NoError(t, err, "Product with ID 8713 should be found")
	require.Greater(t, price, float32(0), "expected price to be greater than 0")
}
