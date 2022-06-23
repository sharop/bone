package log

import (
	badger "github.com/sharop/bone/store"
	"sync"
)

// COmuncation with bleve
type indexer struct {
	path string

	db *badger.Database

	mutex sync.Mutex

	closed bool
}

func (i *indexer) ExistKeyWith(prefix []byte, neq []byte) (bool, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if i.closed {
		return false, badger.ErrDBClosed
	}

	//Implement search for index prefix
	return false, nil
}