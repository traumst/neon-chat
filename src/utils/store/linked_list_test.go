package store

import "testing"

func TestNewLinkedList(t *testing.T) {
	ll := NewLinkedList(10)
	if ll.size != 10 {
		t.Fatalf("Expected size to be 10")
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
}

func TestAddMany(t *testing.T) {
	ll := NewLinkedList(3)
	var node *Node
	var err error
	node, err = ll.Add(1, "one")
	if err != nil || node == nil {
		t.Fatalf("Expected no error on 1, node[%+v] err[%s]", node, err)
	}
	node, err = ll.Add(2, "two")
	if err != nil || node == nil {
		t.Fatalf("Expected no error on 2, node[%+v] err[%s]", node, err)
	}
	node, err = ll.Add(3, "three")
	if err != nil || node == nil {
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
