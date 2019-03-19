package model

//UserStorage introduces user storage service
type UserStorage interface {
	UserByID(id string) (User, error)
	IDByName(name string) (string, error)
	AttachDeviceToken(id, token string) error
	DetachDeviceToken(token string) error
	UserByNamePassword(name, password string) (User, error)
	AddUserByNameAndPassword(name, password string, profile map[string]interface{}) (User, error)
	UserExists(name string) bool
	UserByFederatedID(provider FederatedIdentityProvider, id string) (User, error)
	AddUserWithFederatedID(provider FederatedIdentityProvider, id string) (User, error)
	UpdateUser(userID string, newUser User) (User, error)
	ResetPassword(id, password string) error
	DeleteUser(id string) error
	FetchUsers(search string, skip, limit int) ([]User, error)
	NewUser() User

	RequestScopes(userID string, scopes []string) ([]string, error)
	Scopes() []string
	ImportJSON(data []byte) error
}

//User is abstract representation of the user in auth layer
//everything could be user
//we are not locked on any implementation
type User interface {
	ID() string
	Name() string
	PasswordHash() string
	Profile() map[string]interface{}
	Active() bool
	Sanitize() User
}
