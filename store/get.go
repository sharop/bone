package store

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/timshannon/badgerhold/v4"
	"regexp"
)

// Get Primitive Get function of badger
func (db *Database) Get(key string) (value interface{}) {
	err := db.Service.Badger().View(func(txn *badger.Txn) error {
		item, err := txn.Get(makeKey(key))
		if err != nil {
			return err
		}

		return item.Value(func(valueBytes []byte) error {
			value = string(append([]byte{}, valueBytes...)[:])
			return nil
		})
	})

	if err != nil && err != badger.ErrKeyNotFound {
		return nil
	}

	return
}

func (db *Database) FindBy(field, expression string) ([]Item, error) {
	var result []Item
	err := db.Service.Find(&result, badgerhold.Where(field).RegExp(regexp.MustCompile(expression)))
	if err != nil {
		return nil, err
	}
	return result, err
}

func (db *Database) GetBadgerhold(key string) (*Item, error) {
	var result Item
	err := db.Service.Get(key, &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}
