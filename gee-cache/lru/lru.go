package lru

import "container/list"

type Cache struct {
	maxBytes  int64
	nBytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int64
}

func NewCache(maxBytes int64, onEvited func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		nBytes:    0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvited,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, exist := c.cache[key]; exist {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	} else {
		return nil, false
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += value.Len() - kv.value.Len()
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.nBytes += value.Len()
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		if remove := c.RemoveOldest(); !remove {
			break
		}
	}
}

func (c *Cache) RemoveOldest() bool {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= kv.value.Len()
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
		return true
	}
	return false
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
