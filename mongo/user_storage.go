package mongo

import (
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
	return User{userData: u}, nil
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
	return User{userData: u}, nil
}

//AddNewUser adds new user
func (us *UserStorage) AddNewUser(usr model.User, password string) (model.User, error) {
	u, ok := usr.(User)
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

//data implementation
type userData struct {
	ID      bson.ObjectId          `bson:"_id,omitempty"`
	Name    string                 `bson:"name,omitempty"`
	Pswd    string                 `bson:"pswd,omitempty"`
	Profile map[string]interface{} `bson:"profile,omitempty"`
	Active  bool                   `bson:"active,omitempty"`
}

//User user data structure for mongodb storage
type User struct {
	userData
}

//model.User interface implementation
func (u User) ID() string                      { return u.userData.ID.Hex() }
func (u User) Name() string                    { return u.userData.Name }
func (u User) PasswordHash() string            { return u.userData.Pswd }
func (u User) Profile() map[string]interface{} { return u.userData.Profile }
func (u User) Active() bool                    { return u.userData.Active }

//PasswordHash creates hash with salt for password
func PasswordHash(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash)
}
