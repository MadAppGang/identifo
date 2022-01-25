package boltdb

import (
	"errors"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
)

// database open timeout is 3 secs
const dbOpenTimeout = 3 * time.Second

// pool is a pool of shared database connections to boltdb
// BoltDB does not support concurent access to database file
// so we are reusing the same database connection if different storages
// ask to access to same file
// var pool = make(map[string]*bolt.DB)
var pool = sync.Map{}

// poolCounter works as automatic reference counter
// if we request close for specific database, we decrease the counter unless it goes to zero
// and then actually closing it
var poolCounter = sync.Map{}

var ErrorClosingNonExistentDatabase = errors.New("error closing the database which is already closed or not in the pool")

// InitDB opens database.
func InitDB(file string) (*bolt.DB, error) {
	result, ok := pool.Load(file)
	if ok {
		db := result.(*bolt.DB)
		increaseCounter(file, db)
		return db, nil
	} else {
		db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: dbOpenTimeout})
		if err != nil {
			return nil, err
		}
		pool.Store(file, db)
		// create first reference counter to the object
		poolCounter.Store(db, 1)
		return db, nil
	}
}

func increaseCounter(file string, db *bolt.DB) {
	result, _ := poolCounter.LoadOrStore(db, 0)
	newValue := result.(int) + 1
	poolCounter.Store(db, newValue)
}

// CloseDB closes database.
func CloseDB(db *bolt.DB) error {
	result, ok := poolCounter.Load(db)
	if ok {
		counter := result.(int)
		// it is the last counter, remove it from counter and close it
		if counter <= 1 {
			poolCounter.Delete(db)
			pool.Delete(db.Path())
			return db.Close()
		}
		// if we have more than 1 reference, we just decrease the reference count
		poolCounter.Store(db, counter-1)
		return nil
	} else {
		// we have not database to close
		// let's close it directly
		return ErrorClosingNonExistentDatabase
	}
}
