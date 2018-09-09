package http

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/urfave/negroni"
)

const (
	//SignatureHeaderKey header to store HMAC signature digest
	SignatureHeaderKey = "Digest"
	//SignatureHeaderValuePrefix signature prefix, indicating hash algorithm, hardcoded now, could be dynamic in the future
	SignatureHeaderValuePrefix = "SHA-256="
)

//SignatureHandler returns middleware, that hanles
//Digest: SHA-256=X48E9qOokqqrvdts8nOJRJN3OWDUoyWxBf7kbu9DBPE=
//https://identifo.madappgang.com/#ca6498ab-b3dc-4c1e-a5b0-2dd633831e2d
func (ar *apiRouter) SignatureHandler(secret string) negroni.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		//read and decode request signature im header
		reqMAC := extractSignature(r.Header.Get(SignatureHeaderKey))
		if reqMAC == nil {
			ar.logger.Println("Error extracting signature")
			ar.Error(rw, ErrorRequestSignature, http.StatusBadRequest, "")
			return
		}

		//extract body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ar.logger.Printf("Error reading body: %v", err)
			ar.Error(rw, ErrorWrongInput, http.StatusBadRequest, "")
			return
		}

		//check body signature
		if err := validateBodySignature(body, reqMAC, []byte(secret)); err != nil {
			ar.logger.Printf("Error validating request signature: %v", err)
			ar.Error(rw, err, http.StatusBadRequest, "")
			return
		}

		//return body as Reader to next handlers
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		// call next handler
		next(rw, r)
	}
}

//extractSignature extracts signature from raw header value and returns it's byte representation
//return nil slice if something wrong happens
func extractSignature(b64 string) []byte {
	b64 = strings.TrimSpace(b64)

	if (len(b64) <= len(SignatureHeaderValuePrefix)) ||
		(strings.ToUpper(b64[0:len(SignatureHeaderValuePrefix)]) != SignatureHeaderValuePrefix) {
		return nil
	}
	//extract Base64 part of signature, trim prefix
	b64 = b64[len(SignatureHeaderValuePrefix):]

	//try to decode to byte, could be not valid Base64
	reqMAC, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil
	}
	return reqMAC
}

//validateBodySignature validates signature for request with `body` is match to signature `reqMAC`, signed with `secret`
func validateBodySignature(body, reqMAC, secret []byte) error {
	mac := hmac.New(sha256.New, secret)
	mac.Write(body)
	expectedMAC := mac.Sum(nil)
	if !hmac.Equal(reqMAC, expectedMAC) {
		return ErrorRequestSignature
	}
	return nil
}
