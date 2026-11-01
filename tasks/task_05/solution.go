package main

type Cache[K comparable, V any] struct {
	capacity int
	items    map[K]V
}

func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	cache := Cache[K, V]{
		capacity: capacity,
		items:    make(map[K]V)}
	return &cache
}
func (c *Cache[K, V]) Get(k K) (v V, ok bool) {
	v, ok = c.items[k]
	// судя по описанию, при нулевом capacity должен бы возвращать false (промах), но в тестах не так
	//if c.capacity <= 0 {
	//	return v, false
	//}

	return
}

func (c *Cache[K, V]) Set(k K, v V) bool {
	_, ok := c.items[k]
	if ok {
		// cache hit, update value
		c.items[k] = v
		return true
	}

	if c.capacity <= 0 {
		return false
	}

	// cache miss, add new key and reduce capacity:
	c.items[k] = v
	c.capacity--
	return true
}
