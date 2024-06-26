package providers

import (
	"context"
)

// Provider is an interface that represents a grocery store provider for a single product.
type Provider interface {
	Name() string
	FetchPrice(ctx context.Context, pid ProductCode) (float32, error)
	ValidateCode(code ProductCode) error
}
type ProductCode [3]string

var Default = &Registry{}

type Registry struct {
	data map[string]Provider
}

func (r *Registry) Register(prov Provider) {
	if r.data == nil {
		r.data = make(map[string]Provider)
	}

	r.data[prov.Name()] = prov
}

func (r *Registry) Lookup(name string) (Provider, bool) {
	if f, ok := r.data[name]; ok {
		return f, true
	}

	return nil, false
}

func Register(p Provider) {
	Default.Register(p)
}

func Lookup(name string) (Provider, bool) {
	return Default.Lookup(name)
}
