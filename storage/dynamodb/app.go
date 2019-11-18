package dynamodb

import (
	"encoding/json"
	"log"

	"github.com/madappgang/identifo/model"
	"github.com/rs/xid"
)

// AppData is DynamoDB model for model.AppData.
type AppData struct {
	appData
}

type appData struct {
	ID                           string                 `json:"id,omitempty"`
	Secret                       string                 `json:"secret,omitempty"`
	Active                       bool                   `json:"active"`
	Name                         string                 `json:"name,omitempty"`
	Description                  string                 `json:"description,omitempty"`
	Scopes                       []string               `json:"scopes,omitempty"`
	Offline                      bool                   `json:"offline"`
	Type                         model.AppType          `json:"type,omitempty"`
	RedirectURLs                 []string               `json:"redirect_urls,omitempty"`
	RefreshTokenLifespan         int64                  `json:"refresh_token_lifespan,omitempty"`
	InviteTokenLifespan          int64                  `json:"invite_token_lifespan,omitempty"`
	TokenLifespan                int64                  `json:"token_lifespan,omitempty"`
	TokenPayload                 []string               `json:"token_payload,omitempty"`
	TFAStatus                    model.TFAStatus        `json:"tfa_status"`
	DebugTFACode                 string                 `json:"debug_tfa_code,omitempty"`
	RegistrationForbidden        bool                   `json:"registration_forbidden"`
	AnonymousRegistrationAllowed bool                   `json:"anonymous_registration_allowed"`
	AuthorizationWay             model.AuthorizationWay `json:"authorization_way,omitempty"`
	AuthorizationModel           string                 `json:"authorization_model,omitempty"`
	AuthorizationPolicy          string                 `json:"authorization_policy,omitempty"`
	RolesWhitelist               []string               `json:"roles_whitelist,omitempty"`
	RolesBlacklist               []string               `json:"roles_blacklist,omitempty"`
	NewUserDefaultRole           string                 `json:"new_user_default_role,omitempty"`
	AppleInfo                    *model.AppleInfo       `json:"apple_info,omitempty"`
}

// NewAppData instantiates DynamoDB app data model from the general one.
func NewAppData(data model.AppData) (AppData, error) {
	if _, err := xid.FromString(data.ID()); err != nil {
		log.Println("Incorrect AppID: ", data.ID())
		return AppData{}, model.ErrorWrongDataFormat
	}
	return AppData{appData: appData{
		ID:                           data.ID(),
		Secret:                       data.Secret(),
		Active:                       data.Active(),
		Name:                         data.Name(),
		Description:                  data.Description(),
		Scopes:                       data.Scopes(),
		Offline:                      data.Offline(),
		RedirectURLs:                 data.RedirectURLs(),
		RefreshTokenLifespan:         data.RefreshTokenLifespan(),
		InviteTokenLifespan:          data.InviteTokenLifespan(),
		TokenLifespan:                data.TokenLifespan(),
		TokenPayload:                 data.TokenPayload(),
		RegistrationForbidden:        data.RegistrationForbidden(),
		AnonymousRegistrationAllowed: data.AnonymousRegistrationAllowed(),
	}}, nil
}

// AppDataFromJSON deserializes data from JSON.
func AppDataFromJSON(d []byte) (AppData, error) {
	apd := appData{}
	if err := json.Unmarshal(d, &apd); err != nil {
		log.Println(err)
		return AppData{}, err
	}
	return AppData{appData: apd}, nil
}

// MakeAppData creates new DynamoDB app data instance.
func MakeAppData(id, secret string, active bool, name, description string, scopes []string, offline bool, redirectURLs []string,
	refreshTokenLifespan, inviteTokenLifespan, tokenLifespan int64, tokenPayload []string, registrationForbidden bool, anonymousRegistrationAllowed bool,
	tfaStatus model.TFAStatus, debugTFACode string, authzWay model.AuthorizationWay, authzModel, authzPolicy string, rolesWhitelist, rolesBlacklist []string, newUserDefaultRole string) (AppData, error) {
	if _, err := xid.FromString(id); err != nil {
		log.Println("Cannot create ID from the string representation:", err)
		return AppData{}, model.ErrorWrongDataFormat
	}
	return AppData{appData: appData{
		ID:                           id,
		Secret:                       secret,
		Active:                       active,
		Name:                         name,
		Description:                  description,
		Scopes:                       scopes,
		Offline:                      offline,
		RedirectURLs:                 redirectURLs,
		RefreshTokenLifespan:         refreshTokenLifespan,
		InviteTokenLifespan:          inviteTokenLifespan,
		TokenLifespan:                tokenLifespan,
		TokenPayload:                 tokenPayload,
		RegistrationForbidden:        registrationForbidden,
		AnonymousRegistrationAllowed: anonymousRegistrationAllowed,
		TFAStatus:                    tfaStatus,
		DebugTFACode:                 debugTFACode,
		AuthorizationWay:             authzWay,
		AuthorizationModel:           authzModel,
		AuthorizationPolicy:          authzPolicy,
		RolesWhitelist:               rolesWhitelist,
		RolesBlacklist:               rolesBlacklist,
		NewUserDefaultRole:           newUserDefaultRole,
	}}, nil
}

// Marshal serializes data to byte array.
func (ad AppData) Marshal() ([]byte, error) {
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

// RedirectURLs implements model.AppData interface.
func (ad *AppData) RedirectURLs() []string { return ad.appData.RedirectURLs }

// RefreshTokenLifespan implements model.AppData interface.
func (ad *AppData) RefreshTokenLifespan() int64 { return ad.appData.RefreshTokenLifespan }

// InviteTokenLifespan a inviteToken lifespan in seconds, if 0 - default one is used.
func (ad *AppData) InviteTokenLifespan() int64 { return ad.appData.InviteTokenLifespan }

// TokenLifespan implements model.AppData interface.
func (ad *AppData) TokenLifespan() int64 { return ad.appData.TokenLifespan }

// TokenPayload implements model.AppData interface.
func (ad *AppData) TokenPayload() []string { return ad.appData.TokenPayload }

// TFAStatus implements model.AppData interface.
func (ad *AppData) TFAStatus() model.TFAStatus { return ad.appData.TFAStatus }

// DebugTFACode implements model.AppData interface.
func (ad *AppData) DebugTFACode() string { return ad.appData.DebugTFACode }

// RegistrationForbidden implements model.AppData interface.
func (ad *AppData) RegistrationForbidden() bool { return ad.appData.RegistrationForbidden }

// AnonymousRegistrationAllowed implements model.AppData interface.
func (ad *AppData) AnonymousRegistrationAllowed() bool { return ad.appData.AnonymousRegistrationAllowed }

// AuthzWay implements model.AppData interface.
func (ad *AppData) AuthzWay() model.AuthorizationWay { return ad.appData.AuthorizationWay }

// AuthzModel implements model.AppData interface.
func (ad *AppData) AuthzModel() string { return ad.appData.AuthorizationModel }

// AuthzPolicy implements model.AppData interface.
func (ad *AppData) AuthzPolicy() string { return ad.appData.AuthorizationPolicy }

// RolesWhitelist implements model.AppData interface.
func (ad *AppData) RolesWhitelist() []string { return ad.appData.RolesWhitelist }

// RolesBlacklist implements model.AppData interface.
func (ad *AppData) RolesBlacklist() []string { return ad.appData.RolesBlacklist }

// NewUserDefaultRole implements model.AppData interface.
func (ad *AppData) NewUserDefaultRole() string { return ad.appData.NewUserDefaultRole }

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
