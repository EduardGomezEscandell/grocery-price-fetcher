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

func (Provider) FetchPrice(ctx context.Context, pid providers.ProductCode) (float32, error) {
	if pid[0] == "" {
		return 0, nil
	}

	price, err := strconv.ParseFloat(pid[0], 32)
	if err != nil {
		return 0, fmt.Errorf("could not parse price: %v", err)
	}

	return float32(price), nil
}

func (Provider) ValidateCode(code providers.ProductCode) error {
	if code[1] != "" {
		return fmt.Errorf("unexpected field at index 1: %q", code[2])
	}

	if code[2] != "" {
		return fmt.Errorf("unexpected field at index 2: %q", code[2])
	}

	return nil
}
