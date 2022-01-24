package boltdb

import (
	"os"
	"testing"

	bolt "go.etcd.io/bbolt"
)

func TestInitDB(t *testing.T) {
	db1, err := InitDB("./db1.db")
	if err != nil {
		t.Errorf("error open first connection to DB: %v", err)
	}

	db2, err := InitDB("./db1.db")

	if err != nil {
		t.Errorf("error open second connection to DB: %v", err)
	}

	err = db1.Update(func(tx *bolt.Tx) error {
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
	if err != nil {
		t.Errorf("error update database with error: %v", err)
	}

	err = db2.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Bucket"))
		v := b.Get([]byte("answer"))
		if string(v) != "42" {
			t.Errorf("invalid value, expected 42, got %s", string(v))
		}
		return nil
	})
	if err != nil {
		t.Errorf("error view database with error: %v", err)
	}

	err = db2.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte("Bucket"))
		if err != nil {
			t.Errorf("error deleting the bucket : %s", err)
		}
		return nil
	})
	if err != nil {
		t.Errorf("error delete bucket database with error: %v", err)
	}

	err = CloseDB(db1)
	if err != nil {
		t.Errorf("error closing first connection with error: %v", err)
	}
	err = CloseDB(db2)
	if err != nil {
		t.Errorf("error closing second connection with error: %v", err)
	}
	// trying to close closed database
	err = CloseDB(db1)
	if err == nil {
		t.Errorf("no error closing already closed database ")
	}
	// Recreate new connection with empty poll
	db3, err := InitDB("./db1.db")
	if err != nil {
		t.Errorf("error closing database with error: %v", err)
	}
	err = db3.Update(func(tx *bolt.Tx) error {
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
	if err != nil {
		t.Errorf("error view database with error: %v", err)
	}
	os.Remove("./db1.db")
}
