package store

import (
	"fmt"
	"testing"
)

func TestNewLRUCache(t *testing.T) {
	cache := NewLRUCache(3)
	if cache.Size() != 3 {
		t.Fatalf("Expected size to be 3")
	}
	if cache.Count() != 0 {
		t.Fatalf("Expected count to be 0")
	}
	if cache.dict == nil {
		t.Fatalf("Expected dict to be initialized")
	}
	if cache.list == nil {
		t.Fatalf("Expected list to be initialized")
	}
}

func TestSize(t *testing.T) {
	cache := NewLRUCache(3)
	if cache.Size() != 3 {
		t.Fatalf("Expected size to be 3")
	}
}

func TestCount(t *testing.T) {
	cache := NewLRUCache(3)
	if cache.Count() != 0 {
		t.Fatalf("Expected count to be 0")
	}
	cache.Set(1, "one")
	if cache.Count() != 1 {
		t.Fatalf("Expected count to be 0")
	}
}

func TestSet(t *testing.T) {
	cache := NewLRUCache(3)
	err := cache.Set(1, "one")
	if err != nil {
		t.Fatalf("Expected no error, %s", err.Error())
	}
	if cache.Count() != 1 {
		t.Fatalf("Expected count to be 1")
	}
	storedInMap := cache.dict[1]
	if storedInMap == nil {
		t.Fatalf("Expected storedInMap dict is empty")
	}
	if storedInMap.id != 1 {
		t.Fatalf("Expected storedInMap id to be 1")
	}
	if storedInMap.value != "one" {
		t.Fatalf("Expected storedInMap value to be 'one'")
	}
	storedInList, err := cache.list.Get(1)
	if err != nil {
		t.Fatalf("Expected node to be in list, %s", err.Error())
	}
	if storedInList.id != 1 {
		t.Fatalf("Expected storedInList id to be 1")
	}
	if storedInList.value != "one" {
		t.Fatalf("Expected value to be 'one'")
	}
	if storedInList != storedInMap {
		t.Fatalf("Expected storedInList to be storedInMap")
	}
}

func TestSetMultiple(t *testing.T) {
	cache := NewLRUCache(3)
	err := cache.Set(1, "one")
	if err != nil {
		t.Fatalf("Expected no error 1, %s", err.Error())
	}
	if cache.Count() != 1 {
		t.Fatalf("Expected count to be 1")
	}
	err = cache.Set(2, "two")
	if err != nil {
		t.Fatalf("Expected no error 2, %s", err.Error())
	}
	if cache.Count() != 2 {
		t.Fatalf("Expected count to be 2")
	}
	err = cache.Set(3, "three")
	if err != nil {
		t.Fatalf("Expected no error 3, %s", err.Error())
	}
	if cache.Count() != 3 {
		t.Fatalf("Expected count to be 3")
	}
	err = cache.Set(4, "four")
	if err != nil {
		t.Fatalf("Expected 4 to be added, 1 removed")
	}
	if cache.Count() != 3 {
		t.Fatalf("Expected count to remain 3 but was [%d]", cache.Count())
	}
	_, err = cache.Get(1)
	if err == nil {
		t.Fatalf("Expected 1 to be removed")
	}
	if cache.Count() != 3 {
		t.Fatalf("Expected count to remain 3")
	}
}

func TestGet(t *testing.T) {
	cache := NewLRUCache(3)
	_ = cache.Set(1, "one")
	_ = cache.Set(2, "two")
	_ = cache.Set(3, "three")
	if cache.list.head.id != 3 {
		t.Fatalf("Expected head to be 3")
	}
	if cache.list.tail.id != 1 {
		t.Fatalf("Expected tail to be 1")
	}
	value, err := cache.Get(2)
	if err != nil {
		t.Fatalf("Unexpected error on 2, %s", err.Error())
	}
	if value != "two" {
		t.Fatalf("Expected value to be 'two', but was [%s]", value)
	}
	if cache.list.head.id != 2 {
		t.Fatalf("Expected head to be 3")
	}
	if cache.list.tail.id != 1 {
		t.Fatalf("Expected tail to be 1")
	}
	value, err = cache.Get(4)
	if err == nil {
		t.Fatalf("Expected no key 4 in cache")
	}
	if value != "" {
		t.Fatalf("Expected value to be empty string, but was [%s]", value)
	}
}

func TestScan(t *testing.T) {
	cache := NewLRUCache(3)
	_ = cache.Set(1, "one")
	_ = cache.Set(2, "two")
	_ = cache.Set(3, "three")
	keys := cache.Keys()
	if len(keys) != 3 {
		t.Fatalf("Expected keys to be 3 but was [%d]", len(keys))
	}
	i1 := false
	i2 := false
	i3 := false
	for i := 0; i < 3; i++ {
		switch keys[i] {
		case 1:
			i1 = true
		case 2:
			i2 = true
		case 3:
			i3 = true
		}
	}
	if !i1 || !i2 || !i3 {
		t.Fatalf("Unexpected set of keys %v", keys)
	}
}

func TestTake(t *testing.T) {
	cache := NewLRUCache(3)
	_ = cache.Set(1, "one")
	_ = cache.Set(2, "two")
	_ = cache.Set(3, "three")
	value, err := cache.Get(2)
	if err != nil {
		t.Fatalf("Unexpected error on 2, %s", err.Error())
	}
	if value != "two" {
		t.Fatalf("Expected value to be 'two', but was [%s]", value)
	}
	removed, err := cache.Take(2)
	if err != nil {
		t.Fatalf("Unexpected error on 2, %s", err.Error())
	}
	if removed == nil {
		t.Fatalf("Expected removed to be 2, but was NIL")
	}
	removedValue, ok := removed.(string)
	if !ok {
		t.Fatalf("Expected removed to be string, but was [%+v]", removed)
	}
	if removedValue != "two" {
		t.Fatalf("Expected removed to be 'two', but was [%s]", removedValue)
	}
}

func TestDrop(t *testing.T) {
	cache := NewLRUCache(8)
	var i uint
	for i = 0; i < 8; i++ {
		_ = cache.Set(uint(i), fmt.Sprintf("value-%d", i))
	}
	if cache.Count() != 8 {
		t.Fatalf("Expected count to be 8")
	}
	n, err := cache.Drop(4)
	if err != nil {
		t.Fatalf("Unexpected error on drop, %s", err.Error())
	}
	if n != 4 {
		t.Fatalf("Expected to drop 4, but dropped %d", n)
	}
	if cache.Count() != 4 {
		t.Fatalf("Expected count to be 4")
	}
	if cache.Size() != 8 {
		t.Fatalf("Expected size to remain 8")
	}
	tmp, err := cache.Get(2)
	if err == nil {
		t.Fatalf("Unexpected 2 to be dropped, but found[%v]", tmp)
	}
	s, ok := tmp.(string)
	if !ok {
		t.Fatalf("Expected value to be string, but was [%+v]", tmp)
	}
	if s != "" {
		t.Fatalf("Expected value to be 'value-2', but was [%s]", s)
	}
	tmp, err = cache.Get(5)
	if err != nil {
		t.Fatalf("Unexpected 5 to remain, but it wasn't")
	}
	s, ok = tmp.(string)
	if !ok {
		t.Fatalf("Expected value of 5 to be string, but was [%+v]", tmp)
	}
	if s != "value-5" {
		t.Fatalf("Expected value to be [value-5], but was [%s]", s)
	}
}
