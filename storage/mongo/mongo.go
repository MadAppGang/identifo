package mongo

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/madappgang/identifo/v2/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewDB creates new database connection.
func NewDB(logger *slog.Logger, conn string, dbName string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	db := &DB{
		logger:   logger,
		client:   client,
		database: client.Database(dbName),
	}
	return db, nil
}

// DB is database connection structure.
type DB struct {
	logger   *slog.Logger
	database *mongo.Database
	client   *mongo.Client
}

// Close closes database connection.
func (db *DB) Close() {
	if err := db.client.Disconnect(context.TODO()); err != nil {
		db.logger.Error("Error closing mongo storage", logging.FieldError, err)
	}
}

// EnsureCollectionIndices creates indices on a collection.
func (db *DB) EnsureCollectionIndices(collectionName string, newIndices []mongo.IndexModel) error {
	coll := db.database.Collection(collectionName)

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

	keys, ok := index.Keys.(bson.D)
	if !ok {
		return "", fmt.Errorf("incorrect index keys type - expecting bson.D")
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

		value := fmt.Sprintf("%v", elem.Value)

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
