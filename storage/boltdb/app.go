package boltdb

import (
	"encoding/json"

	"github.com/madappgang/identifo/model"
)

// AppData is an in-memory model for model.AppData.
type AppData struct {
	appData
}

type appData struct {
	ID                    string           `json:"id,omitempty"`
	Secret                string           `json:"secret,omitempty"`
	Active                bool             `json:"active"`
	Name                  string           `json:"name,omitempty"`
	Description           string           `json:"description,omitempty"`
	Scopes                []string         `json:"scopes,omitempty"`
	Offline               bool             `json:"offline"`
	Type                  model.AppType    `json:"type,omitempty"`
	RedirectURL           string           `json:"redirect_url,omitempty"`
	RefreshTokenLifespan  int64            `json:"refresh_token_lifespan,omitempty"`
	InviteTokenLifespan   int64            `json:"invite_token_lifespan,omitempty"`
	TokenLifespan         int64            `json:"token_lifespan,omitempty"`
	TokenPayload          []string         `json:"token_payload,omitempty"`
	TFAEnabled            bool             `json:"tfa_enabled"`
	RegistrationForbidden bool             `json:"registration_forbidden"`
	AppleInfo             *model.AppleInfo `json:"apple_info,omitempty"`
}

// NewAppData instantiates in-memory app data model from the general one.
func NewAppData(data model.AppData) AppData {
	return AppData{appData: appData{
		ID:                    data.ID(),
		Secret:                data.Secret(),
		Active:                data.Active(),
		Name:                  data.Name(),
		Description:           data.Description(),
		Scopes:                data.Scopes(),
		Offline:               data.Offline(),
		RedirectURL:           data.RedirectURL(),
		RefreshTokenLifespan:  data.RefreshTokenLifespan(),
		InviteTokenLifespan:   data.InviteTokenLifespan(),
		TokenLifespan:         data.TokenLifespan(),
		TokenPayload:          data.TokenPayload(),
		RegistrationForbidden: data.RegistrationForbidden(),
	}}
}

// AppDataFromJSON deserializes app data from JSON.
func AppDataFromJSON(d []byte) (*AppData, error) {
	apd := appData{}
	if err := json.Unmarshal(d, &apd); err != nil {
		return &AppData{}, err
	}
	return &AppData{appData: apd}, nil
}

// MakeAppData creates new app data instance.
func MakeAppData(id, secret string, active bool, name, description string, scopes []string, offline bool, redirectURL string, refreshTokenLifespan, inviteTokenLifespan, tokenLifespan int64, tokenPayload []string, registrationForbidden bool) AppData {
	return AppData{appData: appData{
		ID:                    id,
		Secret:                secret,
		Active:                active,
		Name:                  name,
		Description:           description,
		Scopes:                scopes,
		Offline:               offline,
		RedirectURL:           redirectURL,
		RefreshTokenLifespan:  refreshTokenLifespan,
		InviteTokenLifespan:   inviteTokenLifespan,
		TokenLifespan:         tokenLifespan,
		TokenPayload:          tokenPayload,
		RegistrationForbidden: registrationForbidden,
	}}
}

// Marshal serializes data to byte array.
func (ad *AppData) Marshal() ([]byte, error) {
	return json.Marshal(ad.appData)
}

// ID implements model.AppData interface.
func (ad *AppData) ID() string { return ad.appData.ID }

// Secret implements model.AppData interface.
func (ad *AppData) Secret() string { return ad.appData.Secret }

// Active implements model.AppData interface.
func (ad *AppData) Active() bool { return ad.appData.Active }

// Name implements model.AppData interface.
func (ad *AppData) Name() string { return ad.appData.Name }

// Description implements model.AppData interface.
func (ad *AppData) Description() string { return ad.appData.Description }

// Scopes implements model.AppData interface.
func (ad *AppData) Scopes() []string { return ad.appData.Scopes }

// Offline implements model.AppData interface.
func (ad *AppData) Offline() bool { return ad.appData.Offline }

// Type implements model.AppData interface.
func (ad *AppData) Type() model.AppType { return ad.appData.Type }

// RedirectURL implements model.AppData interface.
func (ad *AppData) RedirectURL() string { return ad.appData.RedirectURL }

// RefreshTokenLifespan implements model.AppData interface.
func (ad *AppData) RefreshTokenLifespan() int64 { return ad.appData.RefreshTokenLifespan }

// InviteTokenLifespan a inviteToken lifespan in seconds, if 0 - default one is used.
func (ad *AppData) InviteTokenLifespan() int64 { return ad.appData.InviteTokenLifespan }

// TokenLifespan implements model.AppData interface.
func (ad *AppData) TokenLifespan() int64 { return ad.appData.TokenLifespan }

// TokenPayload implements model.AppData interface.
func (ad *AppData) TokenPayload() []string { return ad.appData.TokenPayload }

// RegistrationForbidden implements model.AppData interface.
func (ad *AppData) RegistrationForbidden() bool { return ad.appData.RegistrationForbidden }

// TFAEnabled implements model.AppData interface.
func (ad *AppData) TFAEnabled() bool { return ad.appData.TFAEnabled }

// AppleInfo implements model.AppData interface.
func (ad *AppData) AppleInfo() *model.AppleInfo { return ad.appData.AppleInfo }

// SetSecret implements model.AppData interface.
func (ad *AppData) SetSecret(secret string) {
	if ad == nil {
		return
	}
	ad.appData.Secret = secret
}

// Sanitize removes all sensitive data.
func (ad *AppData) Sanitize() {
	if ad == nil {
		return
	}
	ad.appData.Secret = ""
	if ad.appData.AppleInfo != nil {
		ad.appData.AppleInfo.ClientSecret = ""
	}
}
