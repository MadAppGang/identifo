package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/madappgang/identifo/v2/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const managementKeysCollectionName = "ManagementKeys"

// ManagementKeysStorage is a MongoDB management keys storage.
type ManagementKeysStorage struct {
	coll *mongo.Collection
}

// NewManagementKeysStorage creates a management keys invite storage.
func NewManagementKeysStorage(settings model.MongoDatabaseSettings) (model.ManagementKeysStorage, error) {
	if len(settings.ConnectionString) == 0 || len(settings.DatabaseName) == 0 {
		return nil, ErrorEmptyConnectionStringDatabase
	}

	// create database
	db, err := NewDB(settings.ConnectionString, settings.DatabaseName)
	if err != nil {
		return nil, err
	}

	coll := db.Database.Collection(managementKeysCollectionName)
	return &ManagementKeysStorage{coll: coll}, nil
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
	key.Secret = "***"
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

	for _, v := range keys {
		v.Secret = "***"
	}

	return keys, nil
}
