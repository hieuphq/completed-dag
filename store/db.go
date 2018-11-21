package store

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// DB interface
type DB interface {
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte) error
	Delete(key []byte) error
	Transfer(db DB) error
	All() (map[string][]byte, error)
	Close() error
}

type levelDBImpl struct {
	DB *leveldb.DB
}

// NewDB for level db
func NewDB(path string) (DB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &levelDBImpl{
		DB: db,
	}, nil
}

func (r *levelDBImpl) Get(key []byte) ([]byte, error) {
	return r.DB.Get(key, nil)
}

func (r *levelDBImpl) Put(key []byte, value []byte) error {
	return r.DB.Put(key, value, nil)
}

func (r *levelDBImpl) Delete(key []byte) error {
	return r.DB.Delete(key, nil)
}

func (r *levelDBImpl) Close() error {
	return r.DB.Close()
}

func (r *levelDBImpl) Transfer(db DB) error {
	all, err := r.All()

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

func (r *levelDBImpl) All() (map[string][]byte, error) {
	rs := map[string][]byte{}
	iter := r.DB.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		rs[string(key)] = value
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return nil, err
	}

	return rs, nil
}
