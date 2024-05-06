package product

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/providers"
	"golang.org/x/exp/maps"
)

type Product struct {
	Name     string
	Price    float32
	Provider providers.Provider
}

func (p *Product) UnmarshalTSV(args []string) (err error) {
	const (
		flieldName = iota
		fieldProvider
		fieldArgv
	)

	if len(args) < 2 {
		p.Name = "COULD NOT PARSE"
		return fmt.Errorf("expected at least 2 fields, got %d", len(args))
	}

	p.Name = args[flieldName]

	p.Provider, err = providers.ParseTSV(args[fieldProvider], args[fieldArgv:])
	if err != nil {
		return fmt.Errorf("could not parse provider for %s: %w", p.Name, err)
	}

	return nil
}

func (p *Product) UnmarshalJSON(b []byte) (err error) {
	var helper struct {
		Name      string
		Providers map[string](map[string]string)
	}

	if err := json.Unmarshal(b, &helper); err != nil {
		return fmt.Errorf("could not unmarshal product: %w", err)
	}

	p.Name = helper.Name
	if len(helper.Providers) != 1 {
		return fmt.Errorf("expected 1 provider, got %d", len(helper.Providers))
	}

	pName := maps.Keys(helper.Providers)[0]
	pArgv := maps.Values(helper.Providers)[0]

	p.Provider, err = providers.ParseMap(pName, pArgv)
	if err != nil {
		return fmt.Errorf("could not parse provider for %s: %w", p.Name, err)
	}

	return nil
}

func (p *Product) FetchPrice(ctx context.Context, log logger.Logger) error {
	log.Tracef("Fetching price for %s", p.Name)
	price, err := p.Provider.FetchPrice(ctx, log)
	if err != nil {
		return fmt.Errorf("could not get price for %s: %w", p.Name, err)
	}
	log.Tracef("Got price for %s", p.Name)

	p.Price = price
	return nil
}
