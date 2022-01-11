package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/madappgang/identifo/v2/model"
)

// OriginCheckFunc the type of function that checks origin

// OriginChecker holds user's AllowOriginRequestFunc and checks dynamically
// added CORS origins from an all app's redirect urls.
type OriginChecker struct {
	sync.RWMutex
	origins map[string]bool
	// checks is a slice of AllowOriginRequestFuncs
	checks []model.OriginCheckFunc
}

// NewOriginChecker creates new instance of an OriginChecker.
func NewOriginChecker() *OriginChecker {
	originChecker := &OriginChecker{
		origins: make(map[string]bool),
		checks:  make([]model.OriginCheckFunc, 1),
	}

	// function that grabs origins from the global origins map
	originChecker.checks[0] = func(r *http.Request, origin string) bool {
		return originChecker.IsPresent(r.Header.Get("Origin"))
	}

	return originChecker
}

// NewOriginChecker creates new instance of an OriginChecker.
func NewOriginCheckerWithFunc(f model.OriginCheckFunc) model.OriginChecker {
	oc := NewOriginChecker()
	oc.AddCheck(f)
	return oc
}

func (os *OriginChecker) Update() error {
	// do nothing
	return nil
}

// IsPresent returns true if the provided origin presented in the origins map, false otherwise.
func (os *OriginChecker) IsPresent(origin string) bool {
	os.RLock()
	defer os.RUnlock()

	_, ok := os.origins[origin]
	return ok
}

// Add adds origin to the list of allowed origins.
func (os *OriginChecker) Add(origin string) {
	os.Lock()
	defer os.Unlock()

	os.origins[origin] = true
}

// AddRawURLs parses and adds urls to the list of allowed origins.
func (os *OriginChecker) AddRawURLs(urls []string) {
	os.Lock()
	defer os.Unlock()

	for _, u := range urls {
		parsed, err := url.Parse(u)
		if err == nil {
			os.origins[fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)] = true
		}
	}
}

// Delete removes origin from the list of allowed origins.
func (os *OriginChecker) Delete(origin string) {
	os.Lock()
	defer os.Unlock()

	delete(os.origins, origin)
}

// DeleteAll removes all origins from the global origin map.
func (os *OriginChecker) DeleteAll() {
	os.Lock()
	defer os.Unlock()

	os.origins = make(map[string]bool)
}

// With adds AllowOriginRequestFunc to list of checks.
func (os *OriginChecker) AddCheck(f model.OriginCheckFunc) {
	os.checks = append(os.checks, f)
}

// CheckOrigin is a custom func for validate origin, checking it with all AllowOriginRequestFuncs,
// including user's provided func.
func (os *OriginChecker) CheckOrigin(r *http.Request, origin string) bool {
	for _, check := range os.checks {
		if check != nil && check(r, origin) {
			return true
		}
	}
	return false
}
