package http

//ContextKey enumerates all context keys
type ContextKey int

const (
	//AppDataContextKey context key to keep requested app data
	AppDataContextKey ContextKey = iota
)
