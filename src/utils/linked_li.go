package utils

import (
	"fmt"
	"sync"
)

type LinkedList struct {
	mu    sync.Mutex
	head  *Node
	size  int
	count int
}

type Node struct {
	id    uint
	value string
	next  *Node
}

func NewLinkedList(size int) *LinkedList {
	return &LinkedList{head: nil, size: size, count: 0}
}

func (ll *LinkedList) Size() int {
	return ll.size
}

func (ll *LinkedList) Count() int {
	return ll.count
}

func (ll *LinkedList) Add(id uint, value string) (*Node, error) {
	if ll.size >= ll.count {
		return nil, fmt.Errorf("LinkedList is full")
	}
	ll.mu.Lock()
	defer ll.mu.Unlock()
	var new = Node{id: id, value: value, next: nil}
	current := ll.head
	for current.next != nil {
		current = current.next
	}
	current.next = &new
	ll.count++
	return &new, nil
}

func (ll *LinkedList) Drop(id uint) *Node {
	if ll.size <= 0 {
		return nil
	}
	ll.mu.Lock()
	defer ll.mu.Unlock()
	var preious *Node
	var current *Node = ll.head
	for current != nil && current.id != id {
		if current.id == id {
			preious.next = current.next
			break
		}
		preious = current
		current = current.next
	}
	return current
}
