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

package cachr

import (
	"container/list"
	"errors"
)

// LRU represents an in-memory Least Recently Used cache.
type LRU struct {
	size  int
	bytes int
	queue *list.List
	items map[string]*list.Element
}

// NewLRU creates and returns a LRU (Least Recently Used) cache with a
// maximum size specified in bytes. The cache size must be a number
// greater than 0, otherwise en error will be returned.
func NewLRU(size int) (*LRU, error) {

	if size <= 0 {
		return nil, errors.New("Size must be a positive number")
	}

	var c *LRU

	c.size = size
	c.queue = list.New()
	c.items = make(map[string]*list.Element)

	return c, nil

}

// Clr removes and clears every item from the cache, and resets its size.
func (c *LRU) Clr() {

	for k := range c.items {
		delete(c.items, k)
	}

	c.queue.Init()

	c.bytes = 0

}

// Has checks to see if the key exists in the cache, returning a true
// if it exists and false if not. It does not affect the cache eviction
// marker on the cache item.
func (c *LRU) Has(key string) bool {

	_, ok := c.items[key]

	return ok

}

// Get looks up a key's value in the cache. If the value exists in the
// cache then the value is returned, otherwise a nil byte slice is
// returned. The item is subsequently marked as most recently used.
func (c *LRU) Get(key string) []byte {

	if item, ok := c.items[key]; ok {
		data := item.Value.([]byte)
		c.queue.MoveToFront(item)
		return data
	}

	return nil

}

// Del deletes a key from the cache. If the value existed in the cache
// then the value is returned, otherwise a nil byte slice is returned.
func (c *LRU) Del(key string) []byte {

	if item, ok := c.items[key]; ok {
		data := item.Value.([]byte)
		c.queue.Remove(item)
		c.bytes -= len(data)
		return data
	}

	return nil

}

// Put puts a new item into the cache using the specified key, and
// marks it as most recently used. If the size of the cache will rise
// above the allowed cache size, the oldest items will be removed.
func (c *LRU) Put(key string, val []byte) {

	c.Del(key)

	for c.queue.Len() > 0 && len(val)+c.bytes > c.size {
		if item := c.queue.Back(); item != nil {
			c.size -= len(item.Value.([]byte))
			c.queue.Remove(item)
		}
	}

	item := c.queue.PushFront(val)

	c.items[key] = item

	return

}
