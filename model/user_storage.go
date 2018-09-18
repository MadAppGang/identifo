package model

//UserStorage introduces user storage service
type UserStorage interface {
	UserByID(id string) (User, error)
	UserBySocialID(id string) (User, error)
	AttachDeviceToken(id, token string) error
	UserByNamePassword(name, password string) (User, error)
	RequestScopes(userID string, scopes []string) ([]string, error)
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
}
