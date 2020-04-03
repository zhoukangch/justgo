package gee

import (
	"github.com/zhoukangch/justgo/gee-cache/lru"
	"sync"
)

type Cache struct {
	lru *lru.Cache
	*sync.RWMutex
}

func (c *Cache) Add(key string, value ByteView) {
	c.Lock()
	c.lru.Add(key, value)
	c.Unlock()
}

func (c *Cache) Get(key string) (value ByteView, ok bool) {
	c.Lock()
	defer c.Unlock()
	if v, exist := c.lru.Get(key); exist {
		return v.(ByteView), true
	} else {
		return
	}
}
