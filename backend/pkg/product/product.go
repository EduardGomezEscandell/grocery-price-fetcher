package product

import (
	"context"
	"fmt"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
)

type Product struct {
	Name      string
	BatchSize float32
	Price     float32

	provider  providers.Provider
	productID providers.ProductID
}

func (p *Product) FetchPrice(ctx context.Context) error {
	price, err := p.provider.FetchPrice(ctx, p.productID)
	if err != nil {
		return fmt.Errorf("could not get price for %s: %w", p.Name, err)
	}

	p.Price = price
	return nil
}
