package facebook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// FacebookAPIPath is a Facebook Graph API base URL path.
	FacebookAPIPath = "https://graph.facebook.com/v3.2/"
)

// NewClient creates new http client with short lived access token.
func NewClient(accessToken string) *Client {
	c := &Client{
		AccessToken: accessToken,
		HTTPClient:  &http.Client{Timeout: 15 * time.Second}, //could be reassigned after
	}
	c.BaseURL, _ = url.Parse(FacebookAPIPath)
	return c
}

// Client is a Facebook SDK client for making GraphAPI requests and handling errors.
type Client struct {
	BaseURL     *url.URL
	AccessToken string
	HTTPClient  *http.Client
}

// MyProfile asks for token's owner public profile information: id and name.
// More functions can be implemented later.
func (c *Client) MyProfile() (User, error) {
	v := url.Values{}
	v.Add("fields", "name")
	v.Add("fields", "id")
	req, err := c.request("GET", "/me", v, nil)
	if err != nil {
		return User{}, err
	}
	var user User
	_, err = c.do(req, &user)
	return user, err
}

func (c *Client) request(method, path string, params url.Values, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	q := params
	if q == nil {
		q = url.Values{}
	}
	q.Add("access_token", c.AccessToken)
	req.URL.RawQuery = q.Encode()
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Facebook response error: %d", resp.StatusCode)
	}
	if err = json.NewDecoder(resp.Body).Decode(v); err != nil {
		return nil, err
	}
	return resp, nil
}
