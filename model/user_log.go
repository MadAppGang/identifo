package model

import "time"

// UserLogEntity is activity atomic action for user auth process.
type UserLogEntity struct {
	ID        string            `json:"id,omitempty"`
	UserID    string            `json:"user_id,omitempty"`
	Email     string            `json:"email,omitempty"`
	DeviceID  string            `json:"device_id,omitempty"`
	Phone     string            `json:"phone,omitempty"`
	Event     UserLogEntityType `json:"event,omitempty"`
	Data      map[string]any    `json:"data,omitempty"`
	Timestamp time.Time         `json:"timestamp,omitempty"`
}

// UserLogEntityType is a type of log for later filter, use structured logs 💪🏽.
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

// UserLifecycleEventSource dedicate user activity regarding the API interaction
type UserLifecycleEventSource string

const (
	UserLifecycleEventSourceUser          UserLifecycleEventSource = "user"
	UserLifecycleEventSourceAdmin         UserLifecycleEventSource = "admin"
	UserLifecycleEventSourceManagementAPI UserLifecycleEventSource = "management_api"
	UserLifecycleEventSourceUnknown       UserLifecycleEventSource = "unknown"
)
