package mongo

import (
	"encoding/json"
	"log"

	"github.com/madappgang/identifo/model"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	// AppsCollection is a collection name for storing apps data.
	AppsCollection = "Applications"
)

// NewAppStorage creates new MongoDB AppStorage implementation.
func NewAppStorage(db *DB) (model.AppStorage, error) {
	return &AppStorage{db: db}, nil
}

// AppStorage is a fully functional app storage for MongoDB.
type AppStorage struct {
	db *DB
}

// NewAppData returns pointer to newly created app data.
func (as *AppStorage) NewAppData() model.AppData {
	return &AppData{appData: appData{}}
}

// AppByID returns app from MongoDB by ID.
func (as *AppStorage) AppByID(id string) (model.AppData, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, model.ErrorWrongDataFormat
	}
	s := as.db.Session(AppsCollection)
	defer s.Close()

	var ad appData
	if err := s.C.FindId(bson.ObjectIdHex(id)).One(&ad); err != nil {
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
	s := as.db.Session(AppsCollection)
	defer s.Close()

	q := bson.M{"name": bson.M{"$regex": bson.RegEx{Pattern: filterString, Options: "i"}}}

	total, err := s.C.Find(q).Count()
	if err != nil {
		return []model.AppData{}, 0, err
	}

	orderByField := "name"

	var appsData []appData
	if err := s.C.Find(q).Sort(orderByField).Limit(limit).Skip(skip).All(&appsData); err != nil {
		return []model.AppData{}, 0, err
	}

	apps := make([]model.AppData, len(appsData))
	for i := 0; i < len(appsData); i++ {
		apps[i] = &AppData{appData: appsData[i]}
	}

	return apps, total, nil
}

// DeleteApp deletes app by id.
func (as *AppStorage) DeleteApp(id string) error {
	if !bson.IsObjectIdHex(id) {
		return model.ErrorWrongDataFormat
	}
	s := as.db.Session(AppsCollection)
	defer s.Close()

	err := s.C.RemoveId(bson.ObjectIdHex(id))
	return err
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
	s := as.db.Session(AppsCollection)
	defer s.Close()

	if !a.appData.ID.Valid() {
		a.appData.ID = bson.NewObjectId()
	}
	if err := s.C.Insert(a.appData); err != nil {
		return nil, err
	}
	return app, nil
}

// DisableApp disables app in MongoDB storage.
func (as *AppStorage) DisableApp(app model.AppData) error {
	if !bson.IsObjectIdHex(app.ID()) {
		return model.ErrorWrongDataFormat
	}
	s := as.db.Session(AppsCollection)
	defer s.Close()

	var ad appData
	update := mgo.Change{
		Update:    bson.M{"$set": bson.M{"active": false}},
		ReturnNew: true,
	}
	if _, err := s.C.FindId(bson.ObjectId(app.ID())).Apply(update, &ad); err != nil {
		return err
	}
	//maybe return updated data?
	return nil
}

// UpdateApp updates app in MongoDB storage.
func (as *AppStorage) UpdateApp(appID string, newApp model.AppData) (model.AppData, error) {
	if !bson.IsObjectIdHex(appID) {
		return nil, model.ErrorWrongDataFormat
	}

	res, ok := newApp.(*AppData)
	if !ok || res == nil {
		return nil, model.ErrorWrongDataFormat
	}

	// use ID from the request if it's not set
	if len(res.ID()) == 0 {
		res.appData.ID = bson.ObjectId(appID)
	}

	s := as.db.Session(AppsCollection)
	defer s.Close()

	var ad appData
	update := mgo.Change{
		Update:    bson.M{"$set": *res},
		ReturnNew: true,
	}
	if _, err := s.C.FindId(bson.ObjectId(appID)).Apply(update, &ad); err != nil {
		return nil, err
	}

	return &AppData{appData: ad}, nil
}

// TestDatabaseConnection checks whether we can fetch the first document in the applications collection.
func (as *AppStorage) TestDatabaseConnection() error {
	s := as.db.Session(AppsCollection)
	defer s.Close()

	var ad appData
	err := s.C.Find(nil).One(&ad)
	if err == mgo.ErrNotFound { // It's OK, collection is empty.
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

// Close closes database connection.
func (as *AppStorage) Close() {
	as.db.Close()
}
