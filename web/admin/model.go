package admin

import "github.com/madappgang/identifo/v2/model"

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
