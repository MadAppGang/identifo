package http

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/madappgang/identifo/model"
)

//NewUserPayloadProvider creates new HTTP webhood  provider
//it basically to the call to 3rd party http service
//to secure this interaction, the receiver should apply some actions to ensure
//the authorized Identity service is doing the request
//To provide that level of verification we are signing the request with HMAC-SHA256 signature
//https://en.wikipedia.org/wiki/HMAC
//We are not using SHA1 here, because SHA2 is more secure.
//We have limited SHA2 with SHA256 simplify the implementation, and SHA256 is the most popular among SHA2 family
//SHA3 is not so popular yet and is limited in client packages available
//Please verify signature on your side
//
//you  can also whitelist identifo's IP as an extra step
func NewUserPayloadProvider(secret string, serviceURL string) (model.UserPayloadProvider, error) {
	if len(secret) < 5 {
		return nil, errors.New("http user payload provider init error, the secret is empty or short, it should be at least 5 chars long")
	}
	_, err := url.Parse(serviceURL)
	if err != nil {
		return nil, fmt.Errorf("http user payload provider init error, bad service URL , %v", err)
	}
	p := provider{
		secret: secret,
		url:    serviceURL,
	}
	return &p, nil
}

type provider struct {
	secret string
	url    string
}

func (p *provider) UserPayloadForApp(appId, appName, userId string) (map[string]interface{}, error) {
	body, _ := json.Marshal(map[string]string{
		"app_id":   appId,
		"app_name": appName,
		"user_id":  userId,
	})
	h := hmac.New(sha256.New, []byte(p.secret))
	h.Write(body)
	sha := hex.EncodeToString(h.Sum(nil))

	client := &http.Client{}
	request, err := http.NewRequest("POST", p.url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("creating http client fro http user payload provider: %v", err)
	}
	request.Header.Set("Digest", "SHA-256="+sha)
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("getting user payload: %v", err)
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("getting user payload from provider, response code expected 200, got: %d", resp.StatusCode)
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("getting user payload from provider, could not read response body with error: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("getting user payload from provider, could not parse response body with error: %v", err)
	}
	return result, nil
}
