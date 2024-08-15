package store

import (
	"fmt"
	"math"
	"sync"
)

type LRUCache struct {
	mu   sync.Mutex
	dict map[uint]*node
	list *DoublyLinkedList
	// 0 to 1, defaults to 1/4 of keys will to be dropped when full
	cleanupRatio float64
}

func NewLRUCache(size int) *LRUCache {
	return &LRUCache{
		dict:         make(map[uint]*node, size),
		list:         NewLinkedList(size),
		cleanupRatio: 0.25,
	}
}

func (cache *LRUCache) Size() int {
	if cache == nil {
		panic("cache is nil")
	}
	if cache.list == nil {
		panic("cache.list is nil")
	}
	return cache.list.Size()
}

func (cache *LRUCache) Count() int {
	if cache == nil {
		panic("cache is nil")
	}
	if cache.list == nil {
		panic("cache.list is nil")
	}
	return cache.list.Count()
}

func (cache *LRUCache) CleanupRatio() int {
	if cache == nil {
		panic("cache is nil")
	}
	return int(math.Ceil(cache.cleanupRatio * float64(cache.Size())))
}

func (cache *LRUCache) Set(key uint, value interface{}) error {
	if cache == nil {
		panic("cache is nil")
	}
	if cache.Count()+1 > cache.Size() {
		n := cache.CleanupRatio()
		_, err := cache.Drop(n)
		if err != nil {
			return fmt.Errorf("cache is full, failed to drop [%d] keys: %s", n, err.Error())
		}
	}
	if _, ok := cache.dict[key]; ok {
		_, err := cache.Take(key)
		if err != nil {
			return fmt.Errorf("failed to pop taken key[%d]: %s", key, err)
		}
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	newNode := node{id: key, value: value}
	err := cache.list.AddHead(&newNode)
	if err != nil {
		return err
	}

	cache.dict[key] = &newNode
	return nil
}

// gets the value stored in cache and bumps it to the head
func (cache *LRUCache) Get(key uint) (interface{}, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	node, ok := cache.dict[key]
	if !ok {
		return "", fmt.Errorf("key not found")
	}
	dropIds, err := cache.list.Bump(node)
	if len(dropIds) > 0 {
		for _, dropId := range dropIds {
			delete(cache.dict, dropId)
		}
	}
	return node.value, err
}

// scans cache and returns all existing keys in no particular order
func (cache *LRUCache) Keys() []uint {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	keys := make([]uint, 0)
	for _, node := range cache.dict {
		keys = append(keys, node.id)
	}
	return keys
}

// removes the value stored in cache and returns it
func (cache *LRUCache) Take(key uint) (interface{}, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	node, ok := cache.dict[key]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}
	removed := cache.list.Remove(node.id)
	if removed == nil {
		panic(fmt.Sprintf("Node[%d] not found in list", key))
	}
	delete(cache.dict, key)
	return node.value, nil
}

// removes at least 1 and up-to n oldest nodes
func (cache *LRUCache) Drop(n int) (int, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	dropCount, dropIds, err := cache.list.Prune(n)
	if err != nil || dropCount == 0 {
		return dropCount, err
	}
	for _, dropId := range dropIds {
		delete(cache.dict, dropId)
	}
	return dropCount, nil
}
