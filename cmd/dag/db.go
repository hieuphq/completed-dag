package main

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// DB interface
type DB interface {
}

// NewDB for level db
func NewDB(path string) (*leveldb.DB, error) {
	return leveldb.OpenFile(path, nil)
}
