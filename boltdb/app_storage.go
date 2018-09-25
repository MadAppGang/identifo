package boltdb

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/madappgang/identifo/model"
)

const (
	//AppBucket bucket name with apps
	AppBucket = "Apps"
)

//NewAppStorage creates new embedded AppStorage implementation
func NewAppStorage(db *bolt.DB) (model.AppStorage, error) {
	as := AppStorage{}
	as.db = db
	//ensure we have app's bucket in the database
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(AppBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &as, nil
}

//AppStorage is fully functional app storage
type AppStorage struct {
	db *bolt.DB
}

//AppByID returns app from memory by ID
func (as *AppStorage) AppByID(id string) (model.AppData, error) {
	var res AppData
	if err := as.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(AppBucket))
		v := b.Get([]byte(id))
		if v == nil {
			return ErrorNotFound
		}
		rr, err := AppDataFromJSON(v)
		res = rr
		return err
	}); err != nil {
		return nil, err
	}
	return res, nil
}

//AddNewApp add new app to memory storage
func (as *AppStorage) AddNewApp(app model.AppData) (model.AppData, error) {
	res, ok := app.(AppData)
	if !ok {
		return nil, ErrorWrongDataFormat
	}
	return app, as.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(AppBucket))
		data, err := res.Marshal()
		fmt.Println("ADDING: " + string(data))
		if err != nil {
			return err
		}
		return b.Put([]byte(res.ID()), data)
	})
}

//DisableApp disables app from storage
func (as *AppStorage) DisableApp(app model.AppData) error {
	res, ok := app.(AppData)
	if !ok {
		return ErrorWrongDataFormat
	}
	res.appData.Active = false
	_, err := as.AddNewApp(res)
	return err
}

//UpdateApp updates app in storage
func (as *AppStorage) UpdateApp(oldAppID string, newApp model.AppData) error {
	res, ok := newApp.(AppData)
	if !ok {
		return ErrorWrongDataFormat
	}
	return as.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(AppBucket))
		if err := b.Delete([]byte(oldAppID)); err != nil {
			return err
		}
		data, err := res.Marshal()
		if err != nil {
			return err
		}
		return b.Put([]byte(res.ID()), data)
	})
}

type appData struct {
	ID                   string   `json:"id,omitempty"`
	Secret               string   `json:"secret,omitempty"`
	Active               bool     `json:"active,omitempty"`
	Description          string   `json:"description,omitempty"`
	Scopes               []string `json:"scopes,omitempty"`
	Offline              bool     `json:"offline,omitempty"`
	RedirectURL          string   `json:"redirect_url,omitempty"`
	RefreshTokenLifespan int64    `json:"refresh_token_lifespan,omitempty"`
	TokenLifespan        int64    `json:"token_lifespan,omitempty"`
}

//AppData is memory model for model.AppData
type AppData struct {
	appData
}

//NewAppData instantiate app data memory model from general one
func NewAppData(data model.AppData) AppData {
	return AppData{appData{
		ID:                   data.ID(),
		Secret:               data.Secret(),
		Active:               data.Active(),
		Description:          data.Description(),
		Scopes:               data.Scopes(),
		Offline:              data.Offline(),
		RedirectURL:          data.RedirectURL(),
		RefreshTokenLifespan: data.RefreshTokenLifespan(),
		TokenLifespan:        data.TokenLifespan(),
	}}
}

//AppDataFromJSON deserializes data from JSON
func AppDataFromJSON(d []byte) (AppData, error) {
	add := appData{}
	if err := json.Unmarshal(d, &add); err != nil {
		return AppData{}, err
	}
	return AppData{add}, nil
}

//Marshal serialize data to byte array
func (ad AppData) Marshal() ([]byte, error) {
	return json.Marshal(ad.appData)
}

//MakeAppData creates new memory app data instance
func MakeAppData(id, secret string, active bool, description string, scopes []string, offline bool, redirectURL string, refreshTokenLifespan, tokenLifespan int64) AppData {
	return AppData{appData{id, secret, active, description, scopes, offline, redirectURL, refreshTokenLifespan, tokenLifespan}}

}

func (ad AppData) ID() string                  { return ad.appData.ID }
func (ad AppData) Secret() string              { return ad.appData.Secret }
func (ad AppData) Active() bool                { return ad.appData.Active }
func (ad AppData) Description() string         { return ad.appData.Description }
func (ad AppData) Scopes() []string            { return ad.appData.Scopes }
func (ad AppData) Offline() bool               { return ad.appData.Offline }
func (ad AppData) RedirectURL() string         { return ad.appData.RedirectURL }
func (ad AppData) RefreshTokenLifespan() int64 { return ad.appData.RefreshTokenLifespan }
func (ad AppData) TokenLifespan() int64        { return ad.appData.TokenLifespan }
