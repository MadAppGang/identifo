package sig

import (
	"fmt"
	"net/http"
)

var (
	DigestHeaderKey       = http.CanonicalHeaderKey("Digest")
	DigestHeaderSHAPrefix = "sha-256="
	ContentMD5HeaderKey   = http.CanonicalHeaderKey("Content-MD5")
	ExpiresHeaderKey      = http.CanonicalHeaderKey("Expires")
	DateHeaderKey         = http.CanonicalHeaderKey("Date")
	ContentTypeHeaderKey  = http.CanonicalHeaderKey("Content-Type")
	KeyIDHeaderKey        = http.CanonicalHeaderKey("X-Nl-Key-Id")
)

type SigningData struct {
	Method      string
	BodyMD5     string
	ContentType string
	Date        string
	Expires     int64
	Host        string
}

func (sd SigningData) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%d\n%s", sd.Method, sd.BodyMD5, sd.ContentType, sd.Date, sd.Expires, sd.Host)
}
