package model

//UserPayloadProvider provides additional user payload to include to the token
type UserPayloadProvider interface {
	UserPayloadForApp(appId, appName, userId string) (map[string]interface{}, error)
}
