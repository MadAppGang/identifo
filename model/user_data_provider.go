package model

//TokenPayloadProvider provides additional user payload to include to the token
type TokenPayloadProvider interface {
	TokenPayloadForApp(appId, appName, userId string) (map[string]interface{}, error)
}
