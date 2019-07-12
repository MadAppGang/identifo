package api

import (
	"fmt"
	"net/http"
	"strings"
)

const anonymousRole = "anonymous"

type authzInfo struct {
	appID       string
	tokenStr    string
	resourceURI string
	method      string
}

// Authorize checks if user has an access to the requested resource.
// If error happens, writes it to ResponseWriter.
// Also, writes an error on failed authorization.
func (ar *Router) Authorize(w http.ResponseWriter, azi authzInfo) error {
	if ar.Authorizer == nil {
		return nil
	}

	userID, err := ar.getTokenSubject(azi.tokenStr)
	if err != nil {
		err = fmt.Errorf("Error getting subject from token: %s", err)
		ar.logger.Println(err)
		ar.Error(w, ErrorAPIAppCannotExtractTokenSubject, http.StatusBadRequest, err.Error(), "Authorizer.GetTokenSubject")
		return err
	}

	user, err := ar.userStorage.UserByID(userID)
	if err != nil {
		err = fmt.Errorf("Error getting user by ID: %s", err)
		ar.logger.Println(err)
		ar.Error(w, ErrorAPIUserNotFound, http.StatusUnauthorized, err.Error(), "Authorizer.UserByID")
		return err
	}

	sub := user.Role()
	if sub == "" {
		sub = anonymousRole
	}
	obj := strings.Join([]string{azi.appID, azi.resourceURI}, ":")
	act := azi.method

	accessGranted, err := ar.Authorizer.Enforce(sub, obj, act)
	if err != nil {
		err = fmt.Errorf("Error calling enforcer: %s", err)
		ar.logger.Println(err)
		ar.Error(w, ErrorAPIUserNotFound, http.StatusInternalServerError, err.Error(), "Authorizer.EnforceSafe")
		return err
	}
	if !accessGranted {
		err := fmt.Errorf("Access denied")
		ar.Error(w, ErrorAPIAppAccessDenied, http.StatusForbidden, err.Error(), "Authorizer.AccessDenied")
		return err
	}
	return nil
}
