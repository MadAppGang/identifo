package mongo

import (
	"encoding/json"
	"strings"

	"github.com/madappgang/identifo/model"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// UsersCollection is a collection name for users.
	UsersCollection = "Users"
)

// NewUserStorage creates and inits MongoDB user storage.
func NewUserStorage(db *DB) (model.UserStorage, error) {
	us := &UserStorage{db: db}

	s := us.db.Session(UsersCollection)
	defer s.Close()

	if err := s.C.EnsureIndex(mgo.Index{
		Key: []string{"username"},
		Collation: &mgo.Collation{
			Locale:   "en",
			Strength: 1,
		},
		Sparse: true,
		Unique: true,
	}); err != nil {
		return nil, err
	}

	if err := s.C.EnsureIndex(mgo.Index{
		Key:    []string{"email"},
		Sparse: true,
		Unique: true,
	}); err != nil {
		return nil, err
	}
	if err := s.C.EnsureIndex(mgo.Index{
		Key:    []string{"phone"},
		Sparse: true,
		Unique: true,
	}); err != nil {
		return nil, err
	}

	return us, nil
}

// UserStorage implements user storage interface.
type UserStorage struct {
	db *DB
}

// NewUser returns pointer to newly created user.
func (us *UserStorage) NewUser() model.User {
	return &User{}
}

// UserByID returns user by its ID.
func (us *UserStorage) UserByID(id string) (model.User, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, model.ErrorWrongDataFormat
	}
	s := us.db.Session(UsersCollection)
	defer s.Close()

	var u userData
	if err := s.C.FindId(bson.ObjectIdHex(id)).One(&u); err != nil {
		return nil, err
	}
	return &User{userData: u}, nil
}

// UserByEmail returns user by its email.
func (us *UserStorage) UserByEmail(email string) (model.User, error) {
	if email == "" {
		return nil, model.ErrorWrongDataFormat
	}
	email = strings.ToLower(email)
	s := us.db.Session(UsersCollection)
	defer s.Close()

	var u userData
	if err := s.C.Find(bson.M{"email": email}).One(&u); err != nil {
		return nil, err
	}
	return &User{userData: u}, nil
}

// UserByFederatedID returns user by federated ID.
func (us *UserStorage) UserByFederatedID(provider model.FederatedIdentityProvider, id string) (model.User, error) {
	s := us.db.Session(UsersCollection)
	defer s.Close()
	sid := string(provider) + ":" + id

	var u userData
	if err := s.C.Find(bson.M{"federatedIDs": sid}).One(&u); err != nil {
		return nil, model.ErrorNotFound
	}
	//clear password hash
	u.Pswd = ""
	return &User{userData: u}, nil
}

// UserExists checks if user with provided name exists.
func (us *UserStorage) UserExists(name string) bool {
	s := us.db.Session(UsersCollection)
	defer s.Close()

	strictPattern := "^" + name + "$"
	q := bson.M{"$regex": bson.RegEx{Pattern: strictPattern, Options: "i"}}
	var u userData
	err := s.C.Find(bson.M{"username": q}).One(&u)

	return err == nil
}

//AttachDeviceToken do nothing here
//TODO: implement device storage
func (us *UserStorage) AttachDeviceToken(id, token string) error {
	//we are not supporting devices for users here
	return model.ErrorNotImplemented
}

//DetachDeviceToken do nothing here yet
//TODO: implement
func (us *UserStorage) DetachDeviceToken(token string) error {
	return model.ErrorNotImplemented
}

//RequestScopes for now returns requested scope
//TODO: implement scope logic
func (us *UserStorage) RequestScopes(userID string, scopes []string) ([]string, error) {
	return scopes, nil
}

// Scopes returns supported scopes, could be static data of database.
func (us *UserStorage) Scopes() []string {
	// we allow all scopes for embedded database, you could implement your own logic in external service.
	return []string{"offline", "user"}
}

func (us *UserStorage) UserByPhone(phone string) (model.User, error) {
	s := us.db.Session(UsersCollection)
	defer s.Close()

	var u userData
	if err := s.C.Find(bson.M{"phone": phone}).One(&u); err != nil {
		return nil, model.ErrorNotFound
	}
	u.Pswd = ""

	return &User{userData: u}, nil
}

