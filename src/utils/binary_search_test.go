package utils

import (
	i "neon-chat/src/interface"
	"testing"
)

type TestItem struct {
	Id uint
}

func (ti TestItem) GetId() uint {
	return ti.Id
}

// running this test we get
// 2024/08/04 id[44] isfound:true check took 13.375µs, 14 loops in 90910 items
// 2024/08/04 id[77] isfound:true check took 125ns, 15 loops in 90910 items
// 2024/08/04 id[99] isfound:true check took 167ns, 17 loops in 90910 items
// 2024/08/04 id[43] isfound:false check took 125ns, 17 loops in 90910 items
// 2024/08/04 id[75] isfound:false check took 167ns, 17 loops in 90910 items
// 2024/08/04 id[91] isfound:false check took 125ns, 17 loops in 90910 items

func TestBinarySearch(t *testing.T) {
	items := make([]i.Identifiable, 0)
	for i := 0; i < 1_000_000; i += 11 {
		item := TestItem{Id: uint(i)}
		items = append(items, item)
	}
	var found i.Identifiable
	var idx int
	found, idx = BinarySearch(items, 0)
	if found == nil {
		t.Error("expected 0, got nil")
	}
	if idx != 0 {
		t.Errorf("expected index 0, got %d", idx)
	}
	if found.GetId() != 0 {
		t.Errorf("expected 0, got %d", found.GetId())
	}
	found, idx = BinarySearch(items, 77)
	if found == nil {
		t.Error("expected 77, got nil")
	}
	if idx != 7 {
		t.Errorf("expected index 7, got %d", idx)
	}
	if found.GetId() != 77 {
		t.Errorf("expected 77, got %d", found.GetId())
	}
	found, idx = BinarySearch(items, 99)
	if found == nil {
		t.Error("expected 99, got nil")
	}
	if idx != 9 {
		t.Errorf("expected index 9, got %d", idx)
	}
	if found.GetId() != 99 {
		t.Errorf("expected 99, got %d", found.GetId())
	}
	notFound, idx := BinarySearch(items, 43)
	if notFound != nil {
		t.Errorf("should not have found %d", notFound)
	}
	if idx != -1 {
		t.Errorf("expected index -1, got %d", idx)
	}
	notFound, idx = BinarySearch(items, 75)
	if notFound != nil {
		t.Errorf("should not have found %d", notFound)
	}
	if idx != -1 {
		t.Errorf("expected index -1, got %d", idx)
	}
	notFound, idx = BinarySearch(items, 91)
	if notFound != nil {
		t.Errorf("should not have found %d", notFound)
	}
	if idx != -1 {
		t.Errorf("expected index -1, got %d", idx)
	}
}
