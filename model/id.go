package model

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

// ID is simplified version of bson's ObjectID implementation and it consists of:
// A 4-byte timestamp, representing the ObjectId's creation, measured in seconds since the Unix epoch.
// A 5-byte random value generated once per process. This random value is unique to the machine and process.
// A 3-byte incrementing counter, initialized to a random value.
// The main difference to bson's implementation is that it not in binary format.
// https://pkg.go.dev/github.com/mongodb/mongo-go-driver/bson/objectid
type ID string

// NewUserID is a special ID for cases when we need to indicate the data belongs to a user who is not registered yet.
// we are using valid ObjectID compatible value here.
var (
	NewUserID       = ID("64acc61bab36eb8395ce5846")
	NewUserUsername = "ephemeral"
)
var (
	processUnique   = processUniqueBytes()
	objectIDCounter = readRandomUint32()
)

func NewID() ID {
	var b [12]byte
	timestamp := time.Now()
	binary.BigEndian.PutUint32(b[0:4], uint32(timestamp.Unix()))
	copy(b[4:9], processUnique[:])
	putUint24(b[9:12], atomic.AddUint32(&objectIDCounter, 1))
	var buf [24]byte
	hex.Encode(buf[:], b[:])
	return ID(buf[:])
}

func (id ID) String() string {
	return string(id)
}

func (id ID) IsNewUserID() bool {
	return id == NewUserID
}

func processUniqueBytes() [5]byte {
	var b [5]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(fmt.Errorf("cannot initialize objectid package with crypto.rand.Reader: %v", err))
	}

	return b
}

func putUint24(b []byte, v uint32) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}

func readRandomUint32() uint32 {
	var b [4]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(fmt.Errorf("cannot initialize objectid package with crypto.rand.Reader: %v", err))
	}

	return (uint32(b[0]) << 0) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24)
}
