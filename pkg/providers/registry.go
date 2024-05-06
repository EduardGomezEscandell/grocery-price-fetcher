package providers

import (
	"context"
	"errors"
	"strings"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/logger"
)

// Provider is an interface that represents a grocery store provider for a single product.
type Provider interface {
	UnmarshalTSV(cols ...string) error
	UnmarshalMap(argv map[string]string) error
	FetchPrice(ctx context.Context, log logger.Logger) (float32, error)
}

// Factory is a function that creates a new Provider.
type Factory func() Provider

var Default = &Registry{}

type Registry struct {
	data map[string]Factory
}

func (r *Registry) Register(name string, prov Factory) {
	if r.data == nil {
		r.data = make(map[string]Factory)
	}

	r.data[name] = prov
}

func (r *Registry) ParseTSV(providerName string, cols []string) (Provider, error) {
	factory, ok := r.data[providerName]
	if !ok {
		return nil, errors.New("provider not found")
	}

	g := factory()

	if err := g.UnmarshalTSV(trim(cols)...); err != nil {
		return nil, err
	}

	return g, nil
}

func (r *Registry) ParseMap(providerName string, argv map[string]string) (Provider, error) {
	factory, ok := r.data[providerName]
	if !ok {
		return nil, errors.New("provider not found")
	}

	g := factory()

	if err := g.UnmarshalMap(argv); err != nil {
		return nil, err
	}

	return g, nil
}

func Register(name string, g Factory) {
	Default.Register(name, g)
}

func ParseTSV(providerName string, cols []string) (Provider, error) {
	return Default.ParseTSV(providerName, cols)
}

func ParseMap(providerName string, argv map[string]string) (Provider, error) {
	return Default.ParseMap(providerName, argv)
}

func trim(arg []string) []string {
	for i, a := range arg {
		arg[i] = strings.TrimSpace(a)
	}

	// Remove trailing empty arguments
	for i := len(arg); i > 0 && arg[i-1] == ""; i = len(arg) {
		arg = arg[:i-1]
	}

	return arg
}
