package mongo

import (
	"encoding/json"

	"gopkg.in/mgo.v2/bson"

	"github.com/madappgang/identifo/model"
	"golang.org/x/crypto/bcrypt"
)

const (
	//UsersCollection collection name for users
	UsersCollection = "Users"
)

//NewUserStorage creates and inits mongodb user storage
func NewUserStorage(db *DB) (model.UserStorage, error) {
	us := UserStorage{}
	us.db = db
	//TODO: ensure indexes
	return &us, nil
}

//UserStorage implements user storage in memory
type UserStorage struct {
	db *DB
}

//UserByID returns user by it's ID
func (us *UserStorage) UserByID(id string) (model.User, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, ErrorWrongDataFormat
	}
	s := us.db.Session(UsersCollection)
	defer s.Close()

	var u userData
	if err := s.C.FindId(bson.ObjectIdHex(id)).One(&u); err != nil {
		return nil, err
	}
	return &User{userData: u}, nil
}

//UserBySocialID returns random generated user
func (us *UserStorage) UserBySocialID(id string) (model.User, error) {
	return nil, ErrorNotImplemented
}

//AttachDeviceToken do nothing here
//TODO: implement device storage
func (us *UserStorage) AttachDeviceToken(id, token string) error {
	//we are not supporting devices for users here
	return nil
}

//RequestScopes for now returns requested scope
//TODO: implement scope logic
func (us *UserStorage) RequestScopes(userID string, scopes []string) ([]string, error) {
	return scopes, nil
}

//UserByNamePassword returns  user by name and password
func (us *UserStorage) UserByNamePassword(name, password string) (model.User, error) {
	s := us.db.Session(UsersCollection)
	defer s.Close()

	var u userData
	q := bson.M{"$regex": bson.RegEx{Pattern: name, Options: "i"}}
	if err := s.C.Find(bson.M{"name": q}).One(&u); err != nil {
		return nil, ErrorNotFound
	}
	if u.Pswd != PasswordHash(password) {
		return nil, ErrorNotFound
	}
	//clear password hash
	u.Pswd = ""
	return &User{userData: u}, nil
}

//AddNewUser adds new user
func (us *UserStorage) AddNewUser(usr model.User, password string) (model.User, error) {
	u, ok := usr.(*User)
	if !ok {
		return nil, ErrorWrongDataFormat
	}
	s := us.db.Session(UsersCollection)
	defer s.Close()
	u.userData.ID = bson.NewObjectId()
	if err := s.C.Insert(u.userData); err != nil {
		return nil, err
	}
	return u, nil
}

//AddUserByNameAndPassword register new user
func (us *UserStorage) AddUserByNameAndPassword(name, password string, profile map[string]interface{}) (model.User, error) {
	//using user name as a key
	_, err := us.UserByID(name)
	//if there is no error, it means user already exists
	if err == nil {
		return nil, ErrorUserExists
	}
	u := userData{}
	u.Active = true
	u.Name = name
	u.Profile = profile
	return us.AddNewUser(&User{u}, password)
}

// UserByName returns user by it's Name
func (us *UserStorage) UserByName(name string) (model.User, error) {
	s := us.db.Session(UsersCollection)
	defer s.Close()

	var u userData
	q := bson.M{"$regex": bson.RegEx{Pattern: name, Options: "i"}}
	if err := s.C.Find(bson.M{"name": q}).One(&u); err != nil {
		return nil, ErrorNotFound
	}

	return &User{userData: u}, nil
}

//data implementation
type userData struct {
	ID      bson.ObjectId          `bson:"_id,omitempty" json:"id,omitempty"`
	Name    string                 `bson:"name,omitempty" json:"name,omitempty"`
	Pswd    string                 `bson:"pswd,omitempty" json:"pswd,omitempty"`
	Profile map[string]interface{} `bson:"profile,omitempty" json:"profile,omitempty"`
	Active  bool                   `bson:"active,omitempty" json:"active,omitempty"`
}

//User user data structure for mongodb storage
type User struct {
	userData
}

//Sanitize removes sensitive data
func (u *User) Sanitize() {
	u.userData.Pswd = ""
	u.userData.Active = false
}

//UserFromJSON deserializes data
func UserFromJSON(d []byte) (*User, error) {
	user := userData{}
	if err := json.Unmarshal(d, &user); err != nil {
		return &User{}, err
	}
	return &User{user}, nil
}

//model.User interface implementation
func (u *User) ID() string                      { return u.userData.ID.Hex() }
func (u *User) Name() string                    { return u.userData.Name }
func (u *User) PasswordHash() string            { return u.userData.Pswd }
func (u *User) Profile() map[string]interface{} { return u.userData.Profile }
func (u *User) Active() bool                    { return u.userData.Active }

//PasswordHash creates hash with salt for password
func PasswordHash(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash)
}
