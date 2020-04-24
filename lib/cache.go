package cache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Item struct {
	payload    string
	created    time.Time
	expireTime int
}

type Cache struct {
	items           map[string]Item
	cleanupInterval int
	lock            sync.RWMutex
}

func (c *Cache) CreateItem(key string, value string, expiry int) {

	c.lock.Lock()
	c.items[key] = Item{payload: value, created: time.Now()}

	c.lock.Unlock()

}

func (c Cache) ReadItem(key string) (string, error) {

	c.lock.RLock()
	item, found := c.items[key]
	if !found {
		c.lock.RUnlock()
		return "Not Found", errors.New("Not Found")
	} else {
		c.lock.RUnlock()
		return item.payload, nil
	}

}

func (c *Cache) DestroyItem(key string) error {
	c.lock.RLock()
	_, ok := c.items[key]

	if !ok {
		c.lock.RUnlock()
		return errors.New("Not Found")
	}
	c.lock.RUnlock()

	c.lock.Lock()
	delete(c.items, key)
	c.lock.Unlock()
	return nil
}

func (c *Cache) Cleanup() {
	for range time.Tick(time.Second * time.Duration(c.cleanupInterval)) {
		c.lock.Lock()
		for k, v := range c.items {

			if v.created.Add(time.Second * time.Duration(v.expireTime)).Before(time.Now()) {
				delete(c.items, k)
			}
		}
		c.lock.Unlock()
	}
}

func (c *Cache) DumpToTerminal() {
	for range time.Tick(time.Second * time.Duration(1)) {
		c.lock.Lock()
		fmt.Println(c.items)
		c.lock.Unlock()
	}
}

func (c *Cache) Purge() {
	c.lock.Lock()
	c.items = map[string]Item{}
	c.lock.Unlock()
}

func New(cleanupInterval int) Cache {
	initItems := make(map[string]Item)
	return Cache{items: initItems, cleanupInterval: cleanupInterval}
}
