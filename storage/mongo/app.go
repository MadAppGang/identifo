package mongo

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/madappgang/identifo/v2/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const appsCollectionName = "Applications"

// NewAppStorage creates new MongoDB AppStorage implementation.
func NewAppStorage(settings model.MongodDatabaseSettings) (model.AppStorage, error) {
	if len(settings.ConnectionString) == 0 || len(settings.DatabaseName) == 0 {
		return nil, ErrorEmptyConnectionStringDatabase
	}

	// create database
	db, err := NewDB(settings.ConnectionString, settings.DatabaseName)
	if err != nil {
		return nil, err
	}

	coll := db.Database.Collection(appsCollectionName)
	return &AppStorage{coll: coll, timeout: 30 * time.Second}, nil
}

// AppStorage is a fully functional app storage for MongoDB.
type AppStorage struct {
	coll    *mongo.Collection
	timeout time.Duration
}

// AppByID returns app from MongoDB by ID.
func (as *AppStorage) AppByID(id string) (model.AppData, error) {
	hexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.AppData{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), as.timeout)
	defer cancel()

	var ad model.AppData
	if err := as.coll.FindOne(ctx, bson.M{"_id": hexID.Hex()}).Decode(&ad); err != nil {
		return model.AppData{}, err
	}
	return ad, nil
}

// ActiveAppByID returns app by id only if it's active.
func (as *AppStorage) ActiveAppByID(appID string) (model.AppData, error) {
	if appID == "" {
		return model.AppData{}, ErrorEmptyAppID
	}

	app, err := as.AppByID(appID)
	if err != nil {
		return model.AppData{}, err
	}

	if !app.Active {
		return model.AppData{}, ErrorInactiveApp
	}

	return app, nil
}

// FetchApps fetches apps which name satisfies provided filterString.
// Supports pagination.
func (as *AppStorage) FetchApps(filterString string) ([]model.AppData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*as.timeout)
	defer cancel()

	q := bson.D{primitive.E{Key: "name", Value: primitive.Regex{Pattern: filterString, Options: "i"}}}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{primitive.E{Key: "name", Value: 1}})

	curr, err := as.coll.Find(ctx, q, findOptions)
	if err != nil {
		return []model.AppData{}, err
	}

	var appsData []model.AppData
	if err = curr.All(ctx, &appsData); err != nil {
		return []model.AppData{}, err
	}

	apps := make([]model.AppData, len(appsData))
	for i := 0; i < len(appsData); i++ {
		apps[i] = appsData[i]
	}
	return apps, nil
}

// DeleteApp deletes app by id.
func (as *AppStorage) DeleteApp(id string) error {
	hexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), as.timeout)
	defer cancel()

	if _, err := as.coll.DeleteOne(ctx, bson.M{"_id": hexID.Hex()}); err != nil {
		return err
	}
	return nil
}

// CreateApp creates new app in MongoDB.
func (as *AppStorage) CreateApp(app model.AppData) (model.AppData, error) {
	if objID, err := primitive.ObjectIDFromHex(app.ID); err != nil || objID == primitive.NilObjectID {
		app.ID = primitive.NewObjectID().Hex()
	}

	ctx, cancel := context.WithTimeout(context.Background(), as.timeout)
	defer cancel()

	if _, err := as.coll.InsertOne(ctx, app); err != nil {
		return model.AppData{}, err
	}
	return app, nil
}

// DisableApp disables app in MongoDB storage.
func (as *AppStorage) DisableApp(app model.AppData) error {
	hexID, err := primitive.ObjectIDFromHex(app.ID)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), as.timeout)
	defer cancel()

	update := bson.M{"$set": bson.M{"active": false}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var ad model.AppData
	if err := as.coll.FindOneAndUpdate(ctx, bson.M{"_id": hexID.Hex()}, update, opts).Decode(&ad); err != nil {
		return err
	}
	// maybe return updated data?
	return nil
}

// UpdateApp updates app in MongoDB storage.
func (as *AppStorage) UpdateApp(appID string, newApp model.AppData) (model.AppData, error) {
	hexID, err := primitive.ObjectIDFromHex(appID)
	if err != nil {
		return model.AppData{}, err
	}

	// use ID from the request if it's not set
	if len(newApp.ID) == 0 {
		newApp.ID = appID
	}

	ctx, cancel := context.WithTimeout(context.Background(), as.timeout)
	defer cancel()

	update := bson.M{"$set": newApp}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var ad model.AppData
	if err = as.coll.FindOneAndUpdate(ctx, bson.M{"_id": hexID.Hex()}, update, opts).Decode(&ad); err != nil {
		return model.AppData{}, err
	}

	return ad, nil
}

// TestDatabaseConnection checks if we can access applications collection.
func (as *AppStorage) TestDatabaseConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), as.timeout)
	defer cancel()

	_, err := as.coll.Find(ctx, bson.D{})
	if isErrNotFound(err) {
		return nil
	}
	return err
}

// ImportJSON imports data from JSON.
func (as *AppStorage) ImportJSON(data []byte) error {
	apd := []model.AppData{}
	if err := json.Unmarshal(data, &apd); err != nil {
		log.Println(err)
		return err
	}
	for _, a := range apd {
		if _, err := as.CreateApp(a); err != nil {
			return err
		}
	}
	return nil
}

// Close is a no-op.
func (as *AppStorage) Close() {}
