package facebook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	//FacebookAPIPath is facebook Graph API base URL path
	FacebookAPIPath = "https://graph.facebook.com/v3.2/"
)

//NewClient create new default client with short lived access token
func NewClient(accessToken string) *Client {
	c := Client{}
	c.BaseURL, _ = url.Parse(FacebookAPIPath)
	c.AccessToken = accessToken
	c.HTTPClient = http.DefaultClient //could be reassigned after
	return &c
}

//Client is Facebook SDK client to make GraphAPI requests and handle errors
type Client struct {
	BaseURL     *url.URL
	AccessToken string
	HTTPClient  *http.Client
}

//MyProfile asks for token owner public profile information: id and name
//more functions could be implemented after
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
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
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
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
