package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const managementKeysCollectionName = "ManagementKeys"

// ManagementKeysStorage is a MongoDB management keys storage.
type ManagementKeysStorage struct {
	logger *slog.Logger
	coll   *mongo.Collection
}

// NewManagementKeysStorage creates a management keys invite storage.
func NewManagementKeysStorage(
	logger *slog.Logger,
	settings model.MongoDatabaseSettings,
) (model.ManagementKeysStorage, error) {
	if len(settings.ConnectionString) == 0 || len(settings.DatabaseName) == 0 {
		return nil, ErrorEmptyConnectionStringDatabase
	}

	// create database
	db, err := NewDB(logger, settings.ConnectionString, settings.DatabaseName)
	if err != nil {
		return nil, err
	}

	coll := db.database.Collection(managementKeysCollectionName)
	return &ManagementKeysStorage{
		logger: logger,
		coll:   coll,
	}, nil
}

func (ms *ManagementKeysStorage) GetKey(ctx context.Context, id string) (model.ManagementKey, error) {
	filter := bson.M{"_id": id}

	var key model.ManagementKey
	if err := ms.coll.FindOne(ctx, filter).Decode(&key); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.ManagementKey{}, model.ErrorNotFound
		}
		return model.ManagementKey{}, err
	}
	return key, nil
}

func (ms *ManagementKeysStorage) CreateKey(ctx context.Context, name string, scopes []string) (model.ManagementKey, error) {
	key := model.ManagementKey{
		Name:      name,
		Scopes:    scopes,
		ID:        primitive.NewObjectID().Hex(),
		Active:    true,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
	}

	_, err := ms.coll.InsertOne(ctx, key)
	return key, err
}

func (ms *ManagementKeysStorage) AddKey(ctx context.Context, key model.ManagementKey) (model.ManagementKey, error) {
	_, err := ms.coll.InsertOne(ctx, key)
	return key, err
}

func (ms *ManagementKeysStorage) DisableKey(ctx context.Context, id string) (model.ManagementKey, error) {
	update := bson.M{"$set": bson.M{"active": false}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var key model.ManagementKey
	if err := ms.coll.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&key); err != nil {
		return key, err
	}

	return key, nil
}

func (ms *ManagementKeysStorage) RenameKey(ctx context.Context, id, name string) (model.ManagementKey, error) {
	update := bson.M{"$set": bson.M{"name": name}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var key model.ManagementKey
	if err := ms.coll.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&key); err != nil {
		return key, err
	}

	return key, nil
}

func (ms *ManagementKeysStorage) ChangeScopesForKey(ctx context.Context, id string, scopes []string) (model.ManagementKey, error) {
	update := bson.M{"$set": bson.M{"scopes": scopes}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var key model.ManagementKey
	if err := ms.coll.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&key); err != nil {
		return key, err
	}

	return key, nil
}

func (ms *ManagementKeysStorage) UseKey(ctx context.Context, id string) (model.ManagementKey, error) {
	update := bson.M{"$set": bson.M{"lastUsed": time.Now()}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var key model.ManagementKey
	if err := ms.coll.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&key); err != nil {
		return key, err
	}

	return key, nil
}

func (ms *ManagementKeysStorage) GeyAllKeys(ctx context.Context) ([]model.ManagementKey, error) {
	curr, err := ms.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var keys []model.ManagementKey
	if err = curr.All(ctx, &keys); err != nil {
		return nil, err
	}

	return keys, nil
}

func (ms *ManagementKeysStorage) ClearAllData() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := ms.coll.DeleteMany(ctx, bson.M{}); err != nil {
		ms.logger.Error("Error cleaning all user data",
			logging.FieldError, err)
	}
}

// ImportJSON imports data from JSON.
func (ms *ManagementKeysStorage) ImportJSON(data []byte, cleanOldData bool) error {
	if cleanOldData {
		ms.ClearAllData()
	}

	keys := []model.ManagementKey{}
	if err := json.Unmarshal(data, &keys); err != nil {
		return fmt.Errorf("failed to unmarshal management keys: %w", err)
	}

	for _, a := range keys {
		if _, err := ms.AddKey(context.TODO(), a); err != nil {
			return fmt.Errorf("failed to add management keys: %w", err)
		}
	}
	return nil
}
