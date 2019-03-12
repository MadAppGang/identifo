package boltdb

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/madappgang/identifo/model"
	"github.com/rs/xid"
)

const (
	// AppBucket is a name for bucket with apps.
	AppBucket = "Apps"
)

// NewAppStorage creates new embedded AppStorage implementation.
func NewAppStorage(db *bolt.DB) (model.AppStorage, error) {
	as := AppStorage{db: db}
	// ensure we have app's bucket in the database
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(AppBucket)); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &as, nil
}

// AppStorage is a fully functional app storage.
type AppStorage struct {
	db *bolt.DB
}

// NewAppData returns newly created app data.
func (as *AppStorage) NewAppData() model.AppData {
	return &AppData{appData: appData{}}
}

// AppByID returns app from memory by ID.
func (as *AppStorage) AppByID(id string) (model.AppData, error) {
	var res AppData
	if err := as.db.View(func(tx *bolt.Tx) error {
		ab := tx.Bucket([]byte(AppBucket))
		app := ab.Get([]byte(id))
		if app == nil {
			return model.ErrorNotFound
		}

		var err error
		res, err = AppDataFromJSON(app)
		return err
	}); err != nil {
		return nil, err
	}
	return res, nil
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

// CreateApp creates new app in BoltDB.
func (as *AppStorage) CreateApp(app model.AppData) (model.AppData, error) {
	res, ok := app.(*AppData)
	if !ok || app == nil {
		return nil, model.ErrorWrongDataFormat
	}
	result, err := as.addNewApp(*res)
	return result, err
}

// addNewApp adds new app to memory storage.
func (as *AppStorage) addNewApp(app model.AppData) (model.AppData, error) {
	res, ok := app.(AppData)
	if !ok {
		return nil, ErrorWrongDataFormat
	}
	// generate new ID if it's not set
	if len(res.ID()) == 0 {
		res.appData.ID = xid.New().String()
	}
	return res, as.db.Update(func(tx *bolt.Tx) error {
		data, err := res.Marshal()
		if err != nil {
			return err
		}

		ab := tx.Bucket([]byte(AppBucket))

		return ab.Put([]byte(res.ID()), data)
	})
}

// DisableApp disables app in the storage.
func (as *AppStorage) DisableApp(app model.AppData) error {
	res, ok := app.(AppData)
	if !ok {
		return ErrorWrongDataFormat
	}
	res.appData.Active = false
	_, err := as.addNewApp(res)
	return err
}

// UpdateApp updates app in the storage.
func (as *AppStorage) UpdateApp(oldAppID string, newApp model.AppData) error {
	res, ok := newApp.(AppData)
	if !ok {
		return ErrorWrongDataFormat
	}
	return as.db.Update(func(tx *bolt.Tx) error {
		data, err := res.Marshal()
		if err != nil {
			return err
		}

		ab := tx.Bucket([]byte(AppBucket))
		if err := ab.Delete([]byte(oldAppID)); err != nil {
			return err
		}

		return ab.Put([]byte(res.ID()), data)
	})
}

// FetchApps fetches apps which name satisfies provided filterString.
// Supports pagination.
func (as *AppStorage) FetchApps(filterString string, skip, limit int) ([]model.AppData, error) {
	var apps []model.AppData

	err := as.db.View(func(tx *bolt.Tx) error {
		ab := tx.Bucket([]byte(AppBucket))

		if iterErr := ab.ForEach(func(k, v []byte) error {
			if strings.Contains(strings.ToLower(string(k)), strings.ToLower(filterString)) {
				app, err := AppDataFromJSON(v)
				if err != nil {
					return err
				}
				apps = append(apps, app)
			}
			return nil
		}); iterErr != nil {
			return iterErr
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return apps, nil
}

// DeleteApp deletes app by ID.
func (as *AppStorage) DeleteApp(id string) error {
	err := as.db.Update(func(tx *bolt.Tx) error {
		ab := tx.Bucket([]byte(AppBucket))
		return ab.Delete([]byte(id))
	})
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
		if _, err := as.addNewApp(AppData{appData: a}); err != nil {
			return err
		}
	}
	return nil
}

type appData struct {
	ID                   string   `json:"id,omitempty"`
	Secret               string   `json:"secret,omitempty"`
	Active               bool     `json:"active,omitempty"`
	Name                 string   `json:"name,omitempty"`
	Description          string   `json:"description,omitempty"`
	Scopes               []string `json:"scopes,omitempty"`
	Offline              bool     `json:"offline,omitempty"`
	RedirectURL          string   `json:"redirect_url,omitempty"`
	RefreshTokenLifespan int64    `json:"refresh_token_lifespan,omitempty"`
	TokenLifespan        int64    `json:"token_lifespan,omitempty"`
	TokenPayload         []string `json:"token_payload,omitempty"`
}

// AppData is an in-memory model for model.AppData.
type AppData struct {
	appData
}

// NewAppData instantiates in-memory app data model from the general one.
func NewAppData(data model.AppData) AppData {
	return AppData{appData: appData{
		ID:                   data.ID(),
		Secret:               data.Secret(),
		Active:               data.Active(),
		Name:                 data.Name(),
		Description:          data.Description(),
		Scopes:               data.Scopes(),
		Offline:              data.Offline(),
		RedirectURL:          data.RedirectURL(),
		RefreshTokenLifespan: data.RefreshTokenLifespan(),
		TokenLifespan:        data.TokenLifespan(),
		TokenPayload:         data.TokenPayload(),
	}}
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

// MakeAppData creates new in-memory app data instance.
func MakeAppData(id, secret string, active bool, name, description string, scopes []string, offline bool, redirectURL string, refreshTokenLifespan, tokenLifespan int64, tokenPayload []string) AppData {
	return AppData{appData: appData{
		ID:                   id,
		Secret:               secret,
		Active:               active,
		Name:                 name,
		Description:          description,
		Scopes:               scopes,
		Offline:              offline,
		RedirectURL:          redirectURL,
		RefreshTokenLifespan: refreshTokenLifespan,
		TokenLifespan:        tokenLifespan,
		TokenPayload:         tokenPayload,
	}}
}

// Sanitize removes all sensitive data.
func (ad AppData) Sanitize() model.AppData {
	ad.appData.Secret = ""
	return ad
}

// ID implements model.AppData interface.
func (ad AppData) ID() string { return ad.appData.ID }

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
