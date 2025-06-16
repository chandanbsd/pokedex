package pokecache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	CreatedAt time.Time
	val       []byte
}

type Cache struct {
	store map[string]CacheEntry
	mutex sync.Mutex
}

//const cacheDuration = 5 * time.Second
//var timeChan chan time.Time = make(chan time.Time)

func NewCache(interval time.Duration) *Cache {

	c := &Cache{
		store: map[string]CacheEntry{},
		mutex: sync.Mutex{},
	}

	ticker := time.NewTicker(interval)

	go c.reapLoop(interval, ticker.C)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.store[key] = CacheEntry{
		CreatedAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	entry, ok := c.store[key]

	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration, timeChan <-chan time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	timeVal := <-timeChan

	for key, val := range c.store {
		if timeVal.Sub(val.CreatedAt) > interval {
			delete(c.store, key)
		}
	}
}
