package store

import (
	"fmt"
	"sync"
)

type LRUCache struct {
	mu   sync.Mutex
	dict map[uint]*Node
	list *DoublyLinkedList
}

func NewLRUCache(size int) *LRUCache {
	return &LRUCache{
		dict: make(map[uint]*Node),
		list: NewLinkedList(size),
	}
}

func (cache *LRUCache) Size() int {
	return cache.list.Size()
}

func (cache *LRUCache) Count() int {
	return cache.list.Count()
}

func (cache *LRUCache) Set(key uint, value interface{}) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	newNode := Node{id: key, value: value}
	err := cache.list.AddHead(&newNode)
	if err != nil {
		return err
	}

	cache.dict[key] = &newNode
	return err
}

// gets the value stored in cache and bumps it to the head
func (cache *LRUCache) Get(key uint) (interface{}, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	node, ok := cache.dict[key]
	if !ok {
		return "", fmt.Errorf("key not found")
	}
	err := cache.list.Bump(node)
	return node.value, err
}

// removes the value stored in cache and returns it
func (cache *LRUCache) Take(key uint) (interface{}, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	node, ok := cache.dict[key]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}
	delete(cache.dict, key)
	removed := cache.list.Remove(node.id)
	if removed == nil {
		panic(fmt.Sprintf("Node[%d] not found in list", key))
	}
	return node.value, nil
}

// removes up-to n nodes from the tail
func (cache *LRUCache) Drop(n int) (int, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	dropCount := 0
	for i := 0; i < n; i++ {
		tail := cache.list.tail
		if tail == nil {
			return dropCount, nil
		}
		delete(cache.dict, tail.id)
		removed := cache.list.Remove(tail.id)
		if removed == nil {
			return dropCount, fmt.Errorf("node[%d] removed as NIL", tail.id)
		}
		dropCount++
	}
	return dropCount, nil
}
