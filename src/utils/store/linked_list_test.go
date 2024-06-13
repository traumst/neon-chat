package store

import (
	"testing"
)

func TestNewLinkedList(t *testing.T) {
	ll := NewLinkedList(3)
	if ll.Size() != 3 {
		t.Fatalf("Expected size to be 3")
	}
	if ll.Count() != 0 {
		t.Fatalf("Expected count to be 0")
	}
	if ll.head != nil {
		t.Fatalf("Expected head to be nil")
	}
	if ll.tail != nil {
		t.Fatalf("Expected tail to be nil")
	}
}

func TestAddHeadFirst(t *testing.T) {
	ll := NewLinkedList(3)
	new := Node{id: 1, value: "one"}
	err := ll.AddHead(&new)
	if err != nil {
		t.Fatalf("Expected new head got NIL")
	}
	if ll.Size() != 3 {
		t.Fatalf("Expected size to remain 3")
	}
	if ll.Count() != 1 {
		t.Fatalf("Expected count to be 1")
	}
	if ll.head != &new {
		t.Fatalf("Expected head to be the new node")
	}
	if ll.tail != &new {
		t.Fatalf("Expected tail to be the new node")
	}
}

func TestAddHeadMultiple(t *testing.T) {
	ll := NewLinkedList(3)
	new1 := Node{id: 1, value: "one"}
	err := ll.AddHead(&new1)
	if err != nil {
		t.Fatalf("Expected new head[%d] got NIL", new1.id)
	}
	if ll.head != ll.tail {
		t.Fatalf("Expected first node should be both head[%v] and tail[%v]", ll.head, ll.tail)
	}
	new2 := Node{id: 2, value: "two"}
	err = ll.AddHead(&new2)
	if err != nil {
		t.Fatalf("Expected new head[%d] got NIL", new2.id)
	}
	if ll.head == ll.tail {
		t.Fatalf("Expected head[%v] to differ from tail[%v]", ll.head, ll.tail)
	}
	if ll.head != &new2 {
		t.Fatalf("Expected head[%d] to be the new node", new2.id)
	}
	new3 := Node{id: 3, value: "three"}
	err = ll.AddHead(&new3)
	if err != nil {
		t.Fatalf("Expected new head[%d] got NIL", new3.id)
	}
	if ll.head == ll.tail {
		t.Fatalf("Expected head[%v] to differ from tail[%v]", ll.head, ll.tail)
	}
	if ll.head != &new3 {
		t.Fatalf("Expected head[%d] to be the new node", new3.id)
	}
}
func TestBump(t *testing.T) {
	ll := NewLinkedList(3)
	node1 := &Node{
		id:    1,
		value: "one",
	}
	_ = ll.AddHead(node1)
	node2 := &Node{
		id:    2,
		value: "two",
	}
	_ = ll.AddHead(node2)
	node3 := &Node{
		id:    3,
		value: "three",
	}
	_ = ll.AddHead(node3)
	if ll.head != node3 {
		t.Fatalf("Expected head to be node[%d], but was node[%d]", node3.id, ll.head.id)
	}
	if ll.tail != node1 {
		t.Fatalf("Expected tail to be node[%d], but was node[%d]", node1.id, ll.tail.id)
	}
	dropIds, err := ll.Bump(node1)
	if err != nil {
		t.Fatalf("Expected no error, %s", err.Error())
	}
	if len(dropIds) != 0 {
		t.Fatalf("Expected no dropped nodes, but got [%v]", dropIds)
	}
	if ll.head != node1 {
		t.Fatalf("Expected head to be node[%d], but was node[%d]", node2.id, ll.head.id)
	}
	if ll.tail != node2 {
		t.Fatalf("Expected tail to be node[%d], but was node[%d]", node1.id, ll.tail.id)
	}
}

