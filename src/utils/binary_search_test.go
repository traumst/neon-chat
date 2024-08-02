package utils

import (
	"prplchat/src/utils/interfaces"
	"testing"
)

type TestItem struct {
	Id uint
}

func (ti TestItem) GetId() uint {
	return ti.Id
}

func TestBinarySearch(t *testing.T) {
	items := make([]interfaces.Identifiable, 0)
	for i := 11; i < 100; i += 11 {
		item := TestItem{Id: uint(i)}
		items = append(items, item)
	}
	found := BinarySearch(items, 44)
	if found == nil {
		t.Error("expected 44, got nil")
	}
	if found.GetId() != 44 {
		t.Errorf("expected 44, got %d", found.GetId())
	}
	found = BinarySearch(items, 77)
	if found == nil {
		t.Error("expected 77, got nil")
	}
	if found.GetId() != 77 {
		t.Errorf("expected 77, got %d", found.GetId())
	}
	found = BinarySearch(items, 99)
	if found == nil {
		t.Error("expected 99, got nil")
	}
	if found.GetId() != 99 {
		t.Errorf("expected 99, got %d", found.GetId())
	}
	found = BinarySearch(items, 43)
	if found != nil {
		t.Errorf("should not have found %d", found)
	}
	found = BinarySearch(items, 75)
	if found != nil {
		t.Errorf("should not have found %d", found)
	}
	found = BinarySearch(items, 91)
	if found != nil {
		t.Errorf("should not have found %d", found)
	}
}
