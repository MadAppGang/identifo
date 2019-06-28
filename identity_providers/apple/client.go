package apple

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/madappgang/identifo/model"
)

const (
	// appleAPIPath is a base URL path for connecting to Apple REST API.
	appleAPIPath = "https://appleid.apple.com/"
)

// NewClient creates new HTTP client for communicating with Apple REST API servers.
func NewClient(authorizationCode string, appleInfo *model.AppleInfo) *Client {
	c := &Client{
		AuthorizationCode: authorizationCode,
		ClientID:          appleInfo.ClientID,
		ClientSecret:      appleInfo.ClientSecret,
		HTTPClient:        &http.Client{Timeout: 15 * time.Second},
	}
	c.BaseURL, _ = url.Parse(appleAPIPath)
	return c
}

// Client is a client for making REST API requests to Apple authorization servers.
type Client struct {
	AuthorizationCode string
	ClientID          string
	ClientSecret      string
	BaseURL           *url.URL
	HTTPClient        *http.Client
}

type appleTokenResponse struct {
	IDToken string `json:"id_token"`
}

// User is what we can get about the user from Apple.
type User struct {
	ID string
}

// MyProfile asks for token's owner public profile information.
// Currently, everything Apple provides us with is an obfuscated unique user identifier.
func (c *Client) MyProfile() (User, error) {
	form := url.Values{}

	form.Set("client_id", c.ClientID)
	form.Set("client_secret", c.ClientSecret)
	form.Set("code", c.AuthorizationCode)
	form.Set("grant_type", "authorization_token")

	var user User

	req, err := c.formRequest("POST", "/auth/token", form)
	if err != nil {
		return user, err
	}

	var resp appleTokenResponse
	if _, err = c.do(req, &resp); err != nil {
		return user, err
	}

	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(resp.IDToken, claims, nil); err != nil {
		return user, err
	}

	user.ID = claims["sub"].(string)
	if user.ID == "" {
		return user, fmt.Errorf("ID token has empty subject")
	}
	return user, nil
}

func (c *Client) formRequest(method, path string, form url.Values) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Apple response error: %d", resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(v); err != nil {
		return nil, err
	}
	return resp, nil
}
