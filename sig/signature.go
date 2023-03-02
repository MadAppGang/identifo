package sig

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/url"
)

// SignString sing the string -> base64 -> url encode
func SignString(s string, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(s))
	str := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	str = url.QueryEscape(str)
	return str
}

// SignString sing the string -> base64 -> url encode
func Sign(s string, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(s))
	str := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	str = url.QueryEscape(str)
	return str
}

func GetMD5(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}
