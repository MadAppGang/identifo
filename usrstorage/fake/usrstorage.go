package fake

import (
	"math/rand"

	"github.com/MadAppGang/identifo"
	"github.com/MadAppGang/identifo/usrstorage"
)

const (
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	userIDLength = 10
	//PasswordKey profile key to keep password
	PasswordKey = "password"
)

//fakeClient is a fake user storage client, basically keeping all the users in memory
type fakeClient struct {
	users []identifo.User
}

//NewClient instantiate and setup new fake user storage client
//users: predefined user list
func NewClient(users []identifo.User) usrstorage.Client {
	c := fakeClient{}
	if users != nil {
		c.users = users
	}
	return &c
}

//Connect initiates and returns empty fake user session, which do nothing :-)
func (c *fakeClient) Connect() (usrstorage.Session, error) {
	u := fakeSession{storage: fakeStorage{}}
	u.storage.client = c //assing reference to the client, because it holds the users list
	return &u, nil
}

//userStorageSession is session with no context and nothing, because it's fake
type fakeSession struct {
	storage fakeStorage
}

//Storage returns
func (uss *fakeSession) Storage() usrstorage.Storage {
	return &uss.storage
}

//fakeStorage is service, that could persist the user sotrage (or at least pretend to)
type fakeStorage struct {
	client *fakeClient
}

//FindUser implements Storage protocol function, find and returns the user from memory
func (s *fakeStorage) FindUser(userID identifo.UserID, password string) (*identifo.User, error) {
	if s.client == nil {
		return nil, ErrorStorageNotConfigured
	}

	//it's better to have map instead of array, but who cares
	for _, u := range s.client.users {

		//password should match
		if u.ID == userID && u.Profile[PasswordKey] == password {
			return &u, nil
		}
	}

	return nil, ErrorUserNotFound
}

func (s *fakeStorage) FindUserWithKey(key string, keyValue interface{}, password string) (*identifo.User, error) {
	if s.client == nil {
		return nil, ErrorStorageNotConfigured
	}

	//it's better to have map instead of array, but who cares
	for _, u := range s.client.users {

		//password should match
		if u.Profile[key] == keyValue && u.Profile[PasswordKey] == password {
			return &u, nil
		}
	}

	return nil, ErrorUserNotFound
}

//CreateUser implements CreateUser protocol function, create the new user in memory
func (s *fakeStorage) CreateUser(user identifo.User, password string) (*identifo.User, error) {
	user.ID = newUserID()
	if user.Profile == nil {
		user.Profile = map[string]interface{}{}
	}
	user.Profile[PasswordKey] = password
	if user.Profile["email"] == nil {
		user.Profile["email"] = "fake@fake.com"
	}
	s.client.users = append(s.client.users, user)
	return &user, nil
}

//newUserID generates new userID
func newUserID() identifo.UserID {
	b := make([]byte, userIDLength)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return identifo.UserID(b)
}
