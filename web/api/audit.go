package api

import "log"

func journal(userID, appID, action string, scopes []string) {
	// TODO: Create an interface for the audit log
	// Implement it for logging to stdout, a database, or a remote service
	log.Printf("audit record | %s | userID=%s appID=%s scopes=%v\n",
		action, userID, appID, scopes)
}
