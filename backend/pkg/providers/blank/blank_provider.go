package blank

import (
	"context"
	"fmt"
	"strconv"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
)

type Provider struct{}

func (Provider) Name() string {
	return "NoProvider"
}

func (Provider) FetchPrice(ctx context.Context, pid providers.ProductID) (float32, error) {
	switch len(pid) {
	case 0:
		return 0, nil
	case 1:
		price, err := strconv.ParseFloat(pid[0], 32)
		if err != nil {
			return 0, fmt.Errorf("could not parse price: %v", err)
		}
		return float32(price), nil
	default:
		return 0, fmt.Errorf("expected 0 or 1 field in product ID, got %d", len(pid))
	}
}

func (Provider) ValidateID(pid providers.ProductID) error {
	if len(pid) > 1 {
		return fmt.Errorf("expected 0 or 1 field in product ID, got %d", len(pid))
	}
	return nil
}
