package mercadona_test

import (
	"context"
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/mercadona"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestMercadonaBadID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	p := mercadona.New(testLogger())

	_, err := p.FetchPrice(ctx, providers.ProductID{"0", "bcn1"})
	require.Error(t, err, "Product with ID 0 should not be found")
}

func TestMercadonaGoodID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	p := mercadona.New(testLogger())

	price, err := p.FetchPrice(ctx, providers.ProductID{"3852", "bcn1"})
	require.NoError(t, err, "Product with ID 3852 should be found")
	require.Greater(t, price, float32(0), "expected price to be greater than 0")
}

func testLogger() logger.Logger {
	l := logger.New()
	l.SetLevel(int(logrus.TraceLevel))
	return l
}
