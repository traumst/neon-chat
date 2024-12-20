package store

import (
	"fmt"
	"log"
	"neon-chat/src/consts"
	i "neon-chat/src/interfaces"
	"sync"
)

type DoublyLinkedList struct {
	mu    sync.Mutex
	head  *node
	tail  *node
	size  int
	count int
}

type node struct {
	id    uint
	value interface{}
	prev  *node
	next  *node
}

func (n node) GetId() uint {
	return n.id
}

func (n node) Next() (bool, i.Identifiable) {
	if n.next == nil {
		return false, nil
	}
	return true, n.next
}

func NewLinkedList(size int) *DoublyLinkedList {
	if size < 1 || consts.MaxCacheSize < size {
		log.Printf("WARN DoublyLinkedList.NewLinkedList overriding stupid size[%d] to default[%d]",
			size, consts.MaxCacheSize)
		size = consts.MaxCacheSize
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

func (ll *DoublyLinkedList) Get(nodeId uint) (*node, error) {
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
func (ll *DoublyLinkedList) AddHead(new *node) error {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	if new == nil {
		return fmt.Errorf("attempt to add NIL node")
	}
	if ll.count >= ll.size {
		return fmt.Errorf("list is full: count[%d] size[%d]", ll.count, ll.size)
	}
	if ll.head == nil && ll.tail == nil {
		new.prev = nil
		new.next = nil
		ll.head = new
		ll.tail = new
	} else if ll.head != nil && ll.tail != nil {
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

// moves node to head
func (ll *DoublyLinkedList) Bump(node *node) ([]uint, error) {
	if node == nil {
		panic("Attempt to bump NIL node")
	}
	if node.id == ll.head.id {
		return nil, nil
	}
	removed := ll.Remove(node.id)
	if removed == nil {
		return nil, fmt.Errorf("node[%d] not found", node.id)
	}
	err := ll.AddHead(removed)
	if err == nil {
		return nil, nil
	}
	log.Printf("WARN DoublyLinkedList.Bump failed to re-add node[%d], %s", node.id, err.Error())
	pruneCount, pruneIds, pruneErr := ll.Prune(0)
	if pruneErr != nil {
		log.Printf("ERROR DoublyLinkedList.Bump failed prunned[%d], %s", pruneCount, pruneErr.Error())
		return pruneIds, fmt.Errorf("%s, %s", err.Error(), pruneErr.Error())
	}
	retryErr := ll.AddHead(removed)
	if retryErr != nil {
		log.Printf("ERROR DoublyLinkedList.Bump failed adding after crop node[%d], %s", node.id, retryErr.Error())
		return pruneIds, fmt.Errorf("%s, %s", err.Error(), retryErr.Error())
	} else {
		log.Printf("TRACE DoublyLinkedList.Bump added after crop node[%d]", node.id)
		return pruneIds, nil
	}
}

// removes n nodes from tail, defaults to 1 + (ll.size / 8)
func (ll *DoublyLinkedList) Prune(drop int) (int, []uint, error) {
	if drop < 1 || drop > consts.MaxCacheSize {
		def := 1 + (ll.Size() / 8)
		log.Printf("WARN DoublyLinkedList.Prune stupid prune size, should fulfill [1 <= %d <= %d], using default[%d]",
			drop, consts.MaxCacheSize, def)
		drop = def
	}
	ll.mu.Lock()
	defer ll.mu.Unlock()
	count := 0
	ids := []uint{}
	for count < drop {
		removed := ll.removeTail()
		if removed == nil {
			break
		}
		ids = append(ids, removed.id)
		count++
		ll.count--
	}
	return count, ids, nil
}

// removes node from middle
func (ll *DoublyLinkedList) Remove(id uint) *node {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	if ll.count <= 0 {
		return nil
	}
	var removed *node
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
func (ll *DoublyLinkedList) removeHead() *node {
	removed := ll.head
	if ll.head == ll.tail {
		log.Printf("TRACE DoublyLinkedList.removeHead removing last[%d]", ll.head.id)
		ll.head = nil
		ll.tail = nil
	} else {
		log.Printf("TRACE DoublyLinkedList.removeHead removing head[%d]", ll.head.id)
		ll.head = removed.next
	}
	return removed
}
func (ll *DoublyLinkedList) removeTail() *node {
	if ll.tail == nil {
		return nil
	}
	removed := ll.tail
	if ll.tail == ll.head {
		log.Printf("TRACE DoublyLinkedList.removeTail removing last[%d]", ll.head.id)
		ll.head = nil
		ll.tail = nil
	} else {
		log.Printf("TRACE DoublyLinkedList.removeTail removing tail[%d]", ll.head.id)
		ll.tail = removed.prev
		ll.tail.next = nil
	}
	return removed
}
func (ll *DoublyLinkedList) removeNode(id uint) *node {
	log.Printf("TRACE DoublyLinkedList.removeNode any node[%d]", id)
	current := ll.head
	for current != nil && current.id != id {
		if current.id == id {
			log.Printf("TRACE DoublyLinkedList.removeNode removing node[%d]", id)
			if current.prev != nil {
				log.Printf("TRACE DoublyLinkedList.removeNode node[%d] link prev[%d] to next[%d]",
					id, current.prev.id, current.next.id)
				current.prev.next = current.next
			}
			if current.next != nil {
				log.Printf("TRACE DoublyLinkedList.removeNode node[%d] link next[%d] to prev[%d]",
					id, current.next.id, current.prev.id)
				current.next.prev = current.prev
			}
			break
		}
		current = current.next
	}
	log.Printf("INFO DoublyLinkedList.removeNode not found node[%d]", id)
	return current
}
