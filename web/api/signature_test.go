package api_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func Signature(data, secret string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))

	if _, err := mac.Write([]byte(data)); err != nil {
		return "", fmt.Errorf("error creating signature for data: %v", err)
	}

	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}
