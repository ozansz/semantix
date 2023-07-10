package semantixdb

import "github.com/oklog/ulid/v2"

// TODO: Implement this
type StringsDB struct {
	// String hash to DB ID index
	// TODO: Consider replacing this with or adding a bloom filter
	index map[string]ulid.ULID
}
