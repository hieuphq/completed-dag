package store

import (
	"sync"

	"github.com/hieuphq/completed-dag/errors"
)

type value []byte

// MemoryStore data store
type MemoryStore struct {
	sync.RWMutex // ← this mutex protects the cache below
	cache        map[string]value
}

// NewMemory create a memory store
func NewMemory() DB {
	return &MemoryStore{
		cache: make(map[string]value),
	}
}

func (ds *MemoryStore) set(key string, value []byte) error {
	ds.Lock()
	defer ds.Unlock()
	ds.cache[key] = value

	return nil
}

func (ds *MemoryStore) get(key string) ([]byte, error) {
	ds.RLock()
	defer ds.RUnlock()
	if item, ok := ds.cache[key]; ok {
		return item, nil
	}
	return nil, errors.ErrNil
}

func (ds *MemoryStore) count() int {
	ds.RLock()
	defer ds.RUnlock()
	return len(ds.cache)
}

// Get ...
func (ds *MemoryStore) Get(key []byte) ([]byte, error) {
	return ds.get(string(key))
}

// Put ...
func (ds *MemoryStore) Put(key []byte, value []byte) error {
	return ds.set(string(key), value)
}

// Delete ...
func (ds *MemoryStore) Delete(key []byte) error {
	return ds.set(string(key), nil)
}

// Close ...
func (ds *MemoryStore) Close() error {
	return nil
}

// Transfer data
func (ds *MemoryStore) Transfer(db DB) error {
	all, err := ds.All()

	if err != nil {
		return err
	}

	if len(all) <= 0 {
		return nil
	}

	for key := range all {
		err = db.Put([]byte(key), all[key])
		if err != nil {
			return err
		}

	}

	return nil
}

// All get all data
func (ds *MemoryStore) All() (map[string][]byte, error) {
	rs := map[string][]byte{}

	for key := range ds.cache {
		rs[key] = ds.cache[key]
	}

	return rs, nil
}

// TryToGetData ...
func (ds *MemoryStore) TryToGetData() (bool, error) {
	return len(ds.cache) > 0, nil
}
