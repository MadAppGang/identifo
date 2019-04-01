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
	return AppData{appData: ad}, nil
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
func (as *AppStorage) FetchApps(filterString string, skip, limit int) ([]model.AppData, error) {
	s := as.db.Session(AppsCollection)
	defer s.Close()

	q := bson.M{"name": bson.M{"$regex": bson.RegEx{Pattern: filterString, Options: "i"}}}

	orderByField := "name"

	var appsData []appData
	if err := s.C.Find(q).Sort(orderByField).Limit(limit).Skip(skip).All(&appsData); err != nil {
		return nil, err
	}

	apps := make([]model.AppData, len(appsData))
	for i := 0; i < len(appsData); i++ {
		apps[i] = AppData{appData: appsData[i]}
	}

	return apps, nil
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
	if !ok || app == nil {
		return nil, model.ErrorWrongDataFormat
	}
	result, err := as.addNewApp(*res)
	return result, err
}

// addNewApp adds new app to MongoDB storage.
func (as *AppStorage) addNewApp(app model.AppData) (model.AppData, error) {
	a, ok := app.(AppData)
	if !ok {
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

	return AppData{appData: ad}, nil
}

// ImportJSON imports data from JSON.
func (as *AppStorage) ImportJSON(data []byte) error {
	apd := []appData{}
	if err := json.Unmarshal(data, &apd); err != nil {
		log.Println(err)
		return err
	}
	for _, a := range apd {
		if _, err := as.addNewApp(AppData{appData: a}); err != nil {
			return err
		}
	}
	return nil
}

type appData struct {
	ID                    bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Secret                string        `bson:"secret,omitempty" json:"secret,omitempty"`
	Active                bool          `bson:"active,omitempty" json:"active,omitempty"`
	Name                  string        `bson:"name,omitempty" json:"name,omitempty"`
	Description           string        `bson:"description,omitempty" json:"description,omitempty"`
	Scopes                []string      `bson:"scopes,omitempty" json:"scopes,omitempty"`
	Offline               bool          `bson:"offline,omitempty" json:"offline,omitempty"`
	RedirectURL           string        `bson:"redirect_url,omitempty" json:"redirect_url,omitempty"`
	RefreshTokenLifespan  int64         `bson:"refresh_token_lifespan,omitempty" json:"refresh_token_lifespan,omitempty"`
	TokenLifespan         int64         `bson:"token_lifespan,omitempty" json:"token_lifespan,omitempty"`
	TokenPayload          []string      `bson:"token_payload,omitempty" json:"token_payload,omitempty"`
	RegistrationForbidden bool          `bson:"registration_forbidden,omitempty" json:"registration_forbidden,omitempty"`
}

// AppData is a MongoDb model that implements model.AppData.
type AppData struct {
	appData
}

// NewAppData instantiates MongoDB app data model from the general one.
func NewAppData(data model.AppData) (AppData, error) {
	if !bson.IsObjectIdHex(data.ID()) {
		return AppData{}, model.ErrorWrongDataFormat
	}
	return AppData{appData: appData{
		ID:                    bson.ObjectIdHex(data.ID()),
		Secret:                data.Secret(),
		Active:                data.Active(),
		Name:                  data.Name(),
		Description:           data.Description(),
		Scopes:                data.Scopes(),
		Offline:               data.Offline(),
		RedirectURL:           data.RedirectURL(),
		RefreshTokenLifespan:  data.RefreshTokenLifespan(),
		TokenLifespan:         data.TokenLifespan(),
		TokenPayload:          data.TokenPayload(),
		RegistrationForbidden: data.RegistrationForbidden(),
	}}, nil
}

// AppDataFromJSON deserializes app data from JSON.
func AppDataFromJSON(d []byte) (AppData, error) {
	apd := appData{}
	if err := json.Unmarshal(d, &apd); err != nil {
		return AppData{}, err
	}
	return AppData{appData: apd}, nil
}

// Marshal serializes data to byte array.
func (ad AppData) Marshal() ([]byte, error) {
	return json.Marshal(ad.appData)
}

// MakeAppData creates new MongoDB app data instance.
func MakeAppData(id, secret string, active bool, name, description string, scopes []string, offline bool, redirectURL string, refreshTokenLifespan, tokenLifespan int64, tokenPayload []string, registrationForbidden bool) (AppData, error) {
	if !bson.IsObjectIdHex(id) {
		return AppData{}, model.ErrorWrongDataFormat
	}
	return AppData{appData: appData{
		ID:                    bson.ObjectIdHex(id),
		Secret:                secret,
		Active:                active,
		Name:                  name,
		Description:           description,
		Scopes:                scopes,
		Offline:               offline,
		RedirectURL:           redirectURL,
		RefreshTokenLifespan:  refreshTokenLifespan,
		TokenLifespan:         tokenLifespan,
		TokenPayload:          tokenPayload,
		RegistrationForbidden: registrationForbidden,
	}}, nil
}

// Sanitize removes all sensitive data.
func (ad AppData) Sanitize() model.AppData {
	ad.appData.Secret = ""
	return ad
}

// ID implements model.AppData interface.
func (ad AppData) ID() string { return ad.appData.ID.Hex() }

// Secret implements model.AppData interface.
func (ad AppData) Secret() string { return ad.appData.Secret }

// Active implements model.AppData interface.
func (ad AppData) Active() bool { return ad.appData.Active }

// Name implements model.AppData interface.
func (ad AppData) Name() string { return ad.appData.Name }

// Description implements model.AppData interface.
func (ad AppData) Description() string { return ad.appData.Description }

// Scopes implements model.AppData interface.
func (ad AppData) Scopes() []string { return ad.appData.Scopes }

// Offline implements model.AppData interface.
func (ad AppData) Offline() bool { return ad.appData.Offline }

// RedirectURL implements model.AppData interface.
func (ad AppData) RedirectURL() string { return ad.appData.RedirectURL }

// RefreshTokenLifespan implements model.AppData interface.
func (ad AppData) RefreshTokenLifespan() int64 { return ad.appData.RefreshTokenLifespan }

// TokenLifespan implements model.AppData interface.
func (ad AppData) TokenLifespan() int64 { return ad.appData.TokenLifespan }

// TokenPayload implements model.AppData interface.
func (ad AppData) TokenPayload() []string { return ad.appData.TokenPayload }

// RegistrationForbidden implements model.AppData interface.
func (ad AppData) RegistrationForbidden() bool { return ad.appData.RegistrationForbidden }
