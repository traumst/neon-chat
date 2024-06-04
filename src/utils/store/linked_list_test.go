package store

import (
	"testing"
)

func TestNewLinkedList(t *testing.T) {
	ll := NewLinkedList(3)
	if ll.size != 3 {
		t.Fatalf("Expected size to be 3")
	}
	if ll.count != 0 {
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
	if ll.size != 3 {
		t.Fatalf("Expected size to remain 3")
	}
	if ll.count != 1 {
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

func TestAddTail(t *testing.T) {
	ll := NewLinkedList(3)
	newTail := "one"
	node, err := ll.AddTail(1, newTail)
	if err != nil {
		t.Fatalf("Expected no error, %s", err.Error())
	}
	if node.id != 1 {
		t.Fatalf("Expected id to be 1")
	}
	if node.value != newTail {
		t.Fatalf("Expected head to be the new node")
	}
	if ll.head != node {
		t.Fatalf("Expected head to be the new node")
	}
	if ll.tail != node {
		t.Fatalf("Expected tail to be the new node")
	}
	if ll.size != 3 {
		t.Fatalf("Expected size to remain 3")
	}
	if ll.count != 1 {
		t.Fatalf("Expected count to be 1")
	}
	if ll.head != node {
		t.Fatalf("Expected head to be the new node")
	}
	if ll.tail != node {
		t.Fatalf("Expected tail to be the new node")
	}
}

func TestMultipleAddTail(t *testing.T) {
	ll := NewLinkedList(3)
	var node *Node
	var err error
	node, err = ll.AddTail(1, "one")
	if err != nil {
		t.Fatalf("Expected no error on 1, node[%+v] err[%s]", node, err)
	}
	node, err = ll.AddTail(2, "one")
	if err != nil {
		t.Fatalf("Expected no error on 2, node[%+v] err[%s]", node, err)
	}
	node, err = ll.AddTail(3, "three")
	if err != nil {
		t.Fatalf("Expected no error on 3, node[%+v] err[%s]", node, err)
	}
	node, err = ll.AddTail(4, "four")
	if err == nil || node != nil {
		t.Fatalf("Expected error on 4, node[%+v] err[%s]", node, err)
	}
	_, err = ll.AddTail(5, "five")
	if err == nil || node != nil {
		t.Fatalf("Expected error on 5, node[%+v] err[%s]", node, err)
	}
}

func TestBump(t *testing.T) {
	ll := NewLinkedList(3)
	node1, _ := ll.AddTail(1, "one")
	node2, _ := ll.AddTail(2, "two")
	node3, _ := ll.AddTail(3, "three")
	if ll.head != node1 {
		t.Fatalf("Expected head to be node[%d], but was node[%d]", node1.id, ll.head.id)
	}
	if ll.tail != node3 {
		t.Fatalf("Expected tail to be node[%d], but was node[%d]", node3.id, ll.tail.id)
	}
	err := ll.Bump(node3)
	if err != nil {
		t.Fatalf("Expected no error, %s", err.Error())
	}
	if ll.head != node3 {
		t.Fatalf("Expected head to be node[%d], but was node[%d]", node2.id, ll.head.id)
	}
	if ll.tail != node2 {
		t.Fatalf("Expected tail to be node[%d], but was node[%d]", node1.id, ll.tail.id)
	}
}

func TestCrop(t *testing.T) {
	ll := NewLinkedList(5)
	ll.AddTail(1, "one")
	ll.AddTail(2, "two")
	ll.AddTail(3, "three")
	ll.AddTail(4, "four")
	ll.AddTail(5, "five")
	if ll.count != 5 {
		t.Fatalf("Expected count to be 5")
	}
	cropped, err := ll.Crop(2)
	if err != nil {
		t.Fatalf("Expected no error, %s", err.Error())
	}
	if cropped != 2 {
		t.Fatalf("Expected 2 nodes to be cropped, but got %d", cropped)
	}
	if ll.count != 3 {
		t.Fatalf("Expected count to be 3")
	}
}

func TestRemove(t *testing.T) {
	ll := NewLinkedList(3)
	node1, _ := ll.AddTail(1, "one")
	node2, _ := ll.AddTail(2, "two")
	node3, _ := ll.AddTail(3, "three")
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
	node1, _ := ll.AddTail(1, TestyTest{1, "one", [4]string{"a", "b", "c", "d"}})
	node2, _ := ll.AddTail(2, TestyTest{2, "two", [4]string{"k", "l", "m", "n"}})
	node3, _ := ll.AddTail(3, TestyTest{3, "three", [4]string{"w", "x", "y", "z"}})
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
	node1, _ := ll.AddTail(1, "one")
	node2, _ := ll.AddTail(2, "two")
	node3, _ := ll.AddTail(3, "three")
	test := ll.removeHead()
	if test == nil {
		t.Fatalf("Expected to drop head but got NIL")
	}
	if test.id != node1.id {
		t.Fatalf("Expected to drop head[%d] but got [%d]", node1.id, test.id)
	}
	if test.id == node2.id || test.id == node3.id {
		t.Fatalf("Expected to drop head[%d] but got [%d]", node1.id, test.id)
	}
}

func TestRemoveTail(t *testing.T) {
	ll := NewLinkedList(3)
	node1, _ := ll.AddTail(1, "one")
	node2, _ := ll.AddTail(2, "two")
	node3, _ := ll.AddTail(3, "three")
	test := ll.removeTail()
	if test == nil {
		t.Fatalf("Expected to drop tail but got NIL")
	}
	if test.id == ll.tail.id {
		t.Fatalf("Expected still points to removed[%d]", test.id)
	}
	if test.id != node3.id {
		t.Fatalf("Expected to drop tail[%d] but got [%d]", node3.id, test.id)
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
	if test.id != node1.id {
		t.Fatalf("Expected to drop tail[%d] but got [%d]", node1.id, test.id)
	}
}
