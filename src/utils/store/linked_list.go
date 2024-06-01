package store

import (
	"fmt"
	"sync"
)

type DoublyLinkedList struct {
	mu    sync.Mutex
	head  *Node
	tail  *Node
	size  int
	count int
}

type Node struct {
	id    uint
	value string
	prev  *Node
	next  *Node
}

func NewLinkedList(size int) *DoublyLinkedList {
	return &DoublyLinkedList{head: nil, tail: nil, size: size, count: 0}
}

func (ll *DoublyLinkedList) Size() int {
	return ll.size
}

func (ll *DoublyLinkedList) Count() int {
	return ll.count
}

func (ll *DoublyLinkedList) Add(id uint, value string) (*Node, error) {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	if ll.count >= ll.size {
		return nil, fmt.Errorf("LinkedList is full")
	}
	new := Node{id: id, value: value, prev: nil, next: nil}
	if ll.head == nil && ll.tail == nil {
		ll.head = &new
		ll.tail = &new
	} else if ll.head != nil && ll.tail != nil {
		new.prev = ll.tail
		ll.tail.next = &new
	} else {
		panic(fmt.Sprintf("LinkedList is in an invalid state, head[%+v] tail[%+v]", ll.head, ll.tail))
	}
	ll.count++
	return &new, nil
}

func (ll *DoublyLinkedList) Drop(id uint) *Node {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	if ll.size <= 0 {
		return nil
	}
	current := ll.head
	for current != nil && current.id != id {
		if current.id == id {
			current.prev = current.next
			break
		}
		current = current.next
	}
	return current
}
