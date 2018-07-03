package fake

//Error is package level error
type Error string

//Error implements standart error protocol
func (e Error) Error() string { return string(e) }

const (
	//ErrorStorageNotConfigured indicated the storage is not configured
	ErrorStorageNotConfigured = Error("storage not configured")

	//ErrorUserNotFound there is not such user in memory
	ErrorUserNotFound = Error("user not found")
)
