// Package cache provides a simple LRU cache used to store Discord objects
// wip
package cache

import (
	"sync"
	"time"
)

type Cache struct {
	sync.Mutex
	holds    map[string]*Item
	capacity int
}

type Item struct {
	lastUsed int64

	item interface{}
}

// NewCache constructs a new cache with the given capacity
func NewCache(capacity int) *Cache {
	return &Cache{
		Mutex:    sync.Mutex{},
		capacity: capacity,
		holds:    make(map[string]*Item),
	}
}

func newItem(item interface{}) *Item {
	return &Item{
		lastUsed: time.Now().UnixNano(),
		item:     item,
	}
}

func (c *Cache) Add(id string, item interface{}) {
	insertable := newItem(item)
	c.holds[id] = insertable
}

func (c *Cache) Remove(id string) {
	c.Lock()
	delete(c.holds, id)
	c.Unlock()
}

func (c *Cache) Update(id string, ele interface{}) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.holds[id]; !ok {
		c.Add(id, ele)
	} else {
		c.holds[id] = &Item{
			lastUsed: time.Now().UnixNano(),
			item:     ele,
		}
	}
}

func (c *Cache) Size() int {
	var i int
	for range c.holds {
		i++
	}

	return i
}

func (c *Cache) Has(key string) bool {
	_, ok := c.holds[key]
	return ok
}

func (c *Cache) clearLRU(exception string) {
	// don't clear the cache if we haven't reached full capacity
	if c.Size() < c.capacity {
		return
	}
	// also don't clear if it's an infinite cache, i.e capacity is 0
	if c.capacity == 0 {
		return
	}
	c.Lock()
	defer c.Unlock()
	// primitive sorting, there's gotta be a better way to do this
	var lowest string
	for key := range c.holds {
		if lowest == "" {
			lowest = key
			continue
		}
		if c.holds[key].lastUsed < c.holds[lowest].lastUsed {
			lowest = key
		}
	}

	delete(c.holds, lowest)
}
