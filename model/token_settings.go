package model

type TokenSettings struct {
	// Token settings
	TokenLifespan                     int64                             `bson:"token_lifespan" json:"token_lifespan"`                 // TokenLifespan is a token lifespan in seconds, if 0 - default one is used.
	InviteTokenLifespan               int64                             `bson:"invite_token_lifespan" json:"invite_token_lifespan"`   // InviteTokenLifespan a inviteToken lifespan in seconds, if 0 - default one is used.
	RefreshTokenLifespan              int64                             `bson:"refresh_token_lifespan" json:"refresh_token_lifespan"` // RefreshTokenLifespan is a refreshToken lifespan in seconds, if 0 - default one is used.
	TokenPayload                      []string                          `bson:"token_payload" json:"token_payload"`                   // Payload is a list of fields that are included in token. If it's empty, there are no fields in payload.
	TokenPayloadService               TokenPayloadServiceType           `json:"token_payload_service" bson:"token_payload_service"`
	TokenPayloadServicePluginSettings TokenPayloadServicePluginSettings `json:"token_payload_service_plugin_settings" bson:"token_payload_service_plugin_settings"`
	TokenPayloadServiceHttpSettings   TokenPayloadServiceHttpSettings   `json:"token_payload_service_http_settings" bson:"token_payload_service_http_settings"`
}

// TokenPayloadServiceType service to allow fetch additional data to include to access token
type TokenPayloadServiceType string

const (
	TokenPayloadServiceNone   = "none"   // TokenPayloadServiceNone no service is used
	TokenPayloadServicePlugin = "plugin" // TokenPayloadServicePlugin user local identifo plugin with specific name to retreive token payload
	TokenPayloadServiceHttp   = "http"   // TokenPayloadServiceHttp use external service to get token paylad
)

// TokenPayloadServicePluginSettings settings for token payload service
type TokenPayloadServicePluginSettings struct {
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}

// TokenPayloadServiceHttpSettings settings for token payload service
type TokenPayloadServiceHttpSettings struct {
	URL    string `json:"url,omitempty" bson:"url,omitempty"`
	Secret string `json:"secret,omitempty" bson:"secret,omitempty"`
}
