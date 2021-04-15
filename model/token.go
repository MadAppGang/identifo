package model

const (
	TokenTypeInvite      = "invite"      // TokenTypeInvite is an invite token type value.
	TokenTypeReset       = "reset"       // TokenTypeReset is an reset token type value.
	TokenTypeWebCookie   = "web-cookie"  // TokenTypeWebCookie is a web-cookie token type value.
	TokenTypeAccess      = "access"      // TokenTypeAccess is an access token type.
	TokenTypeRefresh     = "refresh"     // TokenTypeRefresh is a refresh token type.
	TokenTypeTFAPreauth  = "2fa-preauth" // TokenTypeTFAPreauth is an 2fa preauth token type.
	TokenTFAPreauthScope = "2fa"         // TokenTFAPreauthScope preauth token scope for first step of TFA
)
