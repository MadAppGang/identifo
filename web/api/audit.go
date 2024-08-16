package api

import "log"

func journal(userID, appID, action string, scopes []string) {
	log.Printf("audit record | %s | userID=%s appID=%s scopes=%v\n",
		action, userID, appID, scopes)
}
