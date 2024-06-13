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
	if cache.Count()+1 > cache.Size() {
		_, err := cache.Drop(1 + cache.Size()/8)
		if err != nil {
			return err
		}
	}
	cache.mu.Lock()
	defer cache.mu.Unlock()
	newNode := Node{id: key, value: value}
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
