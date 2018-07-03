package mock

import (
	"github.com/MadAppGang/identifo"
	"github.com/MadAppGang/identifo/usrstorage"
)

type Client struct {
	ConnectFn      func() (usrstorage.Session, error)
	ConnectInvoked bool
}

func (c *Client) Connect() (usrstorage.Session, error) {
	c.ConnectInvoked = true
	return c.ConnectFn()
}

type Session struct {
	StorageFn      func() usrstorage.Storage
	StorageInvoked bool
}

func (s *Session) Storage() usrstorage.Storage {
	s.StorageInvoked = true
	return s.StorageFn()
}

type Storage struct {
	FindUserFn      func(userID identifo.UserID, password string) (*identifo.User, error)
	FindUserInvoked bool

	FindUserWithKeyFn      func(key string, keyValue interface{}, password string) (*identifo.User, error)
	FindUserWithKeyInvoked bool

	CreateUserFn      func(user identifo.User, password string) (*identifo.User, error)
	CreateUserInvoked bool
}

func (s *Storage) FindUser(userID identifo.UserID, password string) (*identifo.User, error) {
	s.FindUserInvoked = true
	return s.FindUserFn(userID, password)
}

func (s *Storage) FindUserWithKeyFindUser(key string, keyValue interface{}, password string) (*identifo.User, error) {
	s.FindUserWithKeyInvoked = true
	return s.FindUserWithKeyFn(key, keyValue, password)
}

func (s *Storage) CreateUser(user identifo.User, password string) (*identifo.User, error) {
	s.CreateUserInvoked = true
	return s.CreateUserFn(user, password)
}
