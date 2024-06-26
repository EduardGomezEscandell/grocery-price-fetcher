package product

import (
	"context"
	"fmt"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
)

type ID uint32

const IDSize = 32

type Product struct {
	ID        ID
	Name      string
	BatchSize float32
	Price     float32

	Provider    providers.Provider
	ProductCode providers.ProductCode
}

func (p *Product) FetchPrice(ctx context.Context) error {
	price, err := p.Provider.FetchPrice(ctx, p.ProductCode)
	if err != nil {
		return fmt.Errorf("could not get price for %s: %w", p.Name, err)
	}

	p.Price = price
	return nil
}
