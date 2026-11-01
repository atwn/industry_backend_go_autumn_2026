package main

import (
	"container/list"
	"sync"
)

type LRU[K comparable, V any] interface {
	Get(key K) (value V, ok bool)
	Set(key K, value V)
}
type entry[K comparable, V any] struct {
	key   K
	value V
}
type LRUCache[K comparable, V any] struct {
	capacity int
	mu       sync.Mutex
	ll       list.List
	items    map[K]*list.Element
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	cache := LRUCache[K, V]{
		capacity: capacity,
		mu:       sync.Mutex{},
		ll:       list.List{},
		items:    make(map[K]*list.Element)}
	return &cache
}

func (c *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ref, ok := c.items[key]
	if ok {
		// cache hit
		c.ll.MoveToFront(ref)
		value = ref.Value.(*entry[K, V]).value
	}
	return
}

func (c *LRUCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ref, ok := c.items[key]
	if ok {
		// cache hit: change value and return
		ref.Value.(*entry[K, V]).value = value
		return
	}
	// free space
	if c.capacity <= 0 {
		if c.ll.Len() > 0 {
			stale := c.ll.Remove(c.ll.Back()).(*entry[K, V]) // back = the oldest
			delete(c.items, stale.key)
			c.capacity++
		}
	}
	// cache miss: add a new item
	if c.capacity > 0 {
		c.items[key] = c.ll.PushFront(&entry[K, V]{key: key, value: value}) // front = the newest
		c.capacity--
	}
}
