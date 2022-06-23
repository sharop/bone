package store

import (
	"errors"
	"fmt"
	"github.com/timshannon/badgerhold/v4"
	inform "log"
	"math"
	"os"
	"runtime"
	"sync/atomic"
)

// DefaultFileMode used as the default database's "fileMode"
// for creating the sessions directory path, opening and write the session file.
var (
	DefaultFileMode = 0755
)

var ErrDBClosed = errors.New("already closed")

// Database the badger(key-value file-based) session storage.
type Database struct {
	// Service is the underline badger database connection,
	// it's initialized at `New` or `NewFromDB`.
	// Can be used to get stats.
	Service *badgerhold.Store

	closed uint32 // if 1 is closed.

	//Options raft.Options
}

// New creates and returns a new badger(key-value file-based) storage
// instance based on the "directoryPath".
// DirectoryPath should is the directory which the badger database will store the sessions,
// i.e ./sessions
//
// It will remove any old session files.
func New(directoryPath string) (*Database, error) {

	if directoryPath == "" {
		return nil, errors.New("directoryPath is missing")
	}

	lindex := directoryPath[len(directoryPath)-1]
	if lindex != os.PathSeparator && lindex != '/' {
		directoryPath += string(os.PathSeparator)
	}
	// create directories if necessary
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		if err := os.MkdirAll(directoryPath, os.FileMode(DefaultFileMode)); err != nil {
			return nil, errors.New(fmt.Sprintf("DB_LEVEL: Could not create a data directory. %+v", err))
		}
	}

	options := badgerhold.DefaultOptions
	options.SyncWrites = false
	options.NumVersionsToKeep = math.MaxInt32

	options.Dir = directoryPath
	options.ValueDir = directoryPath

	service, err := badgerhold.Open(options)

	if err != nil {
		inform.Printf("unable to initialize the badger-based session database: %v\n", err)
		return nil, err
	}

	return NewFromDB(service), nil
}

// NewFromDB same as `New` but accepts an already-created custom badger connection instead.
func NewFromDB(service *badgerhold.Store) *Database {
	db := &Database{Service: service}
	runtime.SetFinalizer(db, closeDB)
	return db
}

var delim = byte('_')
var sid = "sio"

func makePrefix(sid string) []byte {
	return append([]byte(sid), delim)
}

func makeKey(key string) []byte {
	return append(makePrefix(sid), []byte(key)...)
}

// Close shutdowns the badger connection.
func (db *Database) Close() error {
	return closeDB(db)
}

func closeDB(db *Database) error {
	if atomic.LoadUint32(&db.closed) > 0 {
		return nil
	}
	err := db.Service.Close()
	if err == nil {
		atomic.StoreUint32(&db.closed, 1)
	}
	return err

}
