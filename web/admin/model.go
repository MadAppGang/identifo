package admin

import "github.com/madappgang/identifo/v2/model"

// request data to create new user form admin panel.
type registrationData struct {
	Username          string `json:"username"`
	Email             string `json:"email"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	MiddleName        string `json:"middle_name"`
	Nickname          string `json:"nickname"`
	PreferredUsername string `json:"preferred_username"`
	PhoneNumber       string `json:"phone_number"`
	Password          string `json:"password"`
}

// password reset data from admin panel.
type passwordResetData struct {
	UserID       string `json:"user_id,omitempty"`
	AppID        string `json:"app_id,omitempty"`
	ResetPageURL string `json:"reset_page_url,omitempty"`
}

// reset email data????
type resetEmailData struct {
	User  model.User
	Token string
	URL   string
	Host  string
}
