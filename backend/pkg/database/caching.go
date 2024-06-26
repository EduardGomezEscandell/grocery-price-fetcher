package database

type CachedLookup[K comparable, V any] struct {
	cache  map[K]V
	lookup func(K) (V, error)
}

func NewCachedLookup[K comparable, T any](lookup func(K) (T, error)) CachedLookup[K, T] {
	return CachedLookup[K, T]{
		cache:  make(map[K]T),
		lookup: lookup,
	}
}

func (c *CachedLookup[K, V]) Lookup(k K) (V, error) {
	if v, ok := c.cache[k]; ok {
		return v, nil
	}

	v, err := c.lookup(k)
	if err != nil {
		var v V
		return v, err
	}

	c.cache[k] = v
	return v, nil
}
