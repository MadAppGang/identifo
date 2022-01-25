package mongo

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// NewDB creates new database connection.
func NewDB(conn string, dbName string) (*DB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(conn))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		return nil, err
	}
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	db := &DB{
		Client:   client,
		Database: client.Database(dbName),
	}
	return db, nil
}

// DB is database connection structure.
type DB struct {
	Database *mongo.Database
	Client   *mongo.Client
}

// Close closes database connection.
func (db *DB) Close() {
	if err := db.Client.Disconnect(context.TODO()); err != nil {
		log.Printf("Error closing mongo storage: %s\n", err)
	}
}

// EnsureCollectionIndices creates indices on a collection.
func (db *DB) EnsureCollectionIndices(collectionName string, newIndices []mongo.IndexModel) error {
	coll := db.Database.Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	for _, newIndex := range newIndices {
		if _, err := coll.Indexes().CreateOne(ctx, newIndex); err != nil && strings.Contains(err.Error(), "already exists with different options") {
			name, err := generateIndexName(newIndex)
			if err != nil {
				return err
			}
			if _, err = coll.Indexes().DropOne(ctx, name); err != nil {
				return err
			}
			if _, err = coll.Indexes().CreateOne(ctx, newIndex); err != nil {
				return err
			}
		}
	}
	return nil
}

func generateIndexName(index mongo.IndexModel) (string, error) {
	if index.Options != nil && index.Options.Name != nil {
		return *index.Options.Name, nil
	}

	name := bytes.NewBufferString("")
	first := true

	keys, ok := index.Keys.(bsonx.Doc)
	if !ok {
		return "", fmt.Errorf("Incorrect index keys type - expecting bsonx.Doc")
	}
	for _, elem := range keys {
		if !first {
			if _, err := name.WriteRune('_'); err != nil {
				return "", err
			}
		}
		if _, err := name.WriteString(elem.Key); err != nil {
			return "", err
		}
		if _, err := name.WriteRune('_'); err != nil {
			return "", err
		}

		var value string

		switch elem.Value.Type() {
		case bsontype.Int32:
			value = fmt.Sprintf("%d", elem.Value.Int32())
		case bsontype.Int64:
			value = fmt.Sprintf("%d", elem.Value.Int64())
		case bsontype.String:
			value = elem.Value.StringValue()
		default:
			return "", fmt.Errorf("Incorrect index value type %s", elem.Value.Type())
		}

		if _, err := name.WriteString(value); err != nil {
			return "", err
		}

		first = false
	}
	return name.String(), nil
}

func isErrNotFound(err error) bool {
	return strings.Contains(err.Error(), "no documents in result")
}

func isErrDuplication(err error) bool {
	return strings.Contains(err.Error(), "duplicate key")
}
