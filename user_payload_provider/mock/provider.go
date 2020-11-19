package mock

import (
	"github.com/madappgang/identifo/model"
)

func NewUserPayloadProvider(payload map[string]interface{}) model.UserPayloadProvider {
	p := provider{payload: payload}
	return &p
}

type provider struct {
	payload map[string]interface{}
}

func (p *provider) UserPayloadForApp(appId, appName, userId string) (map[string]interface{}, error) {
	return p.payload, nil
}
