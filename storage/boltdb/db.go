package boltdb

import "github.com/boltdb/bolt"

// InitDB opens database.
func InitDB(file string) (*bolt.DB, error) {
	return bolt.Open(file, 0600, nil)
}

// CloseDB closes database.
func CloseDB(db *bolt.DB) error {
	return db.Close()
}
