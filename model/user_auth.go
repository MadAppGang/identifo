package model

// login response for user
type AuthResponse struct {
	IDToken         *string `json:"id_token,omitempty"`
	AccessToken     *string `json:"access_token,omitempty"`
	RefreshToken    *string `json:"refresh_token,omitempty"`
	RedirectURI     *string `json:"redirect_uri,omitempty"`
	ClientChallenge *string `json:"client_challenge,omitempty"`
}
