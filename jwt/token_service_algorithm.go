package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	// ErrWrongSignatureAlgorithm is for unsupported signature algorithm.
	ErrWrongSignatureAlgorithm = errors.New("Unsupported signature algorithm")
	// ErrEmptyToken is when token is empty.
	ErrEmptyToken = errors.New("Token is empty")
	// ErrTokenInvalid is when token is invalid.
	ErrTokenInvalid = errors.New("Token is invalid")
)

// StrToTokenServiceAlg maps string token service algorithm names to values.
var StrToTokenServiceAlg = map[string]TokenServiceAlgorithm{
	"es256": TokenServiceAlgorithmES256,
	"rs256": TokenServiceAlgorithmRS256,
	"auto":  TokenServiceAlgorithmAuto}

// TokenServiceAlgorithm is a signing algorithm used by the token service.
// For now, we only support ES256 and RS256.
type TokenServiceAlgorithm int

const (
	// TokenServiceAlgorithmES256 is a ES256 signature.
	TokenServiceAlgorithmES256 TokenServiceAlgorithm = iota + 1
	// TokenServiceAlgorithmRS256 is a RS256 signature.
	TokenServiceAlgorithmRS256
	// TokenServiceAlgorithmAuto tries to detect algorithm on the fly.
	TokenServiceAlgorithmAuto
)

// String implements Stringer.
func (alg TokenServiceAlgorithm) String() string {
	switch alg {
	case TokenServiceAlgorithmES256:
		return "es256"
	case TokenServiceAlgorithmRS256:
		return "rs256"
	case TokenServiceAlgorithmAuto:
		return "auto"
	default:
		return fmt.Sprintf("TokenServiceAlgorithm(%d)", alg)
	}
}

// MarshalJSON implements json.Marshaller.
func (alg TokenServiceAlgorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(alg.String())
}

// UnmarshalJSON implements json.Unmarshaller.
func (alg *TokenServiceAlgorithm) UnmarshalJSON(data []byte) error {
	if alg == nil {
		return nil
	}

	var aux string
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	algorithm, ok := map[string]TokenServiceAlgorithm{
		"es256": TokenServiceAlgorithmES256,
		"rs256": TokenServiceAlgorithmRS256,
		"auto":  TokenServiceAlgorithmAuto}[aux]
	if !ok {
		return fmt.Errorf("Invalid TokenServiceAlgorithm %v", aux)
	}

	*alg = algorithm
	return nil
}

// MarshalYAML implements yaml.Marshaller.
func (alg TokenServiceAlgorithm) MarshalYAML() (interface{}, error) {
	return alg.String(), nil
}

// UnmarshalYAML implements yaml.Unmarshaller.
func (alg *TokenServiceAlgorithm) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if alg == nil {
		return nil
	}

	var aux string
	if err := unmarshal(&aux); err != nil {
		return err
	}

	algorithm, ok := StrToTokenServiceAlg[aux]
	if !ok {
		return fmt.Errorf("Invalid TokenServiceAlgorithm %v", aux)
	}

	*alg = algorithm
	return nil
}
