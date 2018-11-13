package jwt

import "strings"

//TokenHeaderKeyPrefix token prefix regarding RFCXXX
const TokenHeaderKeyPrefix = "BEARER "

//ExtractTokenFromBearerHeader extracts token from bearer token header value
func ExtractTokenFromBearerHeader(token string) []byte {
	token = strings.TrimSpace(token)
	if (len(token) <= len(TokenHeaderKeyPrefix)) ||
		(strings.ToUpper(token[0:len(TokenHeaderKeyPrefix)]) != TokenHeaderKeyPrefix) {
		return nil
	}

	token = token[len(TokenHeaderKeyPrefix):]
	return []byte(token)
}
