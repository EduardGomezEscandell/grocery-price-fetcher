package bonpreu_test

import (
	"context"
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/providers/bonpreu"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestBonpreuBadID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	p := bonpreu.New(testLogger())

	_, err := p.FetchPrice(ctx, providers.ProductID{"0"})
	require.Error(t, err, "Product with ID 0 should not be found")
}

func TestBonpreuGoodID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	p := bonpreu.New(testLogger())

	price, err := p.FetchPrice(ctx, providers.ProductID{"90041"})
	require.NoError(t, err, "Product with ID 90041 should be found")
	require.Greater(t, price, float32(0), "expected price to be greater than 0")
}

func testLogger() logger.Logger {
	l := logger.New()
	l.SetLevel(int(logrus.TraceLevel))
	return l
}
