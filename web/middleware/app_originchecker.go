package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/madappgang/identifo/v2/model"
)

type AppOriginChecker struct {
	OriginChecker
	apps model.AppStorage
}

// NewOriginChecker creates new instance of an OriginChecker.
func NewAppOriginChecker(appstorage model.AppStorage) (*AppOriginChecker, error) {
	originChecker := &AppOriginChecker{
		OriginChecker: OriginChecker{
			origins: make(map[string]bool),
			checks:  make([]model.OriginCheckFunc, 1),
		},
		apps: appstorage,
	}

	// function that grabs origins from the global origins map
	originChecker.checks[0] = func(r *http.Request, origin string) bool {
		return originChecker.IsPresent(r.Header.Get("Origin"))
	}

	if err := originChecker.Update(); err != nil {
		return nil, err
	}

	return originChecker, nil
}

func (aoc *AppOriginChecker) Update() error {
	if aoc.apps == nil {
		return errors.New("AppOriginChecker has no apps storage configured")
	}
	apps, err := aoc.apps.FetchApps("")
	if err != nil {
		return fmt.Errorf("error occurred during fetching apps: %s", err.Error())
	}

	allURLs := []string{}
	for _, a := range apps {
		allURLs = append(allURLs, a.RedirectURLs...)
	}
	aoc.DeleteAll()
	aoc.AddRawURLs(allURLs)
	return nil
}
