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
	items map[string]Item
	lock  sync.RWMutex
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
	c.shards[shardKey].lock.Unlock()
	return nil
}

func (c *Cache) Cleanup() {
	for range time.Tick(time.Second * time.Duration(c.cleanupInterval)) {

		for k := range c.shards {
			c.shards[k].lock.Lock()
			for k2, v := range c.shards[k].items {
				if v.created.Add(time.Second * time.Duration(v.duration)).Before(time.Now()) {
					delete(c.shards[k].items, k2)
				}
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
		"a": &Shard{items: make(map[string]Item)},
		"b": &Shard{items: make(map[string]Item)},
		"c": &Shard{items: make(map[string]Item)},
		"d": &Shard{items: make(map[string]Item)},
		"e": &Shard{items: make(map[string]Item)},
		"f": &Shard{items: make(map[string]Item)},
		"g": &Shard{items: make(map[string]Item)},
		"h": &Shard{items: make(map[string]Item)},
		"i": &Shard{items: make(map[string]Item)},
		"j": &Shard{items: make(map[string]Item)},
		"k": &Shard{items: make(map[string]Item)},
		"l": &Shard{items: make(map[string]Item)},
		"m": &Shard{items: make(map[string]Item)},
		"n": &Shard{items: make(map[string]Item)},
		"o": &Shard{items: make(map[string]Item)},
		"p": &Shard{items: make(map[string]Item)},
		"q": &Shard{items: make(map[string]Item)},
		"r": &Shard{items: make(map[string]Item)},
		"s": &Shard{items: make(map[string]Item)},
		"t": &Shard{items: make(map[string]Item)},
		"u": &Shard{items: make(map[string]Item)},
		"v": &Shard{items: make(map[string]Item)},
		"w": &Shard{items: make(map[string]Item)},
		"x": &Shard{items: make(map[string]Item)},
		"y": &Shard{items: make(map[string]Item)},
		"z": &Shard{items: make(map[string]Item)},
		"1": &Shard{items: make(map[string]Item)},
		"2": &Shard{items: make(map[string]Item)},
		"3": &Shard{items: make(map[string]Item)},
		"4": &Shard{items: make(map[string]Item)},
		"5": &Shard{items: make(map[string]Item)},
		"6": &Shard{items: make(map[string]Item)},
		"7": &Shard{items: make(map[string]Item)},
		"8": &Shard{items: make(map[string]Item)},
		"9": &Shard{items: make(map[string]Item)},
		"0": &Shard{items: make(map[string]Item)},
	}
}
