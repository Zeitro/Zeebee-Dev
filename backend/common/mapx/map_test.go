package mapx

import (
	"sync"
	"testing"
)

func TestNewMap(t *testing.T) {
	m := NewMap[string, int]()
	if m == nil {
		t.Fatal("NewMap returned nil")
	}
}

func TestMap_SetAndGet(t *testing.T) {
	m := NewMap[string, int]()

	// Test Set and Get existing key
	m.Set("key1", 100)
	val, ok := m.Get("key1")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if val != 100 {
		t.Errorf("expected value 100, got %d", val)
	}

	// Test Get non-existent key
	_, ok = m.Get("nonexistent")
	if ok {
		t.Error("expected key to not exist")
	}
}

func TestMap_Overwrite(t *testing.T) {
	m := NewMap[string, int]()

	m.Set("key1", 100)
	m.Set("key1", 200)

	val, ok := m.Get("key1")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if val != 200 {
		t.Errorf("expected overwritten value 200, got %d", val)
	}
}

func TestMap_Delete(t *testing.T) {
	m := NewMap[string, int]()

	// Set and then delete
	m.Set("key1", 100)
	m.Delete("key1")

	_, ok := m.Get("key1")
	if ok {
		t.Error("expected key to be deleted")
	}

	// Delete non-existent key should not panic
	m.Delete("nonexistent")
}

func TestMap_Clear(t *testing.T) {
	m := NewMap[string, int]()

	m.Set("key1", 100)
	m.Set("key2", 200)
	m.Set("key3", 300)

	m.Clear()

	// All keys should be gone after Clear
	val, ok := m.Get("key1")
	if ok {
		t.Error("expected key1 to be cleared")
	}
	if val != 0 {
		t.Errorf("expected zero value, got %d", val)
	}

	_, ok = m.Get("key2")
	if ok {
		t.Error("expected key2 to be cleared")
	}

	_, ok = m.Get("key3")
	if ok {
		t.Error("expected key3 to be cleared")
	}
}

func TestMap_ClearOnEmpty(t *testing.T) {
	m := NewMap[string, int]()

	// Clear on empty map should not panic
	m.Clear()

	// Should still work after clear
	m.Set("key1", 100)
	val, ok := m.Get("key1")
	if !ok || val != 100 {
		t.Error("map should work after Clear on empty")
	}
}

func TestMap_ConcurrentReads(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("key1", 100)
	m.Set("key2", 200)
	m.Set("key3", 300)

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			m.Get("key1")
			m.Get("key2")
			m.Get("key3")
		}()
	}

	wg.Wait()
}

func TestMap_ConcurrentWrites(t *testing.T) {
	m := NewMap[string, int]()

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(n int) {
			defer wg.Done()
			key := "key"
			m.Set(key, n)
		}(i)
	}

	wg.Wait()

	// Verify final value exists
	val, ok := m.Get("key")
	if !ok {
		t.Fatal("expected key to exist")
	}
	// We don't check exact value due to race condition,
	// but it should be one of the written values
	if val < 0 || val >= goroutines {
		t.Errorf("unexpected value: %d", val)
	}
}

func TestMap_ConcurrentMixed(t *testing.T) {
	m := NewMap[string, int]()

	// Pre-populate
	m.Set("key1", 100)
	m.Set("key2", 200)

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines * 3)

	for i := 0; i < goroutines; i++ {
		// Concurrent reads
		go func() {
			defer wg.Done()
			m.Get("key1")
		}()
		// Concurrent writes
		go func(n int) {
			defer wg.Done()
			m.Set("key3", n)
		}(i)
		// Concurrent deletes
		go func() {
			defer wg.Done()
			m.Delete("key2")
		}()
	}

	wg.Wait()
}

func TestMap_WithIntKey(t *testing.T) {
	m := NewMap[int, string]()

	m.Set(1, "one")
	m.Set(2, "two")
	m.Set(999, "nine nine nine")

	val, ok := m.Get(1)
	if !ok || val != "one" {
		t.Error("expected to get 'one' for key 1")
	}

	val, ok = m.Get(999)
	if !ok || val != "nine nine nine" {
		t.Error("expected to get 'nine nine nine' for key 999")
	}
}

func TestMap_WithStructValue(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	m := NewMap[string, Person]()

	m.Set("p1", Person{Name: "Alice", Age: 30})
	m.Set("p2", Person{Name: "Bob", Age: 25})

	val, ok := m.Get("p1")
	if !ok {
		t.Fatal("expected p1 to exist")
	}
	if val.Name != "Alice" || val.Age != 30 {
		t.Errorf("expected Alice 30, got %s %d", val.Name, val.Age)
	}
}

func TestMap_ZeroValue(t *testing.T) {
	m := NewMap[string, int]()

	m.Set("zero", 0)

	val, ok := m.Get("zero")
	if !ok {
		t.Error("expected zero value key to exist")
	}
	if val != 0 {
		t.Errorf("expected 0, got %d", val)
	}
}

func TestMap_MultipleOperations(t *testing.T) {
	m := NewMap[string, int]()

	// Set multiple values
	keys := []string{"a", "b", "c", "d", "e"}
	for i, key := range keys {
		m.Set(key, i*10)
	}

	// Verify all values
	for i, key := range keys {
		val, ok := m.Get(key)
		if !ok {
			t.Errorf("key %s should exist", key)
			continue
		}
		if val != i*10 {
			t.Errorf("key %s: expected %d, got %d", key, i*10, val)
		}
	}

	// Delete some keys
	m.Delete("b")
	m.Delete("d")

	// Verify deletions
	if _, ok := m.Get("b"); ok {
		t.Error("b should be deleted")
	}
	if _, ok := m.Get("d"); ok {
		t.Error("d should be deleted")
	}

	// Remaining keys should still exist
	if _, ok := m.Get("a"); !ok {
		t.Error("a should still exist")
	}
	if _, ok := m.Get("c"); !ok {
		t.Error("c should still exist")
	}
	if _, ok := m.Get("e"); !ok {
		t.Error("e should still exist")
	}
}
