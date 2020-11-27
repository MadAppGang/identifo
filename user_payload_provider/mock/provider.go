package mock

import (
	"github.com/madappgang/identifo/model"
)

func NewTokenPayloadProvider(payload map[string]interface{}) model.TokenPayloadProvider {
	p := provider{payload: payload}
	return &p
}

type provider struct {
	payload map[string]interface{}
}

func (p *provider) TokenPayloadForApp(appId, appName, userId string) (map[string]interface{}, error) {
	return p.payload, nil
}