func TestCrop(t *testing.T) {
	ll := NewLinkedList(5)
	_ = ll.AddHead(&Node{id: 1, value: "111"})
	_ = ll.AddHead(&Node{id: 2, value: "222"})
	_ = ll.AddHead(&Node{id: 3, value: "333"})
	_ = ll.AddHead(&Node{id: 4, value: "444"})
	_ = ll.AddHead(&Node{id: 5, value: "555"})
	if ll.Count() != 5 {
		t.Fatalf("Expected count to be 5")
	}
	pruneCount, pruneIds, err := ll.Prune(0)
	if err != nil {
		t.Fatalf("Expected no error, %s", err.Error())
	}
	if pruneCount != 1 {
		t.Fatalf("Expected 1 node to be cropped, but got %d", pruneCount)
	}
	if pruneIds[0] != 1 {
		t.Fatalf("Expected nodes to be cropped [1, 2], but got [%v]", pruneIds)
	}
	if ll.Count() != 4 {
		t.Fatalf("Expected count to be 4, but was [%d]", ll.Count())
	}
	var test *Node
	test, err = ll.Get(1)
	if err == nil || test != nil {
		t.Fatalf("Expected node[%d] to be removed", 1)
	}
	test, err = ll.Get(2)
	if err != nil || test == nil {
		t.Fatalf("Expected node[%d] to remain", 2)
	}
	test, err = ll.Get(3)
	if err != nil || test == nil {
		t.Fatalf("Expected node[%d] to remain, %v", 3, err)
	}
	test, err = ll.Get(4)
	if err != nil || test == nil {
		t.Fatalf("Expected node[%d] to remain, %v", 4, err)
	}
	test, err = ll.Get(5)
	if err != nil || test == nil {
		t.Fatalf("Expected node[%d] to remain, %v", 5, err)
	}
}

func TestRemove(t *testing.T) {
	ll := NewLinkedList(3)
	node1 := &Node{id: 1, value: "one"}
	_ = ll.AddHead(node1)
	node2 := &Node{id: 2, value: "two"}
	_ = ll.AddHead(node2)
	node3 := &Node{id: 3, value: "three"}
	_ = ll.AddHead(node3)
	test := ll.Remove(4)
	if test != nil {
		t.Fatalf("Expected nothing to drop")
	}
	test = ll.Remove(node2.id)
	if test == nil {
		t.Fatalf("Expected id[%d] to drop but got NIL", node2.id)
	}
	if test.id != node2.id {
		t.Fatalf("Unexpected id[%d] dropped instead of [%d]", node2.id, test.id)
	}
	if node1.next != node3.prev {
		t.Fatalf("Expected prev[%d] to point to next[%d]", node1.id, node3.id)
	}
	test = ll.Remove(node3.id)
	if test == nil {
		t.Fatalf("Expected id[%d] to drop but got NIL", node3.id)
	}
	if test.id != node3.id {
		t.Fatalf("Unexpected id[%d] dropped instead of [%d]", node3.id, test.id)
	}
	test = ll.Remove(node1.id)
	if test == nil {
		t.Fatalf("Expected id[%d] to drop but got NIL", node1.id)
	}
	if test.id != node1.id {
		t.Fatalf("Unexpected id[%d] dropped instead of [%d]", node1.id, test.id)
	}
	if ll.Count() != 0 {
		t.Fatalf("Expected count to be 0")
	}
}

