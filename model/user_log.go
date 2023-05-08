package model

import "time"

type UserLogEntity struct {
	ID        string
	timestamp time.Time
	UserID    string
	UserEmail string
	Event     UserLogEntityType
	Data      map[string]any
}

// some log entities
type UserLogEntityType string

const (
	UserLogEntityRegistration       UserLogEntityType = "registration"
	UserLogEntityTypeLoginSuccess   UserLogEntityType = "login_success"
	UserLogEntityTypeLoginFailure   UserLogEntityType = "login_failure"
	UserLogEntityTypeLogout         UserLogEntityType = "logout"
	UserLogEntityTypePasswordReset  UserLogEntityType = "password_reset"
	UserLogEntityTypePasswordChange UserLogEntityType = "password_change"
	UserLogEntityTypeEmailChange    UserLogEntityType = "email_change"
	UserLogEntityTypeProfileUpdate  UserLogEntityType = "profile_update"
	UserLogEntityTypeDelete         UserLogEntityType = "delete"
	UserLogEntityTypeBlocked        UserLogEntityType = "blocked"
	UserLogEntityTypeUnblocked      UserLogEntityType = "unblocked"
	UserLogEntityTypeMFAEnrolled    UserLogEntityType = "mfa_enrolled"
)
