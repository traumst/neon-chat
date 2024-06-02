package store

import "testing"

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

func TestAdd(t *testing.T) {
	ll := NewLinkedList(3)
	node, err := ll.Add(1, "one")
	if err != nil {
		t.Fatalf("Expected no error, %s", err.Error())
	}
	if node.id != 1 {
		t.Fatalf("Expected id to be 1")
	}
	if node.value != "one" {
		t.Fatalf("Expected value to be 'one'")
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

func TestMultipleAdd(t *testing.T) {
	ll := NewLinkedList(3)
	var node *Node
	var err error
	node, err = ll.Add(1, "one")
	if err != nil {
		t.Fatalf("Expected no error on 1, node[%+v] err[%s]", node, err)
	}
	node, err = ll.Add(2, "one")
	if err != nil {
		t.Fatalf("Expected no error on 2, node[%+v] err[%s]", node, err)
	}
	node, err = ll.Add(3, "three")
	if err != nil {
		t.Fatalf("Expected no error on 3, node[%+v] err[%s]", node, err)
	}
	node, err = ll.Add(4, "four")
	if err == nil || node != nil {
		t.Fatalf("Expected error on 4, node[%+v] err[%s]", node, err)
	}
	_, err = ll.Add(5, "five")
	if err == nil || node != nil {
		t.Fatalf("Expected error on 5, node[%+v] err[%s]", node, err)
	}
}

func TestBump(t *testing.T) {
	ll := NewLinkedList(3)
	node1, _ := ll.Add(1, "one")
	node2, _ := ll.Add(2, "two")
	node3, _ := ll.Add(3, "three")
	if ll.head != node1 {
		t.Fatalf("Expected head to be node[%d], but was node[%d]", node1.id, ll.head.id)
	}
	if ll.tail != node3 {
		t.Fatalf("Expected tail to be node[%d], but was node[%d]", node3.id, ll.tail.id)
	}
	err := ll.Bump(node1)
	if err != nil {
		t.Fatalf("Expected no error, %s", err.Error())
	}
	if ll.head != node2 {
		t.Fatalf("Expected head to be node[%d], but was node[%d]", node2.id, ll.head.id)
	}
	if ll.tail != node1 {
		t.Fatalf("Expected tail to be node[%d], but was node[%d]", node1.id, ll.tail.id)
	}
}

func TestCrop(t *testing.T) {
	ll := NewLinkedList(5)
	ll.Add(1, "one")
	ll.Add(2, "two")
	ll.Add(3, "three")
	ll.Add(4, "four")
	ll.Add(5, "five")
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
	node1, _ := ll.Add(1, "one")
	node2, _ := ll.Add(2, "two")
	node3, _ := ll.Add(3, "three")
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
}
