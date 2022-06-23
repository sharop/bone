package store

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
)

// Set Primitive Set function of badger
func (db *Database) Set(key string, value string) error {
	valueBytes := []byte(value)
	return db.Service.Badger().Update(func(txn *badger.Txn) error {
		return txn.SetEntry(badger.NewEntry(makeKey(key), valueBytes))
	},
	)
}

func (db *Database) Insert(item Item) error {
	pathTemporal := item.Path
	key := item.Key
	Key := fmt.Sprintf("%s/%s", pathTemporal, key)
	//db.Service.Insert(Key, item)
	return db.Service.Upsert(Key, item)
}
