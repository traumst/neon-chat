package store

import (
	"fmt"
	"log"
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

const (
	MinSize = 1
	MaxSize = 1024
)

func NewLinkedList(size int) *DoublyLinkedList {
	if MinSize > size || size > MaxSize {
		panic(fmt.Sprintf("Invalid size[%d] should be [%d <= size < %d]", size, MinSize, MaxSize))
	}
	return &DoublyLinkedList{head: nil, tail: nil, size: size, count: 0}
}

func (ll *DoublyLinkedList) Size() int {
	return ll.size
}

func (ll *DoublyLinkedList) Count() int {
	return ll.count
}

// creates node and adds to tail
func (ll *DoublyLinkedList) Add(id uint, value string) (*Node, error) {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	if ll.count >= ll.size {
		return nil, fmt.Errorf("LinkedList is full")
	}
	new := &Node{id: id, value: value, prev: nil, next: nil}
	if ll.head == nil && ll.tail == nil {
		ll.head = new
		ll.tail = new
	} else if ll.head != nil && ll.tail != nil {
		new.prev = ll.tail
		ll.tail.next = new
		ll.tail = new
	} else {
		panic(fmt.Sprintf("LinkedList is in an invalid state, head[%+v] tail[%+v]", ll.head, ll.tail))
	}
	ll.count++
	return new, nil
}

// adds node to tail
func (ll *DoublyLinkedList) AddTail(new *Node) error {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	if ll.count >= ll.size {
		return fmt.Errorf("LinkedList is full")
	}
	if ll.head == nil && ll.tail == nil {
		new.prev = nil
		new.next = nil
		ll.head = new
		ll.tail = new
	} else if ll.head != nil && ll.tail != nil {
		new.prev = ll.tail
		new.next = nil
		ll.tail.next = new
		ll.tail = new
	} else {
		panic(fmt.Sprintf("LinkedList is in an invalid state, head[%+v] tail[%+v]", ll.head, ll.tail))
	}
	ll.count++
	return nil
}

// moves node to tail
func (ll *DoublyLinkedList) Bump(node *Node) error {
	removed := ll.Remove(node.id)
	if removed == nil {
		return fmt.Errorf("node[%d] not found", node.id)
	}
	err := ll.AddTail(removed)
	if err != nil {
		log.Printf("DoublyLinkedList.Bump WARN failed to add node[%d], %s", node.id, err.Error())
	} else {
		log.Printf("DoublyLinkedList.Bump TRACE added node[%d]", node.id)
		return nil
	}
	cropped, cropErr := ll.Crop(1 + (ll.Size() / 8))
	if cropErr != nil {
		log.Printf("DoublyLinkedList.Bump ERROR failed cropped[%d], %s",
			cropped, cropErr.Error())
		return fmt.Errorf("%s, %s", err.Error(), cropErr.Error())
	}
	retryErr := ll.AddTail(removed)
	if retryErr != nil {
		log.Printf("DoublyLinkedList.Bump ERROR failed adding after crop node[%d], %s", node.id, retryErr.Error())
		return fmt.Errorf("%s, %s", err.Error(), retryErr.Error())
	} else {
		log.Printf("DoublyLinkedList.Bump TRACE added after crop node[%d]", node.id)
		return nil
	}
}

// removes up-to n nodes from the head
func (ll *DoublyLinkedList) Crop(n int) (int, error) {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	if n < MinSize || n > MaxSize {
		return 0, fmt.Errorf("invalid crop size[%d] should be [%d < n <= %d]", n, MinSize, MaxSize)
	}
	removeCount := 0
	for removeCount < n && ll.head != nil {
		removed := ll.removeHead()
		if removed != nil {
			removeCount++
			ll.count--
		}
	}
	return removeCount, nil
}

// removes node from middle
func (ll *DoublyLinkedList) Remove(id uint) *Node {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	if ll.count <= 0 {
		return nil
	}
	var removed *Node
	if ll.head.id == id {
		removed = ll.removeHead()
	} else if ll.tail.id == id {
		removed = ll.removeTail()
	} else {
		removed = ll.removeNode(id)
	}
	if removed != nil {
		ll.count--
	}
	return removed
}
func (ll *DoublyLinkedList) removeHead() *Node {
	removed := ll.head
	if ll.head == ll.tail {
		ll.head = nil
		ll.tail = nil
	} else {
		ll.head = removed.next
	}
	return removed
}
func (ll *DoublyLinkedList) removeTail() *Node {
	removed := ll.tail
	if ll.tail == ll.head {
		ll.head = nil
		ll.tail = nil
	} else {
		ll.tail = removed.prev
	}
	return removed
}
func (ll *DoublyLinkedList) removeNode(id uint) *Node {
	// TODO check from head and tail
	current := ll.head
	for current != nil && current.id != id {
		if current.id == id {
			if current.prev != nil {
				current.prev.next = current.next
			}
			if current.next != nil {
				current.next.prev = current.prev
			}
			break
		}
		current = current.next
	}
	return current
}