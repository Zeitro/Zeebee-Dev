package mapx

import "sync"

type Map[K comparable, V any] struct {
	m map[K]V
	sync.RWMutex
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		m: make(map[K]V),
	}
}

// Get get data by key
func (m *Map[K, V]) Get(key K) (V, bool) {
	m.RLock()
	defer m.RUnlock()
	val, ok := m.m[key]
	return val, ok
}

// Set set data by key
func (m *Map[K, V]) Set(key K, val V) {
	m.Lock()
	defer m.Unlock()
	m.m[key] = val
}

// Delete delete data by key
func (m *Map[K, V]) Delete(key K) {
	m.Lock()
	defer m.Unlock()
	delete(m.m, key)
}

// Clear  clear all  data
func (m *Map[K, V]) Clear() {
	m.Lock()
	defer m.Unlock()
	m.m = make(map[K]V)
}
