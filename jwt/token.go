package jwt

import jwt "github.com/dgrijalva/jwt-go"

//Token represents JWT token in the system
type Token struct {
	JWT *jwt.Token
	new bool
}

//Validate validates token data, returns nil if all data is valid
func (t *Token) Validate() error {
	if t.JWT == nil {
		return ErrEmptyToken
	}
	if !t.new && !t.JWT.Valid {
		return ErrTokenInvalid
	}

	return nil
}

//Claims extended claims structure
type Claims struct {
	UserProfile string `json:"user_profile,omitempty"`
	Scopes      string `json:"scopes,omitempty"`
	Type        string `json:"type,omitempty"` //could be empty or "refresh" only
	jwt.StandardClaims
}

//how to use JWT tokens full example
//https://github.com/dgrijalva/jwt-go/blob/master/cmd/jwt/app.go
