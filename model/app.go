package model

import "encoding/json"

// AppData represents Application data information.
type AppData struct {
	ID                   string               `bson:"_id" json:"id"`
	Secret               string               `bson:"secret" json:"secret"`
	Active               bool                 `bson:"active" json:"active"`
	Name                 string               `bson:"name" json:"name"`
	Description          string               `bson:"description" json:"description"`
	Scopes               []string             `bson:"scopes" json:"scopes"`   // Scopes is the list of all allowed scopes. If it's empty, no limitations (opaque scope).
	Offline              bool                 `bson:"offline" json:"offline"` // Offline is a boolean value that indicates whether on not the app supports refresh tokens. Do not use refresh tokens with apps that does not have secure storage.
	Type                 AppType              `bson:"type" json:"type"`
	RedirectURLs         []string             `bson:"redirect_urls" json:"redirect_urls"`           // RedirectURLs is the list of allowed urls where user will be redirected after successful login. Useful not only for web apps, mobile and desktop apps could use custom scheme for that.
	LoginAppSettings     *LoginWebAppSettings `bson:"login_app_settings" json:"login_app_settings"` // Rewrite login app settings for custom login, reset password and other settings
	CustomEmailTemplates bool                 `bson:"custom_email_templates" json:"custom_email_templates"`
	AuthStrategies       []AuthStrategy       `bson:"auth_strategies" json:"auth_strategies"`

	// map of map of custom sms message templates
	// root map is language map, the key is a language.Tag.String()
	// one special key is "default", which is language agnostic fall-back.
	// the second map is SMS message templates with a key of SMSMessageType.
	// to get message for OTPCode for english: CustomMessages["en"][SMSTypeOTPCode]
	CustomSMSMessages map[string]map[SMSMessageType]string `bson:"custom_sms_messages" json:"custom_sms_messages"`

	// registration settings
	RegistrationForbidden           bool   `bson:"registration_forbidden" json:"registration_forbidden"`
	PasswordlessRegistrationAllowed bool   `bson:"passwordless_registration_allowed" json:"passwordless_registration_allowed"`
	AnonymousRegistrationAllowed    bool   `bson:"anonymous_registration_allowed" json:"anonymous_registration_allowed"`
	NewUserDefaultRole              string `bson:"new_user_default_role" json:"new_user_default_role"`
	DebugOTPCodeAllowed             bool   `bson:"debug_otp_code_allowed" json:"debug_otp_code_allowed"`
	DebugOTPCodeForRegistration     string `bson:"debug_otp_code_for_registration" json:"debug_otp_code_for_registration"`
}

// AppType is a type of application.
type AppType string

const (
	Web     AppType = "web"     // Web is a web app.
	Android AppType = "android" // Android is an Android app.
	IOS     AppType = "ios"     // IOS is an iOS app.
	Desktop AppType = "desktop" // Desktop is a desktop app.
)

// AppDataFromJSON unmarshal AppData from JSON string
func AppDataFromJSON(d []byte) (AppData, error) {
	var apd AppData
	if err := json.Unmarshal(d, &apd); err != nil {
		return AppData{}, err
	}
	return apd, nil
}

func (a AppData) Sanitized() AppData {
	a.Secret = ""
	return a
}
