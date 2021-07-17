package boltdb

import (
	"os"
	"testing"

	"github.com/boltdb/bolt"
)

func TestInitDB(t *testing.T) {
	db1, _ := InitDB("./db1.db")
	db2, _ := InitDB("./db1.db")

	db1.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Bucket"))
		if err != nil {
			t.Errorf("create bucket: %s", err)
		}
		err = b.Put([]byte("answer"), []byte("42"))
		if err != nil {
			t.Errorf("putting value: %s", err)
		}

		return nil
	})

	db2.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Bucket"))
		v := b.Get([]byte("answer"))
		if string(v) != "42" {
			t.Errorf("invalid value, expected 42, got %s", string(v))
		}
		return nil
	})

	db2.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte("Bucket"))
		if err != nil {
			t.Errorf("error deleting the bucket : %s", err)
		}
		return nil
	})

	err := CloseDB(db1)
	if err != nil {
		t.Errorf("error closing database with error: %v", err)
	}
	err = CloseDB(db2)
	if err != nil {
		t.Errorf("error closing database with error: %v", err)
	}
	// trying to close closed database
	err = CloseDB(db1)
	if err == nil {
		t.Errorf("no error closing closed database ")
	}
	os.Remove("./db1.db")
}
