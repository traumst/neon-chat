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
	value interface{}
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
	ll.mu.Lock()
	defer ll.mu.Unlock()
	return ll.size
}

func (ll *DoublyLinkedList) Count() int {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	return ll.count
}

func (ll *DoublyLinkedList) Get(nodeId uint) (*Node, error) {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	current := ll.head
	for current != nil {
		if current.id == nodeId {
			return current, nil
		}
		current = current.next
	}
	return nil, fmt.Errorf("node[%d] not found", nodeId)
}

// adds node to head
func (ll *DoublyLinkedList) AddHead(new *Node) error {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	if ll.count >= ll.size {
		free := 1 + (ll.size / 8)
		cropped, cropErr := ll.Crop(free)
		if cropErr != nil {
			log.Printf("DoublyLinkedList.AddHead ERROR failed cropped[%d], %s", cropped, cropErr.Error())
			return cropErr
		}
	}
	if ll.head == nil && ll.tail == nil {
		log.Printf("DoublyLinkedList.AddHead TRACE adding first[%d]", new.id)
		new.prev = nil
		new.next = nil
		ll.head = new
		ll.tail = new
	} else if ll.head != nil && ll.tail != nil {
		log.Printf("DoublyLinkedList.AddHead TRACE adding head[%d]", new.id)
		new.next = ll.head
		new.prev = nil
		ll.head.prev = new
		ll.head = new
	} else {
		panic(fmt.Sprintf("LinkedList is in an invalid state, head[%+v] tail[%+v]", ll.head, ll.tail))
	}
	ll.count++
	return nil
}

// creates node and adds to tail
func (ll *DoublyLinkedList) AddTail(id uint, value interface{}) (*Node, error) {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	if ll.count >= ll.size {
		free := 1 + (ll.size / 8)
		cropped, cropErr := ll.Crop(free)
		if cropErr != nil {
			log.Printf("DoublyLinkedList.AddHead ERROR failed cropped[%d], %s", cropped, cropErr.Error())
			return nil, cropErr
		}
	}
	new := &Node{id: id, value: value, prev: nil, next: nil}
	if ll.head == nil && ll.tail == nil {
		log.Printf("DoublyLinkedList.Add TRACE adding first[%d]", new.id)
		ll.head = new
		ll.tail = new
	} else if ll.head != nil && ll.tail != nil {
		log.Printf("DoublyLinkedList.Add TRACE adding tail[%d]", new.id)
		new.prev = ll.tail
		ll.tail.next = new
		ll.tail = new
	} else {
		panic(fmt.Sprintf("LinkedList is in an invalid state, head[%+v] tail[%+v]", ll.head, ll.tail))
	}
	ll.count++
	return new, nil
}

// moves node to head
func (ll *DoublyLinkedList) Bump(node *Node) error {
	removed := ll.Remove(node.id)
	if removed == nil {
		return fmt.Errorf("node[%d] not found", node.id)
	}
	err := ll.AddHead(removed)
	if err != nil {
		log.Printf("DoublyLinkedList.Bump WARN failed to add node[%d], %s", node.id, err.Error())
	} else {
		log.Printf("DoublyLinkedList.Bump TRACE added node[%d]", node.id)
		return nil
	}
	free := 1 + (ll.size / 8)
	cropped, cropErr := ll.Crop(free)
	if cropErr != nil {
		log.Printf("DoublyLinkedList.Bump ERROR failed cropped[%d], %s", cropped, cropErr.Error())
		return fmt.Errorf("%s, %s", err.Error(), cropErr.Error())
	}
	retryErr := ll.AddHead(removed)
	if retryErr != nil {
		log.Printf("DoublyLinkedList.Bump ERROR failed adding after crop node[%d], %s", node.id, retryErr.Error())
		return fmt.Errorf("%s, %s", err.Error(), retryErr.Error())
	} else {
		log.Printf("DoublyLinkedList.Bump TRACE added after crop node[%d]", node.id)
		return nil
	}
}

// removes up-to n nodes from the tail
func (ll *DoublyLinkedList) Crop(n int) (int, error) {
	log.Printf("DoublyLinkedList.Crop TRACE removing [%d] nodes from the head", n)
	ll.mu.Lock()
	defer ll.mu.Unlock()
	if n < MinSize || n > MaxSize {
		return 0, fmt.Errorf("invalid crop size[%d] should be [%d < n <= %d]", n, MinSize, MaxSize)
	}
	removeCount := 0
	for removeCount < n && ll.head != nil {
		removed := ll.removeTail()
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
		log.Printf("DoublyLinkedList.removeHead TRACE removing last[%d]", ll.head.id)
		ll.head = nil
		ll.tail = nil
	} else {
		log.Printf("DoublyLinkedList.removeHead TRACE removing head[%d]", ll.head.id)
		ll.head = removed.next
	}
	return removed
}
func (ll *DoublyLinkedList) removeTail() *Node {
	removed := ll.tail
	if ll.tail == ll.head {
		log.Printf("DoublyLinkedList.removeTail TRACE removing last[%d]", ll.head.id)
		ll.head = nil
		ll.tail = nil
	} else {
		log.Printf("DoublyLinkedList.removeTail TRACE removing tail[%d]", ll.head.id)
		ll.tail = removed.prev
		ll.tail.next = nil
	}
	return removed
}
func (ll *DoublyLinkedList) removeNode(id uint) *Node {
	log.Printf("DoublyLinkedList.removeNode TRACE any node[%d]", id)
	// TODO check from head and tail
	current := ll.head
	for current != nil && current.id != id {
		if current.id == id {
			log.Printf("DoublyLinkedList.removeNode TRACE removing node[%d]", id)
			if current.prev != nil {
				log.Printf("DoublyLinkedList.removeNode TRACE node[%d] link prev[%d] to next[%d]",
					id, current.prev.id, current.next.id)
				current.prev.next = current.next
			}
			if current.next != nil {
				log.Printf("DoublyLinkedList.removeNode TRACE node[%d] link next[%d] to prev[%d]",
					id, current.next.id, current.prev.id)
				current.next.prev = current.prev
			}
			break
		}
		current = current.next
	}
	log.Printf("DoublyLinkedList.removeNode INFO not found node[%d]", id)
	return current
}
