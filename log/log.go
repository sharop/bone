package log

import (
	"github.com/google/uuid"
	api "github.com/sharop/bone/pb/v1/log"
	badger "github.com/sharop/bone/store"
	"log"
	"os"
	"sync"
)

type Log struct {
	mu     sync.RWMutex
	Dir    string
	Config Config
	db     *badger.Database
	idx    *indexer
}

func NewLog(dir string, c Config) (*Log, error) {

	l := &Log{
		Dir:    dir,
		Config: c,
	}
	return l, l.setup()
}

func (l *Log) setup() error {
	db, err := badger.New(l.Dir)
	if err != nil {
		log.Fatal(err)
	}

	l.db = db
	return nil
}

func (l *Log) Set(record *api.Record) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	var key string
	if record.Key == "" {
		uid, _ := uuid.NewUUID()
		key = uid.String()
	} else {
		key = record.Key
	}
	err := l.db.Set(key, record.Value[:])
	if err != nil {
		return "", err
	}

	return key, nil
}

func (l *Log) Get(key string) (*api.Record, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	val := l.db.Get(key)
	if val == nil {
		return nil, nil
	}

	//Type assertion http://golang.org/ref/spec#Type_assertions
	if str, ok := val.(string); ok {
		record := &api.Record{
			Key:   key,
			Value: str,
		}
		return record, nil
	} else {
		return nil, nil
	}

}

func (l *Log) Insert(item badger.Item) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.db.Insert(item)
}

func (l *Log) FindBy(field, expression string) ([]badger.Item, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	items, err := l.db.FindBy(field, expression)
	if err != nil {
		return nil, err
	}
	return items, err
}

type cError struct{}

func (e *cError) Error() string {
	return "Log Error"
}

func (l *Log) Delete(key string) error {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if err := l.db.Delete(key); !err {
		return &cError{}
	}
	return os.RemoveAll(l.Dir)
}

func (l *Log) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if err := l.db.Close(); err != nil {
		return err
	}

	return nil
}

func (l *Log) ExistKeyWith(prefix []byte, neq []byte) (bool, error) {
	return l.idx.ExistKeyWith(prefix, neq)

}

type StoreDB interface {
	// Set key and value
	Set(key string, value string) error
	// Get value by key
	Get(key string) interface{}
	// Delete removes a session key value based on its key.
	Delete(key string) (deleted bool)
	// Close connection to store
	Close() error
}

// This interface hold the functions that help to manage the cluster database
type CoreDB interface {
	Insert(item badger.Item) error
	FindBy(field, expression string) ([]badger.Item, error)
}
