package product

import (
	"context"
	"fmt"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
)

type Product struct {
	ID        uint32
	Name      string
	BatchSize float32
	Price     float32

	Provider  providers.Provider
	ProductID providers.ProductID
}

func (p *Product) FetchPrice(ctx context.Context) error {
	price, err := p.Provider.FetchPrice(ctx, p.ProductID)
	if err != nil {
		return fmt.Errorf("could not get price for %s: %w", p.Name, err)
	}

	p.Price = price
	return nil
}
