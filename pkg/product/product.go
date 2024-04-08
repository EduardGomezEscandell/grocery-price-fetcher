package product

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Product struct {
	Name     string
	Provider string
	Count    float32
	Price    float32

	providerArgs []string
}

const (
	flieldName = iota
	fieldProvider
	fieldCount
	fieldArgv
)

func (p *Product) Parse(line string) error {
	fields := strings.Split(line, "\t")
	if len(fields) < 2 {
		p.Name = "COULD NOT PARSE"
		return fmt.Errorf("expected at least 2 fields, got %d", len(fields))
	}

	p.Name = fields[flieldName]

	if _, err := fmt.Sscanf(fields[fieldCount], "%f", &p.Count); err != nil {
		return fmt.Errorf("could not parse count for %s: %w", p.Name, err)
	}

	p.Provider = fields[fieldProvider]
	p.providerArgs = fields[fieldArgv:]

	// Remove trailing empty arguments
	for p.providerArgs[len(p.providerArgs)-1] == "" {
		p.providerArgs = p.providerArgs[:len(p.providerArgs)-1]
	}

	return nil
}

func (p *Product) Get(ctx context.Context, r *Registry) error {
	getter, ok := r.data[p.Provider]
	if !ok {
		return fmt.Errorf("%s: provider %s not registered", p.Name, p.Provider)
	}

	log.Debugf("Provider %s: fetching price for %s", p.Provider, p.Name)
	price, err := getter(ctx, p.Name, p.providerArgs...)
	if err != nil {
		return fmt.Errorf("could not get price for %s: %w", p.Name, err)
	}
	log.Debugf("Provider %s: got price for %s", p.Provider, p.Name)

	p.Price = price / p.Count
	return nil
}

type Getter = func(ctx context.Context, name string, args ...string) (float32, error)

type Registry struct {
	data map[string]Getter
}

func (r *Registry) Register(name string, g Getter) {
	if r.data == nil {
		r.data = make(map[string]Getter)
	}

	r.data[name] = g
}
