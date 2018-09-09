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
	if a, ok := as.storage[id]; ok != true {
		return nil, ErrorNotFound
	} else {
		return a, nil
	}
}

//AddNewApp add new app to memory storage
func (as *AppStorage) AddNewApp(app model.AppData) error {
	as.storage[app.ID()] = NewAppData(app)
	return nil
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

//AppData is memory model for model.AppData
type AppData struct {
	id          string
	secret      string
	active      bool
	description string
	scopes      []string
}

//NewAppData instantiate app data memory model from general one
func NewAppData(data model.AppData) AppData {
	return AppData{
		id:          data.ID(),
		secret:      data.Secret(),
		active:      data.Active(),
		description: data.Description(),
		scopes:      data.Scopes(),
	}
}

//MakeAppData creates new memory app data instance
func MakeAppData(id, secret string, active bool, description string, scopes []string) AppData {
	return AppData{id, secret, active, description, scopes}

}

func (ad AppData) ID() string          { return ad.id }
func (ad AppData) Secret() string      { return ad.secret }
func (ad AppData) Active() bool        { return ad.active }
func (ad AppData) Description() string { return ad.description }
func (ad AppData) Scopes() []string    { return ad.scopes }
