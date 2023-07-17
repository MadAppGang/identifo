package model

// login response for user
type AuthResponse struct {
	IDToken         *string `json:"id_token,omitempty"`
	AccessToken     *string `json:"access_token,omitempty"`
	RefreshToken    *string `json:"refresh_token,omitempty"`
	RedirectURI     *string `json:"redirect_uri,omitempty"`
	ClientChallenge *string `json:"client_challenge,omitempty"`
}

// LoginRequest login request for user for first factor login with secondary id: email, phone, username
type LoginRequest struct {
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Device   string   `json:"device"`
	Scopes   []string `json:"scopes"`
	OS       string   `json:"os"`
}

// LoginStrategy returns login strategy for LoginRequest data
func (l LoginRequest) Strategy() (FirstFactorInternalStrategy, string) {
	idType := AuthIdentityTypePhone
	idValue := l.Phone
	transport := AuthTransportTypeNone

	if len(l.Email) > 0 {
		idType = AuthIdentityTypeEmail
		idValue = l.Email
	} else if len(l.Username) > 0 {
		idType = AuthIdentityTypeUsername
		idValue = l.Username
	}

	return FirstFactorInternalStrategy{
		Transport: transport,
		Identity:  idType,
		Challenge: AuthChallengeTypePassword,
	}, idValue
}
