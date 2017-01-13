// Copyright © 2016 Abcum Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lifo

import (
	"container/list"
	"errors"
	"sync"
)

type elem struct {
	sze int
	key string
	val []byte
}

// Cache represents an in-memory Last In First Out cache that is
// safe to use for concurrent writes from multiple goroutines.
type Cache struct {
	size  int
	lock  sync.Mutex
	bytes int
	queue *list.List
	items map[string]*list.Element
}

// New creates and returns a LIFO (Last In First Out) cache with a
// maximum size specified in bytes. The cache size must be a number
// greater than 0, otherwise en error will be returned.
func New(size int) (*Cache, error) {

	if size <= 0 {
		return nil, errors.New("Size must be a positive number")
	}

	c := &Cache{
		size:  size,
		queue: list.New(),
		items: make(map[string]*list.Element),
	}

	return c, nil

}

// Clr removes and clears every item from the cache, and resets its size.
func (c *Cache) Clr() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.clr()
}

// Has checks to see if the key exists in the cache, returning a true
// if it exists and false if not.
func (c *Cache) Has(key string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.has(key)
}

// Get looks up a key's value in the cache. If the value exists in the
// cache then the value is returned, otherwise a nil byte slice is
// returned.
func (c *Cache) Get(key string) []byte {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.get(key)
}

// Del deletes a key from the cache. If the value existed in the cache
// then the value is returned, otherwise a nil byte slice is returned.
func (c *Cache) Del(key string) []byte {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.del(key)
}

// Put puts a new item into the cache using the specified key. If the
// size of the cache will rise above the allowed cache size, the oldest
// items will be removed.
func (c *Cache) Put(key string, val []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.put(key, val)
}

// ---------------------------------------------------------------------------

func (c *Cache) clr() {

	for k := range c.items {
		delete(c.items, k)
	}

	c.queue.Init()

	c.bytes = 0

}

func (c *Cache) has(key string) bool {

	_, ok := c.items[key]

	return ok

}

func (c *Cache) get(key string) []byte {

	if item, ok := c.items[key]; ok {
		return item.Value.(*elem).val
	}

	return nil

}

func (c *Cache) del(key string) []byte {

	if item, ok := c.items[key]; ok {
		c.bytes -= item.Value.(*elem).sze
		delete(c.items, key)
		c.queue.Remove(item)
		return item.Value.(*elem).val
	}

	return nil

}

func (c *Cache) put(key string, val []byte) {

	// Delete the item

	c.del(key)

	// Get the length of the data

	sze := len(val)

	// The item is too big for caching

	if sze > c.size {
		return
	}

	// Free up some data from the cache

	for c.queue.Len() > 0 && sze+c.bytes > c.size {
		if item := c.queue.Front(); item != nil {
			c.del(item.Value.(*elem).key)
		}
	}

	// Insert the element into the cache

	elem := &elem{sze: sze, key: key, val: val}

	item := c.queue.PushFront(elem)

	c.items[key] = item

	c.bytes += sze

	return

}
