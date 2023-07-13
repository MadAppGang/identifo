package model

// LoginWebAppSettings settings for login web app: usually it is a SPA app which handles login, register, reset password, etc.
type LoginWebAppSettings struct {
	LoginURL           string `json:"login_url,omitempty" bson:"login_url,omitempty"`
	RegisterURL        string `json:"register_url,omitempty" bson:"register_url,omitempty"`
	ResetPasswordURL   string `json:"reset_password_url,omitempty" bson:"reset_password_url,omitempty"`
	ConfirmEmailURL    string `json:"confirm_email_url,omitempty" bson:"confirm_email_url,omitempty"`
	ErrorURL           string `json:"error_url,omitempty" bson:"error_url,omitempty"`
	WelcomePageURL     string `json:"welcome_page_url,omitempty" bson:"welcome_page_url,omitempty"`
	OTPConfirmationURL string `json:"otp_confirmation_url,omitempty" bson:"otp_confirmation_url,omitempty"`
}

// DefaultLoginWebAppSettings default settings for self-hosted SPA login app by Identifo.
var DefaultLoginWebAppSettings = LoginWebAppSettings{
	LoginURL:           "/web",
	RegisterURL:        "/web/register",
	ResetPasswordURL:   "/web/password/reset",
	ConfirmEmailURL:    "/web/email_confirm",
	ErrorURL:           "/web/misconfiguration",
	WelcomePageURL:     "/web/welcome",
	OTPConfirmationURL: "/web/otp_confirm",
}
