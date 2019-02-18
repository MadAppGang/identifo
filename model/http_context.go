package model

//ContextKey enumerates all context keys
type ContextKey int

const (
	//AppDataContextKey context key to keep requested app data
	AppDataContextKey ContextKey = iota + 1
	//TokenContextKey bearer token context key
	TokenContextKey
	//TokenRawContextKey bearer token context key in raw format
	TokenRawContextKey
	// AppIDError context key to keep error from app_middleware
	AppIDError
)
