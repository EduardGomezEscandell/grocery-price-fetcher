package database

type CachedLookup[T any] struct {
	cache  map[string]T
	lookup func(string) (T, bool)
}

func NewCachedLookup[T any](lookup func(string) (T, bool)) CachedLookup[T] {
	return CachedLookup[T]{
		cache:  make(map[string]T),
		lookup: lookup,
	}
}

func (c *CachedLookup[T]) Lookup(name string) (T, bool) {
	if v, ok := c.cache[name]; ok {
		return v, true
	}

	v, ok := c.lookup(name)
	if !ok {
		var t T
		return t, false
	}

	c.cache[name] = v
	return v, true
}
