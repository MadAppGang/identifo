package mem

import "github.com/madappgang/identifo/model"

//NewAppStorage creates new memory AppStorage implementation
func NewAppStorage() model.AppStorage {
	as := AppStorage{}
	as.storage = make(map[string]AppData)
	return &as
}

//AppStorage is fully functional app storage
type AppStorage struct {
	storage map[string]AppData
}

//AppByID returns app from memory by ID
func (as *AppStorage) AppByID(id string) (model.AppData, error) {
	if a, ok := as.storage[id]; !ok {
		return nil, ErrorNotFound
	} else {
		return a, nil
	}
}

//AddNewApp add new app to memory storage
func (as *AppStorage) AddNewApp(app model.AppData) (model.AppData, error) {
	as.storage[app.ID()] = NewAppData(app)
	return app, nil
}

//DisableApp disables app from storage
func (as *AppStorage) DisableApp(app model.AppData) error {
	delete(as.storage, app.ID())
	return nil
}

//UpdateApp updates app in storage
func (as *AppStorage) UpdateApp(oldAppID string, newApp model.AppData) error {
	delete(as.storage, oldAppID)
	as.storage[newApp.ID()] = NewAppData(newApp)
	return nil
}

type appData struct {
	ID                         string   `json:"id,omitempty"`
	Secret                     string   `json:"secret,omitempty"`
	Active                     bool     `json:"active,omitempty"`
	Description                string   `json:"description,omitempty"`
	Scopes                     []string `json:"scopes,omitempty"`
	Offline                    bool     `json:"offline,omitempty"`
	RedirectURL                string   `json:"redirect_url,omitempty"`
	RefreshTokenLifespan       int64    `json:"refresh_token_lifespan,omitempty"`
	ResetPasswordTokenLifespan int64    `json:"reset_password_token_lifespan,omitempty"`
	TokenLifespan              int64    `json:"token_lifespan,omitempty"`
}

//AppData is memory model for model.AppData
type AppData struct {
	appData
}

//NewAppData instantiate app data memory model from general one
func NewAppData(data model.AppData) AppData {
	return AppData{appData: appData{
		ID:                         data.ID(),
		Secret:                     data.Secret(),
		Active:                     data.Active(),
		Description:                data.Description(),
		Scopes:                     data.Scopes(),
		Offline:                    data.Offline(),
		RedirectURL:                data.RedirectURL(),
		RefreshTokenLifespan:       data.RefreshTokenLifespan(),
		ResetPasswordTokenLifespan: data.ResetPasswordTokenLifespan(),
		TokenLifespan:              data.TokenLifespan(),
	}}
}

//MakeAppData creates new memory app data instance
func MakeAppData(id, secret string, active bool, description string, scopes []string, offline bool, redirectURL string, refreshTokenLifespan, resetPasswordTokenLifespan, tokenLifespan int64) AppData {
	return AppData{appData: appData{
		ID:                         id,
		Secret:                     secret,
		Active:                     active,
		Description:                description,
		Scopes:                     scopes,
		Offline:                    offline,
		RedirectURL:                redirectURL,
		RefreshTokenLifespan:       refreshTokenLifespan,
		ResetPasswordTokenLifespan: resetPasswordTokenLifespan,
		TokenLifespan:              tokenLifespan,
	}}
}

func (ad AppData) ID() string                        { return ad.appData.ID }
func (ad AppData) Secret() string                    { return ad.appData.Secret }
func (ad AppData) Active() bool                      { return ad.appData.Active }
func (ad AppData) Description() string               { return ad.appData.Description }
func (ad AppData) Scopes() []string                  { return ad.appData.Scopes }
func (ad AppData) Offline() bool                     { return ad.appData.Offline }
func (ad AppData) RedirectURL() string               { return ad.appData.RedirectURL }
func (ad AppData) RefreshTokenLifespan() int64       { return ad.appData.RefreshTokenLifespan }
func (ad AppData) ResetPasswordTokenLifespan() int64 { return ad.appData.ResetPasswordTokenLifespan }
func (ad AppData) TokenLifespan() int64              { return ad.appData.TokenLifespan }

//AddAppDataFromFile loads appdata from JSON file and save it to the storage
func AddAppDataFromFile(as model.AppStorage, file string) {

}
