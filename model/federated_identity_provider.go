package model

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
)

var FederatedProviders = map[string]FederatedProvider{
	"facebook": {Name: "Facebook", New: func(params map[string]string, redirectURL string, scopes ...string) (*facebook.Provider, error) {
		return facebook.New(params["ClientId"], params["Secret"], redirectURL, scopes...), nil
	}, Params: []string{"ClientId", "Secret"}},
	"google": {Name: "Google", New: func(params map[string]string, redirectURL string, scopes ...string) (*google.Provider, error) {
		return google.New(params["ClientId"], params["Secret"], redirectURL, scopes...), nil
	}, Params: []string{"ClientId", "Secret"}},
	"apple": {Name: "Apple", New: func(params map[string]string, redirectURL string, scopes ...string) (*apple.Provider, error) {
		jwt.TimeFunc = func() time.Time {
			return time.Now().Add(time.Second * 10)
		}

		secret, err := apple.MakeSecret(apple.SecretParams{
			PKCS8PrivateKey: params["PKCS8PrivateKey"],
			TeamId:          params["TeamId"],
			KeyId:           params["KeyId"],
			ClientId:        params["ClientId"],
			Iat:             int(time.Now().Unix()),
			// Valid for 10 minutes
			Exp: int(time.Now().Unix()) + 10*60,
		})
		if err != nil {
			return nil, err
		}
		return apple.New(params["ClientId"], *secret, redirectURL, nil, scopes...), nil
	}, Params: []string{"ClientId", "PKCS8PrivateKey,textarea", "TeamId", "KeyId"}},
}

type FederatedProvider struct {
	New           interface{} `bson:"-" json:"-"`
	Name          string      `bson:"string,omitempty" json:"string,omitempty"`
	DefaultScopes []string    `bson:"default_scopes,omitempty" json:"default_scopes,omitempty"`
	Params        []string    `bson:"params,omitempty" json:"params,omitempty"`
}

type FederatedProviderSettings struct {
	Params map[string]string `bson:"params,omitempty" json:"params,omitempty"`
	Scopes []string          `bson:"scopes,omitempty" json:"scopes,omitempty"`
}

// Session stores data during the auth process with Google.
type FederatedSession struct {
	ProviderSession string
	CallbackUrl     string
	RedirectUrl     string
	AppId           string
	ProviderName    string
	Scopes          []string
}

// Marshal the session into a string
func (s FederatedSession) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s FederatedSession) String() string {
	return s.Marshal()
}

// UnmarshalSession will unmarshal a JSON string into a session.
func UnmarshalFederatedSession(data string) (*FederatedSession, error) {
	sess := &FederatedSession{}
	err := json.NewDecoder(strings.NewReader(data)).Decode(sess)
	return sess, err
}