func TestRemoveStruct(t *testing.T) {
	type TestyTest struct {
		a int
		b string
		c [4]string
	}

	ll := NewLinkedList(3)
	node1 := &Node{
		id:    1,
		value: TestyTest{1, "one", [4]string{"a", "b", "c", "d"}},
		prev:  &Node{},
		next:  &Node{},
	}
	_ = ll.AddHead(node1)
	node2 := &Node{
		id:    2,
		value: TestyTest{2, "two", [4]string{"k", "l", "m", "n"}},
		prev:  &Node{},
		next:  &Node{},
	}
	_ = ll.AddHead(node2)
	node3 := &Node{
		id:    3,
		value: TestyTest{3, "three", [4]string{"w", "x", "y", "z"}},
		prev:  &Node{},
		next:  &Node{},
	}
	_ = ll.AddHead(node3)
	test := ll.Remove(4)
	if test != nil {
		t.Fatalf("Expected nothing to drop")
	}
	test = ll.Remove(node2.id)
	if test == nil {
		t.Fatalf("Expected id[%d] to drop but got NIL", node2.id)
	}
	if test.value.(TestyTest).a != int(node2.id) {
		t.Fatalf("Unexpected id[%d] dropped instead of [%d]", node2.id, test.id)
	}
	if test.id != node2.id {
		t.Fatalf("Unexpected id[%d] dropped instead of [%d]", node2.id, test.id)
	}
	if node1.next != node3.prev {
		t.Fatalf("Expected prev[%d] to point to next[%d]", node1.id, node3.id)
	}
	test = ll.Remove(node3.id)
	if test == nil {
		t.Fatalf("Expected id[%d] to drop but got NIL", node3.id)
	}
	if test.value.(TestyTest).a != int(node3.id) {
		t.Fatalf("Unexpected id[%d] dropped instead of [%d]", node3.id, test.id)
	}
	if test.value.(TestyTest).b != node3.value.(TestyTest).b {
		t.Fatalf("Unexpected id[%d] dropped instead of [%d]", node3.id, test.id)
	}
	if test.id != node3.id {
		t.Fatalf("Unexpected id[%d] dropped instead of [%d]", node3.id, test.id)
	}
	test = ll.Remove(node1.id)
	if test == nil {
		t.Fatalf("Expected id[%d] to drop but got NIL", node1.id)
	}
	if test.value.(TestyTest).a != int(node1.id) {
		t.Fatalf("Unexpected id[%d] dropped instead of [%d]", node1.id, test.id)
	}
	if test.value.(TestyTest).b != node1.value.(TestyTest).b {
		t.Fatalf("Unexpected id[%d] dropped instead of [%d]", node1.id, test.id)
	}
	if test.id != node1.id {
		t.Fatalf("Unexpected id[%d] dropped instead of [%d]", node1.id, test.id)
	}
	if ll.Count() != 0 {
		t.Fatalf("Expected count to be 0")
	}
}

func TestRemoveHead(t *testing.T) {
	ll := NewLinkedList(3)
	node1 := &Node{id: 1, value: "one"}
	_ = ll.AddHead(node1)
	node2 := &Node{id: 2, value: "two"}
	_ = ll.AddHead(node2)
	node3 := &Node{id: 3, value: "three"}
	_ = ll.AddHead(node3)
	test := ll.removeHead()
	if test == nil {
		t.Fatalf("Expected to drop head but got NIL")
	}
	if test.id != node3.id {
		t.Fatalf("Expected to drop head[%d] but got [%d]", node3.id, test.id)
	}
	if test.id == node2.id || test.id == node1.id {
		t.Fatalf("Expected to drop head[%d] but got [%d]", node3.id, test.id)
	}
}

func TestRemoveTail(t *testing.T) {
	ll := NewLinkedList(3)
	node1 := &Node{id: 1, value: "one"}
	_ = ll.AddHead(node1)
	node2 := &Node{id: 2, value: "two"}
	_ = ll.AddHead(node2)
	node3 := &Node{id: 3, value: "three"}
	_ = ll.AddHead(node3)
	test := ll.removeTail()
	if test == nil {
		t.Fatalf("Expected to drop tail but got NIL")
	}
	if test.id == ll.tail.id {
		t.Fatalf("Expected still points to removed[%d]", test.id)
	}
	if test.id != node1.id {
		t.Fatalf("Expected to drop tail[%d] but got [%d]", node1.id, test.id)
	}
	test = ll.removeTail()
	if test == nil {
		t.Fatalf("Expected to drop tail but got NIL")
	}
	if test.id == ll.tail.id {
		t.Fatalf("Expected still points to removed")
	}
	if test.id != node2.id {
		t.Fatalf("Expected to drop tail[%d] but got [%d]", node2.id, test.id)
	}
	test = ll.removeTail()
	if test == nil {
		t.Fatalf("Expected to remove last but got NIL")
	}
	if ll.tail != nil {
		t.Fatalf("Expected still points to [%d]", ll.tail.id)
	}
	if test.id != node3.id {
		t.Fatalf("Expected to drop tail[%d] but got [%d]", node3.id, test.id)
	}
	test = ll.removeTail()
	if test != nil {
		t.Fatalf("Expected should have received NIL but got [%v]", test)
	}
}
