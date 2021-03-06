package model

const (
	// TokenTypeInvite is an invite token type value.
	TokenTypeInvite = "invite"
	// TokenTypeReset is an reset token type value.
	TokenTypeReset = "reset"
	// TokenTypeWebCookie is a web-cookie token type value.
	TokenTypeWebCookie = "web-cookie"
	// TokenTypeAccess is an access token type.
	TokenTypeAccess = "access"
	// TokenTypeRefresh is a refresh token type.
	TokenTypeRefresh = "refresh"
	// TokenTypeTFAPreauth is an 2fa preauth token type.
	TokenTypeTFAPreauth = "2fa-preauth"
	// TokenTFAPreauthScope preauth token scope for first step of TFA
	TokenTFAPreauthScope = "2fa"
)