// UserByNamePassword returns user by name and password.
func (us *UserStorage) UserByNamePassword(name, password string) (model.User, error) {
	s := us.db.Session(UsersCollection)
	defer s.Close()

	var u userData
	strictPattern := "^" + name + "$"
	q := bson.M{"$regex": bson.RegEx{Pattern: strictPattern, Options: "i"}}
	if err := s.C.Find(bson.M{"username": q}).One(&u); err != nil {
		return nil, model.ErrorNotFound
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Pswd), []byte(password)) != nil {
		return nil, model.ErrorNotFound
	}
	//clear password hash
	u.Pswd = ""
	return &User{userData: u}, nil
}

// AddNewUser adds new user to the database.
func (us *UserStorage) AddNewUser(usr model.User, password string) (model.User, error) {
	usr.SetEmail(strings.ToLower(usr.Email()))
	u, ok := usr.(*User)
	if !ok {
		return nil, model.ErrorWrongDataFormat
	}

	s := us.db.Session(UsersCollection)
	defer s.Close()

	u.userData.ID = bson.NewObjectId()
	if len(password) > 0 {
		u.userData.Pswd = PasswordHash(password)
	}

	err := s.C.Insert(u.userData)
	if mgo.IsDup(err) {
		return nil, model.ErrorUserExists
	}

	return u, err
}

// AddUserByPhone registers new user.
func (us *UserStorage) AddUserByPhone(phone string) (model.User, error) {
	u := userData{Active: true, Phone: phone}

	s := us.db.Session(UsersCollection)
	defer s.Close()

	err := s.C.Insert(u)
	if mgo.IsDup(err) {
		return nil, model.ErrorUserExists
	}

	return &User{userData: u}, err
}

// AddUserByNameAndPassword registers new user.
func (us *UserStorage) AddUserByNameAndPassword(username, password string, profile map[string]interface{}) (model.User, error) {
	u := userData{Active: true, Username: username, Profile: profile}
	return us.AddNewUser(&User{userData: u}, password)
}

// AddUserWithFederatedID adds new user with social ID.
func (us *UserStorage) AddUserWithFederatedID(provider model.FederatedIdentityProvider, federatedID string) (model.User, error) {
	// If there is no error, it means user already exists.
	if _, err := us.UserByFederatedID(provider, federatedID); err == nil {
		return nil, model.ErrorUserExists
	}
	sid := string(provider) + ":" + federatedID
	u := userData{Active: true, Username: sid, FederatedIDs: []string{sid}}
	return us.AddNewUser(&User{userData: u}, "")
}

// UpdateUser updates user in MongoDB storage.
func (us *UserStorage) UpdateUser(userID string, newUser model.User) (model.User, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, model.ErrorWrongDataFormat
	}
	newUser.SetEmail(strings.ToLower(newUser.Email()))
	res, ok := newUser.(*User)
	if !ok || res == nil {
		return nil, model.ErrorWrongDataFormat
	}

	// use ID from the request
	res.userData.ID = bson.ObjectIdHex(userID)

	s := us.db.Session(UsersCollection)
	defer s.Close()

	var ud userData
	update := mgo.Change{
		Update:    bson.M{"$set": res.userData},
		ReturnNew: true,
	}
	if _, err := s.C.FindId(bson.ObjectIdHex(userID)).Apply(update, &ud); err != nil {
		return nil, err
	}

	return &User{userData: ud}, nil
}

// ResetPassword sets new user's password.
func (us *UserStorage) ResetPassword(id, password string) error {
	if !bson.IsObjectIdHex(id) {
		return model.ErrorWrongDataFormat
	}
	s := us.db.Session(UsersCollection)
	defer s.Close()

	hash := PasswordHash(password)
	update := bson.M{"$set": bson.M{"pswd": hash}}
	return s.C.UpdateId(bson.ObjectIdHex(id), update)
}

// ResetUsername sets new user's username.
func (us *UserStorage) ResetUsername(id, username string) error {
	if !bson.IsObjectIdHex(id) {
		return model.ErrorWrongDataFormat
	}
	s := us.db.Session(UsersCollection)
	defer s.Close()

	update := bson.M{"$set": bson.M{"username": username}}
	return s.C.UpdateId(bson.ObjectIdHex(id), update)
}

