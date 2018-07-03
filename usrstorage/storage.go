package usrstorage

import "github.com/MadAppGang/identifo"

//Client is a client for user storage service
type Client interface {
	Connect() (Session, error)
}

//Session is a Client's session created with context
type Session interface {
	Storage() Storage
}

//Storage is service, that could persist the user sotrage (or at least pretend to)
type Storage interface {
	FindUser(userID identifo.UserID, password string) (*identifo.User, error)
	FindUserWithKey(key string, keyValue interface{}, password string) (*identifo.User, error)
	CreateUser(user identifo.User, password string) (*identifo.User, error)
}
