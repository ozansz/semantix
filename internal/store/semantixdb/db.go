package semantixdb

import (
	"github.com/ozansz/semantix/internal/parser"
	"github.com/ozansz/semantix/internal/store"
)

// TODO: Implement this
type DB struct {
	strings *StringsDB
	facts   *FactsDB
}

// TODO: Implement this
func (db *DB) Add(t *parser.Fact) error {
	return nil
}

// TODO: Implement this
func (db *DB) Get(q *store.Query) (map[uint32]*parser.Fact, error) {
	return nil, nil
}
