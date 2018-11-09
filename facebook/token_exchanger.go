package facebook

import (
	"net/http"
	"net/url"
)

//ExchangeTokenURL is facebook endpoint to exchange shot-lived token
const ExchangeTokenURL = "https://graph.facebook.com/oauth/access_token"

//ExchangeToken exchange short living token to long living token
//https://developers.facebook.com/docs/facebook-login/access-tokens/refreshing
func ExchangeToken(appID, appSecret, shortToken string) (string, error) {

	req, err := http.NewRequest("GET", ExchangeTokenURL, nil)
	if err != nil {
		return "", err
	}
	q := url.Values{}
	q.Add("grant_type", "fb_exchange_token")
	q.Add("client_id", appID)
	q.Add("client_secret", appSecret)
	q.Add("fb_exchange_token", shortToken)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Accept", "application/json")
	return "", nil
}