// IDByName returns userID by name.
func (us *UserStorage) IDByName(name string) (string, error) {
	s := us.db.Session(UsersCollection)
	defer s.Close()

	var u userData
	strictPattern := "^" + name + "$"
	q := bson.M{"$regex": bson.RegEx{Pattern: strictPattern, Options: "i"}}
	if err := s.C.Find(bson.M{"username": q}).One(&u); err != nil {
		return "", model.ErrorNotFound
	}

	user := &User{userData: u}

	if !user.Active() {
		return "", ErrorInactiveUser
	}

	return user.ID(), nil
}

// DeleteUser deletes user by id.
func (us *UserStorage) DeleteUser(id string) error {
	if !bson.IsObjectIdHex(id) {
		return model.ErrorWrongDataFormat
	}
	s := us.db.Session(UsersCollection)
	defer s.Close()

	err := s.C.RemoveId(bson.ObjectIdHex(id))
	return err
}

// FetchUsers fetches users which name satisfies provided filterString.
// Supports pagination.
func (us *UserStorage) FetchUsers(filterString string, skip, limit int) ([]model.User, error) {
	s := us.db.Session(UsersCollection)
	defer s.Close()

	q := bson.M{"username": bson.M{"$regex": bson.RegEx{Pattern: filterString, Options: "i"}}}

	orderByField := "username"

	var users []model.User
	err := s.C.Find(q).Sort(orderByField).Limit(limit).Skip(skip).All(&users)
	return users, err
}

// ImportJSON imports data from JSON.
func (us *UserStorage) ImportJSON(data []byte) error {
	ud := []userData{}
	if err := json.Unmarshal(data, &ud); err != nil {
		return err
	}
	for _, u := range ud {
		pswd := u.Pswd
		u.Pswd = ""
		if _, err := us.AddNewUser(&User{userData: u}, pswd); err != nil {
			return err
		}
	}
	return nil
}

// User data implementation.
type userData struct {
	ID           bson.ObjectId          `bson:"_id,omitempty" json:"id,omitempty"`
	Username     string                 `bson:"username,omitempty" json:"username,omitempty"`
	Email        string                 `bson:"email,omitempty" json:"email,omitempty"`
	Phone        string                 `bson:"phone,omitempty" json:"phone,omitempty"`
	Pswd         string                 `bson:"pswd,omitempty" json:"pswd,omitempty"`
	Profile      map[string]interface{} `bson:"profile,omitempty" json:"profile,omitempty"`
	Active       bool                   `bson:"active,omitempty" json:"active,omitempty"`
	FederatedIDs []string               `bson:"federated_ids,omitempty" json:"federated_ids,omitempty"`
}

// User is a data structure for MongoDB storage.
type User struct {
	userData
}

// Sanitize removes sensitive data.
func (u *User) Sanitize() model.User {
	u.userData.Pswd = ""
	return u
}

// UserFromJSON deserializes user from JSON.
func UserFromJSON(d []byte) (*User, error) {
	user := userData{}
	if err := json.Unmarshal(d, &user); err != nil {
		return &User{}, err
	}
	return &User{userData: user}, nil
}

// ID implements model.User interface.
func (u *User) ID() string { return u.userData.ID.Hex() }

// Username implements model.User interface.
func (u *User) Username() string { return u.userData.Username }

// SetUsername implements model.User interface.
func (u *User) SetUsername(username string) { u.userData.Username = username }

// Email implements model.Email interface.
func (u *User) Email() string { return u.userData.Email }

// SetEmail implements model.Email interface.
func (u *User) SetEmail(email string) { u.userData.Email = email }

// PasswordHash implements model.User interface.
func (u *User) PasswordHash() string { return u.userData.Pswd }

// Profile implements model.User interface.
func (u *User) Profile() map[string]interface{} { return u.userData.Profile }

// Active implements model.User interface.
func (u *User) Active() bool { return u.userData.Active }

// PasswordHash creates hash with salt for password.
func PasswordHash(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash)
}
