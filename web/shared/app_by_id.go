package shared

import (
	"errors"

	"github.com/madappgang/identifo/model"
)

var (
	errorInactiveApp = errors.New("App is inactive")
	errorEmptyAppID  = errors.New("Empty appID param")
)

// AppByID gets app by id and checks if it's active.
func AppByID(as model.AppStorage, appID string) (model.AppData, error) {
	if appID == "" {
		return nil, errorEmptyAppID
	}

	app, err := as.AppByID(appID)
	if err != nil {
		return nil, err
	}

	if !app.Active() {
		return nil, errorInactiveApp
	}

	return app, nil
}
