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

//UserID returns user ID
func (t *Token) UserID() string {
	claims, ok := t.JWT.Claims.(*Claims)
	if !ok {
		return ""
	}
	return claims.Subject
}

//Type returns token type, could be empty or "refresh" only
func (t *Token) Type() string {
	claims, ok := t.JWT.Claims.(*Claims)
	if !ok {
		return ""
	}
	return claims.Type
}

//Claims extended claims structure
type Claims struct {
	UserProfile string `json:"user_profile,omitempty"`
	Scopes      string `json:"scopes,omitempty"`
	Type        string `json:"type,omitempty"` //could be empty, "access" or "refresh" or "reset-password" only
	KeyID       string `json:"kid,omitempty"`  //optional keyID
	jwt.StandardClaims
}

//how to use JWT tokens full example
//https://github.com/dgrijalva/jwt-go/blob/master/cmd/jwt/app.go
