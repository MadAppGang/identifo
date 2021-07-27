package model

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

// StrToTokenSignAlg maps string token service algorithm names to values.
var StrToTokenSignAlg = map[string]TokenSignatureAlgorithm{
	"es256": TokenSignatureAlgorithmES256,
	"rs256": TokenSignatureAlgorithmRS256,
	"auto":  TokenSignatureAlgorithmAuto,
}

// TokenSignatureAlgorithm is a signing algorithm used by the token service.
// For now, we only support ES256 and RS256.
type TokenSignatureAlgorithm int

const (
	// TokenSignatureAlgorithmES256 is a ES256 signature.
	TokenSignatureAlgorithmES256 TokenSignatureAlgorithm = iota + 1
	// TokenSignatureAlgorithmRS256 is a RS256 signature.
	TokenSignatureAlgorithmRS256
	// TokenSignatureAlgorithmAuto tries to detect algorithm on the fly.
	TokenSignatureAlgorithmAuto
)

// String implements Stringer.
func (alg TokenSignatureAlgorithm) String() string {
	switch alg {
	case TokenSignatureAlgorithmES256:
		return "es256"
	case TokenSignatureAlgorithmRS256:
		return "rs256"
	case TokenSignatureAlgorithmAuto:
		return "auto"
	default:
		return fmt.Sprintf("TokenSignatureAlgorithm(%d)", alg)
	}
}

// MarshalJSON implements json.Marshaller.
func (alg TokenSignatureAlgorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(alg.String())
}

// UnmarshalJSON implements json.Unmarshaller.
func (alg *TokenSignatureAlgorithm) UnmarshalJSON(data []byte) error {
	if alg == nil {
		return nil
	}

	var aux string
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	algorithm, ok := StrToTokenSignAlg[aux]
	if !ok {
		return fmt.Errorf("Invalid TokenSignatureAlgorithm %v", aux)
	}

	*alg = algorithm
	return nil
}

// MarshalYAML implements yaml.Marshaller.
func (alg TokenSignatureAlgorithm) MarshalYAML() (interface{}, error) {
	return alg.String(), nil
}

// UnmarshalYAML implements yaml.Unmarshaller.
func (alg *TokenSignatureAlgorithm) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if alg == nil {
		return nil
	}

	var aux string
	if err := unmarshal(&aux); err != nil {
		return err
	}

	algorithm, ok := StrToTokenSignAlg[aux]
	if !ok {
		return fmt.Errorf("Invalid TokenSignatureAlgorithm %v", aux)
	}

	*alg = algorithm
	return nil
}
