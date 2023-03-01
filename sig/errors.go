package sig

import "errors"

var (
	ErrorEmptyHost                = errors.New("empty host value")
	ErrorMissingMD5Header         = errors.New("missing 'Content-MD5' header in requests")
	ErrorMissingContentTypeHeader = errors.New("missing 'Content-Type' header in requests")
	ErrorIncorrectMD5Header       = errors.New("incorrect 'Content-MD5' header in requests")
	ErrorMissingDateHeader        = errors.New("missing 'Date' header in requests")
	ErrorMissingExpiresHeader     = errors.New("missing 'Expires' header in requests")
	ErrorMissingDigestHeader      = errors.New("missing 'Digest' header in requests")
	ErrorIncorrectDigestHeader    = errors.New("incorrect 'Digest' header in requests")
	ErrorIncorrectExpireHeader    = errors.New("incorrect 'Expires' header in requests")
	ErrorExpiredRequest           = errors.New("the request is expired")
	ErrorSignatureMismatch        = errors.New("request signature mismatch")
)
