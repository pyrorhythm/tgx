package fsm

import (
	"sync"

	"pyrorhythm.dev/fn/res"
)

// dataStorage is a type for default data storage
type dataStorage[K comparable, V any] struct {
	mu      sync.Mutex
	Storage map[int64]map[K]V
}

// initialDataStorage creates in memory storage for user's data
func initialDataStorage[K comparable, V any]() *dataStorage[K, V] {
	return &dataStorage[K, V]{
		Storage: make(map[int64]map[K]V),
	}
}

// Set sets user's data to data storage
func (d *dataStorage[K, V]) Set(userID int64, key K, value V) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	s, ok := d.Storage[userID]
	if !ok {
		s = make(map[K]V)
		d.Storage[userID] = s
	}

	s[key] = value

	return nil
}

// Get gets user's data from data storage
func (d *dataStorage[K, V]) Get(userID int64, key K) res.Of[V] {
	d.mu.Lock()
	defer d.mu.Unlock()

	m, ok := d.Storage[userID]
	if !ok {
		return res.Err[V](ErrNoUserData)
	}

	v, ok := m[key]
	if !ok {
		return res.Err[V](ErrNoUserData)
	}

	return res.OKAny(v)
}

// Delete deletes user's data from data storage
func (d *dataStorage[K, V]) Delete(userID int64, key K) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.Storage[userID]; !ok {
		return nil
	}

	delete(d.Storage[userID], key)

	return nil
}
