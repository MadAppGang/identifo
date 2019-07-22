package mongo

import (
	"encoding/json"

	"github.com/madappgang/identifo/model"
	"gopkg.in/mgo.v2/bson"
)

// AppData is a MongoDb model that implements model.AppData.
type AppData struct {
	appData
}

type appData struct {
	ID                    bson.ObjectId          `bson:"_id,omitempty" json:"id,omitempty"`
	Secret                string                 `bson:"secret,omitempty" json:"secret,omitempty"`
	Active                bool                   `bson:"active" json:"active"`
	Name                  string                 `bson:"name,omitempty" json:"name,omitempty"`
	Description           string                 `bson:"description,omitempty" json:"description,omitempty"`
	Scopes                []string               `bson:"scopes,omitempty" json:"scopes,omitempty"`
	Offline               bool                   `bson:"offline" json:"offline"`
	Type                  model.AppType          `bson:"type,omitempty" json:"type,omitempty"`
	RedirectURL           string                 `bson:"redirect_url,omitempty" json:"redirect_url,omitempty"`
	RefreshTokenLifespan  int64                  `bson:"refresh_token_lifespan,omitempty" json:"refresh_token_lifespan,omitempty"`
	InviteTokenLifespan   int64                  `bson:"invite_token_lifespan,omitempty" json:"invite_token_lifespan,omitempty"`
	TokenLifespan         int64                  `bson:"token_lifespan,omitempty" json:"token_lifespan,omitempty"`
	TokenPayload          []string               `bson:"token_payload,omitempty" json:"token_payload,omitempty"`
	RegistrationForbidden bool                   `bson:"registration_forbidden" json:"registration_forbidden"`
	AuthorizationWay      model.AuthorizationWay `bson:"authorization_way,omitempty" json:"authorization_way,omitempty"`
	AuthorizationModel    string                 `bson:"authorization_model,omitempty" json:"authorization_model,omitempty"`
	AuthorizationPolicy   string                 `bson:"authorization_policy,omitempty" json:"authorization_policy,omitempty"`
	AppleInfo             *model.AppleInfo       `bson:"apple_info,omitempty" json:"apple_info,omitempty"`
}

// NewAppData instantiates MongoDB app data model from the general one.
func NewAppData(data model.AppData) (AppData, error) {
	if !bson.IsObjectIdHex(data.ID()) {
		return AppData{}, model.ErrorWrongDataFormat
	}
	return AppData{appData: appData{
		ID:                    bson.ObjectIdHex(data.ID()),
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
	}}, nil
}

// AppDataFromJSON deserializes app data from JSON.
func AppDataFromJSON(d []byte) (AppData, error) {
	apd := appData{}
	if err := json.Unmarshal(d, &apd); err != nil {
		return AppData{}, err
	}
	return AppData{appData: apd}, nil
}

// MakeAppData creates new MongoDB app data instance.
func MakeAppData(id, secret string, active bool, name, description string, scopes []string, offline bool, redirectURL string, refreshTokenLifespan, inviteTokenLifespan, tokenLifespan int64, tokenPayload []string, registrationForbidden bool) (AppData, error) {
	if !bson.IsObjectIdHex(id) {
		return AppData{}, model.ErrorWrongDataFormat
	}
	return AppData{appData: appData{
		ID:                    bson.ObjectIdHex(id),
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
	}}, nil
}

// Marshal serializes data to byte array.
func (ad AppData) Marshal() ([]byte, error) {
	return json.Marshal(ad.appData)
}

// ID implements model.AppData interface.
func (ad *AppData) ID() string { return ad.appData.ID.Hex() }

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

// InviteTokenLifespan implements model.AppData interface.
func (ad *AppData) InviteTokenLifespan() int64 { return ad.appData.InviteTokenLifespan }

// RefreshTokenLifespan implements model.AppData interface.
func (ad *AppData) RefreshTokenLifespan() int64 { return ad.appData.RefreshTokenLifespan }

// TokenLifespan implements model.AppData interface.
func (ad *AppData) TokenLifespan() int64 { return ad.appData.TokenLifespan }

// TokenPayload implements model.AppData interface.
func (ad *AppData) TokenPayload() []string { return ad.appData.TokenPayload }

// RegistrationForbidden implements model.AppData interface.
func (ad *AppData) RegistrationForbidden() bool { return ad.appData.RegistrationForbidden }

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
