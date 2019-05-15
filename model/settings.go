package model

// JWTSettings are abstract JWT-settings.
type JWTSettings interface {
	Algorithm() string
	SignatureSecret() []byte
}

// https://github.com/dgrijalva/jwt-go/blob/master/cmd/jwt/app.go
// https://godoc.org/github.com/dgrijalva/jwt-go#MapClaims
