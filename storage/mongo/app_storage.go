package mongo

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/madappgang/identifo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const appsCollectionName = "Applications"

// NewAppStorage creates new MongoDB AppStorage implementation.
func NewAppStorage(db *DB) (model.AppStorage, error) {
	coll := db.Database.Collection(appsCollectionName)
	return &AppStorage{coll: coll, timeout: 30 * time.Second}, nil
}

// AppStorage is a fully functional app storage for MongoDB.
type AppStorage struct {
	coll    *mongo.Collection
	timeout time.Duration
}

// NewAppData returns pointer to newly created app data.
func (as *AppStorage) NewAppData() model.AppData {
	return &AppData{appData: appData{}}
}

// AppByID returns app from MongoDB by ID.
func (as *AppStorage) AppByID(id string) (model.AppData, error) {
	hexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), as.timeout)
	defer cancel()

	var ad appData
	if err := as.coll.FindOne(ctx, bson.M{"_id": hexID}).Decode(&ad); err != nil {
		return nil, err
	}
	return &AppData{appData: ad}, nil
}

// ActiveAppByID returns app by id only if it's active.
func (as *AppStorage) ActiveAppByID(appID string) (model.AppData, error) {
	if appID == "" {
		return nil, ErrorEmptyAppID
	}

	app, err := as.AppByID(appID)
	if err != nil {
		return nil, err
	}

	if !app.Active() {
		return nil, ErrorInactiveApp
	}

	return app, nil
}

// FetchApps fetches apps which name satisfies provided filterString.
// Supports pagination.
func (as *AppStorage) FetchApps(filterString string, skip, limit int) ([]model.AppData, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*as.timeout)
	defer cancel()

	q := bson.D{primitive.E{Key: "name", Value: primitive.Regex{Pattern: filterString, Options: "i"}}}

	total, err := as.coll.CountDocuments(ctx, q)
	if err != nil {
		return []model.AppData{}, 0, err
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{primitive.E{Key: "name", Value: 1}})
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))

	curr, err := as.coll.Find(ctx, q, findOptions)
	if err != nil {
		return []model.AppData{}, 0, err
	}

	var appsData []appData
	if err = curr.All(ctx, &appsData); err != nil {
		return []model.AppData{}, 0, err
	}

	apps := make([]model.AppData, len(appsData))
	for i := 0; i < len(appsData); i++ {
		apps[i] = &AppData{appData: appsData[i]}
	}
	return apps, int(total), nil
}

// DeleteApp deletes app by id.
func (as *AppStorage) DeleteApp(id string) error {
	hexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), as.timeout)
	defer cancel()

	if _, err := as.coll.DeleteOne(ctx, bson.M{"_id": hexID}); err != nil {
		return err
	}
	return nil
}

// CreateApp creates new app in MongoDB.
func (as *AppStorage) CreateApp(app model.AppData) (model.AppData, error) {
	res, ok := app.(*AppData)
	if !ok || res == nil {
		return nil, model.ErrorWrongDataFormat
	}
	result, err := as.addNewApp(res)
	return result, err
}

// addNewApp adds new app to MongoDB storage.
func (as *AppStorage) addNewApp(app model.AppData) (model.AppData, error) {
	a, ok := app.(*AppData)
	if !ok || a == nil {
		return nil, model.ErrorWrongDataFormat
	}

	if _, err := primitive.ObjectIDFromHex(app.ID()); err != nil {
		a.appData.ID = primitive.NewObjectID()
	}

	ctx, cancel := context.WithTimeout(context.Background(), as.timeout)
	defer cancel()

	if _, err := as.coll.InsertOne(ctx, a.appData); err != nil {
		return nil, err
	}
	return app, nil
}

// DisableApp disables app in MongoDB storage.
func (as *AppStorage) DisableApp(app model.AppData) error {
	hexID, err := primitive.ObjectIDFromHex(app.ID())
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), as.timeout)
	defer cancel()

	update := bson.M{"$set": bson.M{"active": false}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var ad appData
	if err := as.coll.FindOneAndUpdate(ctx, bson.M{"_id": hexID}, update, opts).Decode(&ad); err != nil {
		return err
	}
	//maybe return updated data?
	return nil
}

// UpdateApp updates app in MongoDB storage.
func (as *AppStorage) UpdateApp(appID string, newApp model.AppData) (model.AppData, error) {
	res, ok := newApp.(*AppData)
	if !ok || res == nil {
		return nil, model.ErrorWrongDataFormat
	}

	hexID, err := primitive.ObjectIDFromHex(appID)
	if err != nil {
		return nil, err
	}

	// use ID from the request if it's not set
	if len(res.ID()) == 0 {
		res.appData.ID = hexID
	}

	ctx, cancel := context.WithTimeout(context.Background(), as.timeout)
	defer cancel()

	update := bson.M{"$set": res.appData}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var ad appData
	if err = as.coll.FindOneAndUpdate(ctx, bson.M{"_id": hexID}, update, opts).Decode(&ad); err != nil {
		return nil, err
	}

	return &AppData{appData: ad}, nil
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
	apd := []appData{}
	if err := json.Unmarshal(data, &apd); err != nil {
		log.Println(err)
		return err
	}
	for _, a := range apd {
		if _, err := as.addNewApp(&AppData{appData: a}); err != nil {
			return err
		}
	}
	return nil
}

// Close is a no-op.
func (as *AppStorage) Close() {}
