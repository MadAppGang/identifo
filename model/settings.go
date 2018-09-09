package model

//JWTSettings holds JWT specific settings
type JWTSettings interface {
	//Algorithm returns default signature algorithm, could be RS256, ES256, none
	Algorithm() string
	//SignatureSecret returns signature secret to sign and verify
	SignatureSecret() []byte
}

// https://github.com/dgrijalva/jwt-go/blob/master/cmd/jwt/app.go
// https://godoc.org/github.com/dgrijalva/jwt-go#MapClaims
