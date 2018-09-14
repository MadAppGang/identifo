package model

//TokenService manage tokens abstraction layer
type TokenService interface {
	NewToken(u User, scopes []string) (Token, error)
	Parse(string) (Token, error)
}

//Token is app token to give user chan
type Token interface {
	Validate() error
	String() string
}
