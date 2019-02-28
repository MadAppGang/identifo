package shared

import (
	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
)

// ParseToken parses token and validates it.
func ParseToken(tstr string, ts model.TokenService, tokenType string) (model.Token, error) {
	token, err := ts.Parse(string(tstr))
	if err != nil {
		return nil, err
	}

	v := jwt.NewValidator("identifo", ts.Issuer(), "")
	if err := v.Validate(token); err != nil {
		return nil, err
	}

	if token.Type() != tokenType {
		return nil, err
	}

	return token, nil
}
