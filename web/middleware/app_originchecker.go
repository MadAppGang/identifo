package middleware

import (
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
	apps, _, err := aoc.apps.FetchApps("", 0, 0)
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
