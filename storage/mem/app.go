package mem

import (
	"github.com/madappgang/identifo/model"
)

// AppData is an in-memory model for model.AppData.
type AppData struct {
	appData
}

type appData struct {
	ID                    string                 `json:"id,omitempty"`
	Secret                string                 `json:"secret,omitempty"`
	Active                bool                   `json:"active"`
	Name                  string                 `json:"name,omitempty"`
	Description           string                 `json:"description,omitempty"`
	Scopes                []string               `json:"scopes,omitempty"`
	Offline               bool                   `json:"offline"`
	Type                  model.AppType          `json:"type,omitempty"`
	RedirectURL           string                 `json:"redirect_url,omitempty"`
	RefreshTokenLifespan  int64                  `json:"refresh_token_lifespan,omitempty"`
	InviteTokenLifespan   int64                  `json:"invite_token_lifespan,omitempty"`
	TokenLifespan         int64                  `json:"token_lifespan,omitempty"`
	TokenPayload          []string               `json:"token_payload,omitempty"`
	RegistrationForbidden bool                   `json:"registration_forbidden"`
	AuthorizationWay      model.AuthorizationWay `json:"authorization_way,omitempty"`
	AuthorizationModel    string                 `json:"authorization_model,omitempty"`
	AuthorizationPolicy   string                 `json:"authorization_policy,omitempty"`
	AppleInfo             *model.AppleInfo       `json:"apple_info,omitempty"`
}

// NewAppData instantiates app data in-memory model from the general one.
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

// MakeAppData creates new in-memory app data instance.
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

// AuthzWay implements model.AppData interface.
func (ad *AppData) AuthzWay() model.AuthorizationWay { return ad.appData.AuthorizationWay }

// AuthzModel implements model.AppData interface.
func (ad *AppData) AuthzModel() string { return ad.appData.AuthorizationModel }

// AuthzPolicy implements model.AppData interface.
func (ad *AppData) AuthzPolicy() string { return ad.appData.AuthorizationPolicy }

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

	ad.appData.AuthorizationWay = ""
	ad.appData.AuthorizationModel = ""
	ad.appData.AuthorizationPolicy = ""
}
