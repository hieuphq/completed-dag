package store

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// DB interface
type DB interface {
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte) error
	Delete(key []byte) error
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
