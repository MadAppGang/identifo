package model

const (
	//OfflineScope scope value to request refresh token
	OfflineScope = "offline"
	//RefrestTokenType refresh token type value
	RefrestTokenType = "refresh"
	//AccessTokenType access token type value
	AccessTokenType = "access"
)

//TokenServiceAlgorithm - we support only two now
type TokenServiceAlgorithm int

const (
	//TokenServiceAlgorithmES256 ES256 signature
	TokenServiceAlgorithmES256 TokenServiceAlgorithm = iota
	//TokenServiceAlgorithmRS256 RS256 signature
	TokenServiceAlgorithmRS256
	//TokenServiceAlgorithmAuto try to detect algorithm on the fly
	TokenServiceAlgorithmAuto
)

//TokenService manage tokens abstraction layer
type TokenService interface {
	//NewToken creates new access token for the user
	NewToken(u User, scopes []string, app AppData) (Token, error)
	//NewRefreshToken creates new refresh token for the user
	NewRefreshToken(u User, scopes []string, app AppData) (Token, error)
	//RefreshToken issues the new access token with access token
	RefreshToken(token Token) (Token, error)
	Parse(string) (Token, error)
	String(Token) (string, error)
	Issuer() string
	Algorithm() string
	PublicKey() interface{} //we are not using crypto.PublicKey here to avoid dependencies
	KeyID() string
}

//Token is app token to give user chan
type Token interface {
	Validate() error
	UserID() string
	Type() string
}

//Validator calidate token with external requester
type Validator interface {
	Validate(Token) error
}

//TokenMapping is service to match tokens to services. etc
type TokenMapping interface {
}
