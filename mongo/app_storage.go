package mongo

import (
	"encoding/json"

	"github.com/madappgang/identifo/model"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	//AppsCollection collection name to store apps data
	AppsCollection = "Applications"
)

//NewAppStorage creates new mongo AppStorage implementation
func NewAppStorage(db *DB) (model.AppStorage, error) {
	return &AppStorage{db: db}, nil
}

//AppStorage is fully functional app storage in mongo
type AppStorage struct {
	db *DB
}

//AppByID returns app from mongo by ID
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

//AddNewApp add new app to mongo storage
func (as *AppStorage) AddNewApp(app model.AppData) (model.AppData, error) {
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

//DisableApp disables app in mongo storage
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

//UpdateApp updates app in mongo storage
func (as *AppStorage) UpdateApp(oldAppID string, newApp model.AppData) error {
	if !bson.IsObjectIdHex(oldAppID) {
		return model.ErrorWrongDataFormat
	}
	s := as.db.Session(AppsCollection)
	defer s.Close()

	var ad appData
	update := mgo.Change{
		Update:    bson.M{"$set": newApp},
		ReturnNew: true,
	}
	if _, err := s.C.FindId(bson.ObjectId(oldAppID)).Apply(update, &ad); err != nil {
		return err
	}
	//maybe return updated data?
	return nil
}

type appData struct {
	ID                         bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Secret                     string        `bson:"secret,omitempty" json:"secret,omitempty"`
	Active                     bool          `bson:"active,omitempty" json:"active,omitempty"`
	Description                string        `bson:"description,omitempty" json:"description,omitempty"`
	Scopes                     []string      `bson:"scopes,omitempty" json:"scopes,omitempty"`
	Offline                    bool          `bson:"offline,omitempty" json:"offline,omitempty"`
	RedirectURL                string        `bson:"redirect_url,omitempty" json:"redirect_url,omitempty"`
	RefreshTokenLifespan       int64         `bson:"refresh_token_lifespan,omitempty" json:"refresh_token_lifespan,omitempty"`
	ResetPasswordTokenLifespan int64         `bson:"reset_password_token_lifespan,omitempty" json:"reset_password_token_lifespan,omitempty"`
	TokenLifespan              int64         `bson:"token_lifespan,omitempty" json:"token_lifespan,omitempty"`
}

//AppData is mongo model for model.AppData
type AppData struct {
	appData
}

//NewAppData instantiate app data mongo model from general one
func NewAppData(data model.AppData) (AppData, error) {
	if !bson.IsObjectIdHex(data.ID()) {
		return AppData{}, model.ErrorWrongDataFormat
	}
	return AppData{appData: appData{
		ID:                         bson.ObjectIdHex(data.ID()),
		Secret:                     data.Secret(),
		Active:                     data.Active(),
		Description:                data.Description(),
		Scopes:                     data.Scopes(),
		Offline:                    data.Offline(),
		RedirectURL:                data.RedirectURL(),
		RefreshTokenLifespan:       data.RefreshTokenLifespan(),
		ResetPasswordTokenLifespan: data.ResetPasswordTokenLifespan(),
		TokenLifespan:              data.TokenLifespan(),
	}}, nil
}

//AppDataFromJSON deserializes data from JSON
func AppDataFromJSON(d []byte) (AppData, error) {
	apd := appData{}
	if err := json.Unmarshal(d, &apd); err != nil {
		return AppData{}, err
	}
	return AppData{appData: apd}, nil
}

//Marshal serialize data to byte array
func (ad AppData) Marshal() ([]byte, error) {
	return json.Marshal(ad.appData)
}

//MakeAppData creates new mongo app data instance
func MakeAppData(id, secret string, active bool, description string, scopes []string, offline bool, redirectURL string, refreshTokenLifespan, resetPasswordTokenLifespan, tokenLifespan int64) (AppData, error) {
	if !bson.IsObjectIdHex(id) {
		return AppData{}, model.ErrorWrongDataFormat
	}
	return AppData{appData: appData{
		ID:                         bson.ObjectIdHex(id),
		Secret:                     secret,
		Active:                     active,
		Description:                description,
		Scopes:                     scopes,
		Offline:                    offline,
		RedirectURL:                redirectURL,
		RefreshTokenLifespan:       refreshTokenLifespan,
		ResetPasswordTokenLifespan: resetPasswordTokenLifespan,
		TokenLifespan:              tokenLifespan,
	}}, nil
}

func (ad AppData) ID() string                        { return ad.appData.ID.Hex() }
func (ad AppData) Secret() string                    { return ad.appData.Secret }
func (ad AppData) Active() bool                      { return ad.appData.Active }
func (ad AppData) Description() string               { return ad.appData.Description }
func (ad AppData) Scopes() []string                  { return ad.appData.Scopes }
func (ad AppData) Offline() bool                     { return ad.appData.Offline }
func (ad AppData) RedirectURL() string               { return ad.appData.RedirectURL }
func (ad AppData) RefreshTokenLifespan() int64       { return ad.appData.RefreshTokenLifespan }
func (ad AppData) ResetPasswordTokenLifespan() int64 { return ad.appData.ResetPasswordTokenLifespan }
func (ad AppData) TokenLifespan() int64              { return ad.appData.TokenLifespan }
