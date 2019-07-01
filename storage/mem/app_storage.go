package mem

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/madappgang/identifo/model"
	"github.com/rs/xid"
)

// NewAppStorage creates new in-memory AppStorage implementation.
func NewAppStorage() (model.AppStorage, error) {
	return &AppStorage{storage: make(map[string]AppData)}, nil
}

// AppStorage is a fully functional app storage.
type AppStorage struct {
	storage map[string]AppData
}

// NewAppData returns pointer to newly created app data.
func (as *AppStorage) NewAppData() model.AppData {
	return &AppData{appData: appData{}}
}

// AppByID returns app by ID from the in-memory storage.
func (as *AppStorage) AppByID(id string) (model.AppData, error) {
	a, ok := as.storage[id]
	if !ok {
		return nil, ErrorNotFound
	}
	return &a, nil
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

// CreateApp creates new app in memory.
func (as *AppStorage) CreateApp(app model.AppData) (model.AppData, error) {
	res, ok := app.(*AppData)
	if !ok || res == nil {
		return nil, model.ErrorWrongDataFormat
	}
	result, err := as.addNewApp(res)
	return result, err
}

// addNewApp adds new app to in-memory storage.
func (as *AppStorage) addNewApp(app model.AppData) (model.AppData, error) {
	a, ok := app.(*AppData)
	if !ok || a == nil {
		return nil, model.ErrorWrongDataFormat
	}
	// generate new ID if it's not set
	if len(a.ID()) == 0 {
		a.appData.ID = xid.New().String()
	}
	as.storage[a.ID()] = NewAppData(app)
	return a, nil
}

// DisableApp deletes app from in-memory storage.
func (as *AppStorage) DisableApp(app model.AppData) error {
	delete(as.storage, app.ID())
	return nil
}

// UpdateApp updates app in the storage.
func (as *AppStorage) UpdateApp(appID string, newApp model.AppData) (model.AppData, error) {
	delete(as.storage, appID)
	updatedApp := NewAppData(newApp)
	as.storage[newApp.ID()] = updatedApp
	return &updatedApp, nil
}

// FetchApps fetches apps which name satisfies provided filterString.
// Supports pagination.
func (as *AppStorage) FetchApps(filterString string, skip, limit int) ([]model.AppData, int, error) {
	apps := []model.AppData{}
	var total int

	for _, app := range as.storage {
		total++
		skip--
		if skip > -1 {
			continue
		}
		if limit != 0 && len(apps) == limit {
			break
		}
		if strings.Contains(strings.ToLower(app.Name()), strings.ToLower(filterString)) {
			apps = append(apps, &app)
		}
	}
	return apps, total, nil
}

// DeleteApp does nothing here.
func (as *AppStorage) DeleteApp(id string) error {
	return nil
}

// TestDatabaseConnection is always optimistic about the database connection.
func (as *AppStorage) TestDatabaseConnection() error {
	return nil
}

// ImportJSON imports data from JSON.
func (as *AppStorage) ImportJSON(data []byte) error {
	apd := []appData{}
	if err := json.Unmarshal(data, &apd); err != nil {
		log.Println("Error unmarshalling app data:", err)
		return err
	}
	for _, a := range apd {
		if _, err := as.addNewApp(&AppData{appData: a}); err != nil {
			return err
		}
	}
	return nil
}

// Close clears storage.
func (as *AppStorage) Close() {
	for k := range as.storage {
		delete(as.storage, k)
	}
}

type appData struct {
	ID                    string           `json:"id,omitempty"`
	Secret                string           `json:"secret,omitempty"`
	Active                bool             `json:"active"`
	Name                  string           `json:"name,omitempty"`
	Description           string           `json:"description,omitempty"`
	Scopes                []string         `json:"scopes,omitempty"`
	Offline               bool             `json:"offline"`
	Type                  model.AppType    `json:"type,omitempty"`
	RedirectURL           string           `json:"redirect_url,omitempty"`
	RefreshTokenLifespan  int64            `json:"refresh_token_lifespan,omitempty"`
	InviteTokenLifespan   int64            `json:"invite_token_lifespan,omitempty"`
	TokenLifespan         int64            `json:"token_lifespan,omitempty"`
	TokenPayload          []string         `json:"token_payload,omitempty"`
	RegistrationForbidden bool             `json:"registration_forbidden"`
	AppleInfo             *model.AppleInfo `json:"apple_info,omitempty"`
}

// AppData is an in-memory model for model.AppData.
type AppData struct {
	appData
}

// NewAppData instantiates app data in-memory model from the general one.
func NewAppData(data model.AppData) AppData {
	return AppData{appData: appData{
		ID:                    data.ID(),
		Secret:                data.Secret(),
		Active:                data.Active(),
		Name:                  data.Name(),
		Description:           data.Description(),
		Scopes:                data.Scopes(),
		Offline:               data.Offline(),
		RedirectURL:           data.RedirectURL(),
		RefreshTokenLifespan:  data.RefreshTokenLifespan(),
		InviteTokenLifespan:   data.InviteTokenLifespan(),
		TokenLifespan:         data.TokenLifespan(),
		TokenPayload:          data.TokenPayload(),
		RegistrationForbidden: data.RegistrationForbidden(),
	}}
}

// MakeAppData creates new in-memory app data instance.
func MakeAppData(id, secret string, active bool, name, description string, scopes []string, offline bool, redirectURL string, refreshTokenLifespan, inviteTokenLifespan, tokenLifespan int64, tokenPayload []string, registrationForbidden bool) AppData {
	return AppData{appData: appData{
		ID:                    id,
		Secret:                secret,
		Active:                active,
		Name:                  name,
		Description:           description,
		Scopes:                scopes,
		Offline:               offline,
		RedirectURL:           redirectURL,
		RefreshTokenLifespan:  refreshTokenLifespan,
		InviteTokenLifespan:   inviteTokenLifespan,
		TokenLifespan:         tokenLifespan,
		TokenPayload:          tokenPayload,
		RegistrationForbidden: registrationForbidden,
	}}
}

// Sanitize removes all sensitive data.
func (ad *AppData) Sanitize() {
	if ad == nil {
		return
	}
	ad.appData.Secret = ""
	if ad.appData.AppleInfo != nil {
		ad.appData.AppleInfo.ClientSecret = ""
	}
}

// ID implements model.AppData interface.
func (ad *AppData) ID() string { return ad.appData.ID }

// Secret implements model.AppData interface.
func (ad *AppData) Secret() string { return ad.appData.Secret }

// Active implements model.AppData interface.
func (ad *AppData) Active() bool { return ad.appData.Active }

// Name implements model.AppData interface.
func (ad *AppData) Name() string { return ad.appData.Name }

// Description implements model.AppData interface.
func (ad *AppData) Description() string { return ad.appData.Description }

// Scopes implements model.AppData interface.
func (ad *AppData) Scopes() []string { return ad.appData.Scopes }

// Offline implements model.AppData interface.
func (ad *AppData) Offline() bool { return ad.appData.Offline }

// Type implements model.AppData interface.
func (ad *AppData) Type() model.AppType { return ad.appData.Type }

// RedirectURL implements model.AppData interface.
func (ad *AppData) RedirectURL() string { return ad.appData.RedirectURL }

// RefreshTokenLifespan implements model.AppData interface.
func (ad *AppData) RefreshTokenLifespan() int64 { return ad.appData.RefreshTokenLifespan }

// InviteTokenLifespan a inviteToken lifespan in seconds, if 0 - default one is used.
func (ad *AppData) InviteTokenLifespan() int64 { return ad.appData.InviteTokenLifespan }

// TokenLifespan implements model.AppData interface.
func (ad *AppData) TokenLifespan() int64 { return ad.appData.TokenLifespan }

// TokenPayload implements model.AppData interface.
func (ad *AppData) TokenPayload() []string { return ad.appData.TokenPayload }

// RegistrationForbidden implements model.AppData interface.
func (ad *AppData) RegistrationForbidden() bool { return ad.appData.RegistrationForbidden }

// AppleInfo implements model.AppData interface.
func (ad *AppData) AppleInfo() *model.AppleInfo { return ad.appData.AppleInfo }

// SetSecret implements model.AppData interface.
func (ad *AppData) SetSecret(secret string) {
	if ad == nil {
		return
	}
	ad.appData.Secret = secret
}
