package originchecker

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

// origins is global for the origin_checker package to be used across all instances of an OriginChecker
var origins = make(map[string]struct{})

// OriginChecker holds user's AllowOriginRequestFunc and checks dynamically
// added CORS origins from an all app's redirect urls.
type OriginChecker struct {
	sync.RWMutex
	// checks is a slice of AllowOriginRequestFuncs
	checks []func(r *http.Request, origin string) bool
}

// NewOriginChecker creates new instance of an OriginChecker.
func NewOriginChecker() *OriginChecker {
	originChecker := &OriginChecker{
		checks: make([]func(r *http.Request, origin string) bool, 1),
	}

	// function that grabs origins from the global origins map
	originChecker.checks[0] = func(r *http.Request, origin string) bool {
		return originChecker.IsPresent(r.Header.Get("Origin"))
	}

	return originChecker
}

// IsPresent returns true if the provided origin presented in the origins map, false otherwise.
func (os *OriginChecker) IsPresent(origin string) bool {
	os.RLock()
	defer os.RUnlock()

	_, ok := origins[origin]
	return ok
}

// Add adds origin to the list of allowed origins.
func (os *OriginChecker) Add(origin string) {
	os.Lock()
	defer os.Unlock()

	origins[origin] = struct{}{}
}

// AddRawURLs parses and adds urls to the list of allowed origins.
func (os *OriginChecker) AddRawURLs(urls []string) {
	os.Lock()
	defer os.Unlock()

	for _, u := range urls {
		parsed, err := url.Parse(u)
		if err == nil {
			origins[fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)] = struct{}{}
		}
	}
}

// Delete removes origin from the list of allowed origins.
func (os *OriginChecker) Delete(origin string) {
	os.Lock()
	defer os.Unlock()

	delete(origins, origin)
}

// DeleteAll removes all origins from the global origin map.
func (os *OriginChecker) DeleteAll() {
	os.Lock()
	defer os.Unlock()

	origins = make(map[string]struct{})
}

// With adds AllowOriginRequestFunc to list of checks.
func (os *OriginChecker) With(f func(r *http.Request, origin string) bool) *OriginChecker {
	if f != nil {
		newChecker := NewOriginChecker()
		newChecker.checks = append(newChecker.checks, f)
		return newChecker
	}
	return os
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
