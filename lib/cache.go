package cache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Item struct {
	payload  string
	created  time.Time
	duration int
}

type Cache struct {
	shards          map[string]*Shard
	cleanupInterval int
}

type Shard struct {
	items       map[string]Item
	deleteCount int
	lock        *sync.RWMutex
}

func (c *Cache) CreateItem(key string, value string, duration int) {
	shardKey := string(key[0])
	c.shards[shardKey].lock.Lock()
	c.shards[shardKey].items[key] = Item{payload: value, created: time.Now(), duration: duration}
	c.shards[shardKey].lock.Unlock()
}

func (c Cache) ReadItem(key string) (string, error) {
	shardKey := string(key[0])
	c.shards[shardKey].lock.RLock()

	item, found := c.shards[shardKey].items[key]
	if !found {
		c.shards[shardKey].lock.RUnlock()
		return "Not Found", errors.New("Not Found")
	}
	c.shards[shardKey].lock.RUnlock()
	return item.payload, nil

}

func (c *Cache) DestroyItem(key string) error {
	shardKey := string(key[1])
	c.shards[shardKey].lock.RLock()

	_, ok := c.shards[shardKey].items[key]

	if !ok {
		c.shards[shardKey].lock.RUnlock()
		return errors.New("Not Found")
	}
	c.shards[shardKey].lock.RUnlock()

	c.shards[shardKey].lock.Lock()
	delete(c.shards[shardKey].items, key)
	c.shards[shardKey].deleteCount++
	c.shards[shardKey].lock.Unlock()
	return nil
}

func (c *Cache) Cleanup() {
	for range time.Tick(time.Second * time.Duration(c.cleanupInterval)) {
		for k := range c.shards {
			c.shards[k].lock.Lock()
			for k2, v := range c.shards[k].items {
				if v.created.Add(time.Second * time.Duration(v.duration)).Before(time.Now()) {
					c.shards[k].deleteCount++
					delete(c.shards[k].items, k2)
				}
			}
			// Fixes go map memory leak
			if c.shards[k].deleteCount > 10 {
				tempMap := c.shards[k].items

				c.shards[k].items = make(map[string]Item)

				for k2, v := range tempMap {
					c.shards[k].items[k2] = v

				}

				c.shards[k].deleteCount = 0
			}

			c.shards[k].lock.Unlock()
		}

	}
}

func (c *Cache) DumpToTerminal() {
	for range time.Tick(time.Second * time.Duration(1)) {
		for k, v := range c.shards {
			c.shards[k].lock.Lock()
			if len(v.items) != 0 {
				fmt.Println(k)
				fmt.Println(v.items)
			}
			c.shards[k].lock.Unlock()
		}
	}
}

func (c *Cache) Purge() {
	for k := range c.shards {
		c.shards[k].lock.Lock()
	}

	for k := range c.shards {
		c.shards[k].items = make(map[string]Item)
	}

	for k := range c.shards {
		c.shards[k].lock.Unlock()
	}
}

func New(cleanupInterval int) Cache {

	shards := createShards()
	return Cache{shards: shards, cleanupInterval: cleanupInterval}
}

func createShards() map[string]*Shard {
	return map[string]*Shard{
		"a": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"b": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"c": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"d": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"e": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"f": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"g": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"h": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"i": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"j": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"k": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"l": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"m": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"n": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"o": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"p": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"q": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"r": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"s": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"t": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"u": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"v": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"w": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"x": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"y": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"z": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"1": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"2": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"3": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"4": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"5": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"6": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"7": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"8": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"9": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
		"0": &Shard{items: make(map[string]Item), lock: new(sync.RWMutex)},
	}
}
