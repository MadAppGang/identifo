package mock

import (
	"crypto/rand"
	"io"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"golang.org/x/exp/maps"
)

var (
	appNotFoundError = l.NewError(l.ErrorNotFound, "app")

	_a model.AppStorage = &App{}
)

// App is mock implementation of model.AppStorage interface.
// ! User it for tests only!
type App struct {
	Apps map[string]model.AppData
}

func NewApp() *App {
	a := App{
		Apps: map[string]model.AppData{},
	}
	return &a
}

func (a *App) AppByID(id string) (model.AppData, error) {
	for _, app := range a.Apps {
		if app.ID == id {
			return app, nil
		}
	}
	return model.AppData{}, appNotFoundError
}

func (a *App) ActiveAppByID(appID string) (model.AppData, error) {
	app, err := a.AppByID(appID)
	if err != nil {
		return model.AppData{}, err
	}
	if !app.Active {
		return model.AppData{}, appNotFoundError
	}
	return app, nil
}

func (a *App) CreateApp(app model.AppData) (model.AppData, error) {
	if len(app.ID) == 0 {
		app.ID = randomID(10)
	}
	a.Apps[app.ID] = app
	return app, nil
}

func (a *App) DisableApp(app model.AppData) error {
	app, err := a.AppByID(app.ID)
	if err != nil {
		return err
	}
	app.Active = false
	a.Apps[app.ID] = app
	return nil
}

func (a *App) UpdateApp(appID string, newApp model.AppData) (model.AppData, error) {
	a.Apps[appID] = newApp
	return newApp, nil
}

func (a *App) FetchApps(filter string) ([]model.AppData, error) {
	return maps.Values(a.Apps), nil
}

func (a *App) DeleteApp(id string) error {
	delete(a.Apps, id)
	return nil
}

func (a *App) ImportJSON(data []byte, cleanOldData bool) error {
	// TODO: implement
	return l.ErrorNotImplemented
}

func (a *App) TestDatabaseConnection() error {
	return nil
}

func (a *App) Close() {
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func randomID(length int) string {
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
