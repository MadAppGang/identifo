package sig

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const defaultExpireInSec = 30 * time.Second

// VerifySignature Check the signature from request, if signature is valid, not error returns
//
// Signature = URL-Encode( Base64( HMAC-SHA1( YourSecretAccessKey, UTF-8-Encoding-Of( StringToSign ) ) ) );
// StringToSign = HTTP-VERB + "\n" +
//
//	Content-MD5 + "\n" +
//	Content-Type + "\n" +
//	Date + "\n" +
//	Expires+ "\n" +
//	HTTP-HOST
//
//	func SignRequest(r *http.Request, secret, body []byte) (*http.Request, error) {
//		if len(r.Host) == 0 {
//			return r, ErrorEmptyHost
//		}
//	}
func VerifySignature(r *http.Request, secret []byte) error {
	dh := r.Header[DigestHeaderKey]
	hashB := []byte{}
	if len(dh) > 0 {
		digest := dh[0]
		if !strings.HasPrefix(digest, DigestHeaderSHAPrefix) {
			return ErrorIncorrectDigestHeader
		}
		hash := strings.TrimPrefix(digest, DigestHeaderSHAPrefix)
		var err error
		hash, err = url.QueryUnescape(hash)
		if err != nil {
			return ErrorIncorrectDigestHeader
		}
		hashB, err = base64.StdEncoding.DecodeString(hash)
		if err != nil {
			return ErrorIncorrectDigestHeader
		}
	} else if len(dh) > 0 {
		return ErrorMissingDigestHeader
	}

	stringToSign, err := stringToSignFromRequest(r, "")
	if err != nil {
		return err
	}
	fmt.Println(stringToSign.String())
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(stringToSign.String()))

	equal := hmac.Equal(hashB, mac.Sum(nil))
	if equal == false {
		return ErrorSignatureMismatch
	}
	return nil
}

func AddHeadersAndSignRequest(r *http.Request, secret []byte, bodyMD5 string) error {
	if len(bodyMD5) == 0 {
		bodyMD5 = GetBodyMD5(r)
	}
	if len(bodyMD5) > 0 {
		r.Header[ContentMD5HeaderKey] = []string{bodyMD5}
	}
	r.Header[ExpiresHeaderKey] = []string{fmt.Sprintf("%d", time.Now().Add(defaultExpireInSec).Unix())}
	r.Header[DateHeaderKey] = []string{time.Now().Format(time.RFC3339)}
	return SignRequest(r, secret, bodyMD5)
}

func SignRequest(r *http.Request, secret []byte, bodyMD5 string) error {
	stringToSign, err := stringToSignFromRequest(r, bodyMD5)
	if err != nil {
		return err
	}
	signature := SignString(stringToSign.String(), secret)
	r.Header[DigestHeaderKey] = []string{fmt.Sprintf("%s%s", DigestHeaderSHAPrefix, signature)}
	return nil
}

func stringToSignFromRequest(r *http.Request, bodyMD5 string) (SigningData, error) {
	data := SigningData{}
	data.Method = r.Method
	bmd5 := bodyMD5
	if len(bodyMD5) == 0 {
		bmd5 = GetBodyMD5(r)
	}
	md5 := r.Header[ContentMD5HeaderKey]
	if len(md5) > 0 {
		data.BodyMD5 = md5[0]
		if bmd5 != md5[0] {
			return data, ErrorIncorrectMD5Header
		}
	} else if len(bmd5) > 0 {
		return data, ErrorMissingMD5Header
	}

	ct := r.Header[ContentTypeHeaderKey]
	if len(ct) > 0 {
		data.ContentType = ct[0]
	} else {
		return data, ErrorMissingContentTypeHeader
	}

	eh := r.Header[ExpiresHeaderKey]
	if len(eh) > 0 {
		exp, err := strconv.ParseInt(eh[0], 10, 0)
		if err != nil {
			return data, ErrorIncorrectExpireHeader
		}
		if time.Now().After(time.Unix(exp, 0)) {
			return data, ErrorExpiredRequest
		}

		data.Expires = exp
	} else {
		return data, ErrorMissingExpiresHeader
	}

	dh := r.Header[DateHeaderKey]
	if len(dh) > 0 {
		data.Date = dh[0]
	} else {
		return data, ErrorMissingDateHeader
	}

	data.Host = r.Host
	return data, nil
}

func GetBodyMD5(r *http.Request) string {
	// Read the Body content
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body)
	}

	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	if len(bodyBytes) == 0 {
		return ""
	}

	return GetMD5(bodyBytes)
}
