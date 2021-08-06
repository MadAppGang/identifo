package mongo

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/madappgang/identifo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"golang.org/x/crypto/bcrypt"
)

const usersCollectionName = "Users"

// NewUserStorage creates and inits MongoDB user storage.
func NewUserStorage(settings model.MongodDatabaseSettings) (model.UserStorage, error) {
	if len(settings.ConnectionString) == 0 || len(settings.DatabaseName) == 0 {
		return nil, ErrorEmptyConnectionStringDatabase
	}

	// create database
	db, err := NewDB(settings.ConnectionString, settings.DatabaseName)
	if err != nil {
		return nil, err
	}

	coll := db.Database.Collection(usersCollectionName)
	us := &UserStorage{coll: coll, timeout: 30 * time.Second}

	userNameIndexOptions := &options.IndexOptions{}
	userNameIndexOptions.SetUnique(true)
	userNameIndexOptions.SetSparse(true)
	userNameIndexOptions.SetCollation(&options.Collation{Locale: "en", Strength: 1})

	userNameIndex := &mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "username", Value: bsonx.Int32(int32(1))}},
		Options: userNameIndexOptions,
	}

	emailIndexOptions := &options.IndexOptions{}
	emailIndexOptions.SetUnique(true)
	emailIndexOptions.SetSparse(true)

	emailIndex := &mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "email", Value: bsonx.Int32(int32(1))}},
		Options: emailIndexOptions,
	}

	phoneIndexOptions := &options.IndexOptions{}
	phoneIndexOptions.SetUnique(true)
	phoneIndexOptions.SetSparse(true)

	phoneIndex := &mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "phone", Value: bsonx.Int32(int32(1))}},
		Options: phoneIndexOptions,
	}

	err = db.EnsureCollectionIndices(usersCollectionName, []mongo.IndexModel{*userNameIndex, *emailIndex, *phoneIndex})
	return us, err
}

// UserStorage implements user storage interface.
type UserStorage struct {
	coll    *mongo.Collection
	timeout time.Duration
}

// UserByID returns user by its ID.
func (us *UserStorage) UserByID(id string) (model.User, error) {
	hexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.User{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	var u model.User
	if err := us.coll.FindOne(ctx, bson.M{"_id": hexID.Hex()}).Decode(&u); err != nil {
		return model.User{}, err
	}
	return u, nil
}

// UserByEmail returns user by their email.
func (us *UserStorage) UserByEmail(email string) (model.User, error) {
	if email == "" {
		return model.User{}, model.ErrorWrongDataFormat
	}
	email = strings.ToLower(email)

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	var u model.User
	if err := us.coll.FindOne(ctx, bson.M{"email": email}).Decode(&u); err != nil {
		return model.User{}, err
	}
	// clear password hash
	u.Pswd = ""
	return u, nil
}

// UserByFederatedID returns user by federated ID.
func (us *UserStorage) UserByFederatedID(provider string, id string) (model.User, error) {
	sid := string(provider) + ":" + id

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	var u model.User
	if err := us.coll.FindOne(ctx, bson.M{"federated_ids": sid}).Decode(&u); err != nil {
		return model.User{}, model.ErrUserNotFound
	}
	// clear password hash
	u.Pswd = ""
	return u, nil
}

// UserExists checks if user with provided name exists.
func (us *UserStorage) UserExists(name string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	strictPattern := "^" + name + "$"
	q := bson.D{primitive.E{Key: "username", Value: primitive.Regex{Pattern: strictPattern, Options: "i"}}}

	var u model.User
	err := us.coll.FindOne(ctx, q).Decode(&u)
	return err == nil
}

// AttachDeviceToken do nothing here
// TODO: implement device storage
func (us *UserStorage) AttachDeviceToken(id, token string) error {
	// we are not supporting devices for users here
	return model.ErrorNotImplemented
}

// DetachDeviceToken do nothing here yet
// TODO: implement
func (us *UserStorage) DetachDeviceToken(token string) error {
	return model.ErrorNotImplemented
}

// RequestScopes for now returns requested scope
// TODO: implement scope logic
func (us *UserStorage) RequestScopes(userID string, scopes []string) ([]string, error) {
	return scopes, nil
}

// Scopes returns supported scopes, could be static data of database.
func (us *UserStorage) Scopes() []string {
	// we allow all scopes for embedded database, you could implement your own logic in external service.
	return []string{"offline", "user"}
}

// UserByPhone fetches user by phone number.
func (us *UserStorage) UserByPhone(phone string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	var u model.User
	if err := us.coll.FindOne(ctx, bson.M{"phone": phone}).Decode(&u); err != nil {
		return model.User{}, err
	}
	u.Pswd = ""
	return u, nil
}

// UserByUsername returns user by name.
func (us *UserStorage) UserByUsername(username string) (model.User, error) {
	strictPattern := "^" + strings.ReplaceAll(username, "+", "\\+") + "$"
	q := bson.D{primitive.E{Key: "username", Value: primitive.Regex{Pattern: strictPattern, Options: "i"}}}

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	var u model.User
	if err := us.coll.FindOne(ctx, q).Decode(&u); err != nil {
		return model.User{}, model.ErrUserNotFound
	}

	// clear password hash
	u.Pswd = ""
	return u, nil
}

// AddNewUser adds new user to the database.
func (us *UserStorage) AddNewUser(user model.User, password string) (model.User, error) {
	user.Email = strings.ToLower(user.Email)

	user.ID = primitive.NewObjectID().Hex()
	if len(password) > 0 {
		user.Pswd = model.PasswordHash(password)
	}
	user.NumOfLogins = 0

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	if _, err := us.coll.InsertOne(ctx, user); err != nil {
		if isErrDuplication(err) {
			return model.User{}, model.ErrorUserExists
		}
		return model.User{}, err
	}
	return user, nil
}

// AddUserWithPassword creates new user and saves it in the database.
func (us *UserStorage) AddUserWithPassword(user model.User, password, role string, isAnonymous bool) (model.User, error) {
	if _, err := us.UserByUsername(user.Username); err == nil {
		return model.User{}, model.ErrorUserExists
	}
	if _, err := us.UserByEmail(user.Email); err == nil {
		return model.User{}, model.ErrorUserExists
	}
	if _, err := us.UserByPhone(user.Phone); err == nil {
		return model.User{}, model.ErrorUserExists
	}

	u := model.User{
		ID:         primitive.NewObjectID().Hex(),
		Active:     true,
		Username:   user.Username,
		Phone:      user.Phone,
		Email:      user.Email,
		AccessRole: role,
		Anonymous:  isAnonymous,
	}

	return us.AddNewUser(u, password)
}

// AddUserWithFederatedID adds new user with social ID.
func (us *UserStorage) AddUserWithFederatedID(user model.User, provider string, federatedID, role string) (model.User, error) {
	// If there is no error, it means user already exists.
	if _, err := us.UserByFederatedID(provider, federatedID); err == nil {
		return model.User{}, model.ErrorUserExists
	}

	user.ID = primitive.NewObjectID().Hex()
	user.Active = true
	user.AccessRole = role
	user.AddFederatedId(provider, federatedID)

	return us.AddNewUser(user, "")
}

// UpdateUser updates user in MongoDB storage.
func (us *UserStorage) UpdateUser(userID string, newUser model.User) (model.User, error) {
	hexID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return model.User{}, err
	}

	newUser.Email = strings.ToLower(newUser.Email)
	// use ID from the request
	newUser.ID = userID

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	update := bson.M{"$set": newUser}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var ud model.User
	if err := us.coll.FindOneAndUpdate(ctx, bson.M{"_id": hexID.Hex()}, update, opts).Decode(&ud); err != nil {
		return model.User{}, err
	}
	return ud, nil
}

// ResetPassword sets new user's password.
func (us *UserStorage) ResetPassword(id, password string) error {
	hexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": bson.M{"pswd": model.PasswordHash(password)}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	var ud model.User
	err = us.coll.FindOneAndUpdate(ctx, bson.M{"_id": hexID.Hex()}, update, opts).Decode(&ud)
	return err
}

// CheckPassword check that password is valid for user id.
func (us *UserStorage) CheckPassword(id, password string) error {
	user, err := us.UserByID(id)
	if err != nil {
		return model.ErrUserNotFound
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Pswd), []byte(password)); err != nil {
		// return this error to hide the existence of the user.
		return model.ErrUserNotFound
	}
	return nil
}

// ResetUsername sets new user's username.
func (us *UserStorage) ResetUsername(id, username string) error {
	hexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	update := bson.M{"$set": bson.M{"username": username}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var ud model.User
	err = us.coll.FindOneAndUpdate(ctx, bson.M{"_id": hexID.Hex()}, update, opts).Decode(&ud)
	return err
}

// IDByName returns userID by name.
func (us *UserStorage) IDByName(name string) (string, error) {
	strictPattern := "^" + name + "$"
	q := bson.D{primitive.E{Key: "username", Value: primitive.Regex{Pattern: strictPattern, Options: "i"}}}

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	var u model.User
	if err := us.coll.FindOne(ctx, q).Decode(&u); err != nil {
		return "", model.ErrorNotFound
	}

	if !u.Active {
		return "", ErrorInactiveUser
	}
	return u.ID, nil
}

// DeleteUser deletes user by id.
func (us *UserStorage) DeleteUser(id string) error {
	hexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	_, err = us.coll.DeleteOne(ctx, bson.M{"_id": hexID.Hex()})
	return err
}

// FetchUsers fetches users which name satisfies provided filterString.
// Supports pagination.
func (us *UserStorage) FetchUsers(filterString string, skip, limit int) ([]model.User, int, error) {
	q := bson.D{primitive.E{Key: "username", Value: primitive.Regex{Pattern: filterString, Options: "i"}}}

	ctx, cancel := context.WithTimeout(context.Background(), 2*us.timeout)
	defer cancel()

	total, err := us.coll.CountDocuments(ctx, q)
	if err != nil {
		return []model.User{}, 0, err
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{primitive.E{Key: "username", Value: 1}})
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))

	curr, err := us.coll.Find(ctx, q, findOptions)
	if err != nil {
		return []model.User{}, 0, err
	}

	usersData := []model.User{}
	if err = curr.All(ctx, &usersData); err != nil {
		return []model.User{}, 0, err
	}

	return usersData, int(total), err
}

// ImportJSON imports data from JSON.
func (us *UserStorage) ImportJSON(data []byte) error {
	ud := []model.User{}
	if err := json.Unmarshal(data, &ud); err != nil {
		return err
	}
	for _, u := range ud {
		pswd := u.Pswd
		u.Pswd = ""
		if _, err := us.AddNewUser(u, pswd); err != nil {
			return err
		}
	}
	return nil
}

// UpdateLoginMetadata updates user's login metadata.
func (us *UserStorage) UpdateLoginMetadata(userID string) {
	hexID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Cannot update login metadata of user %s: %s\n", userID, err)
		return
	}

	update := bson.M{
		"$set": bson.M{"latest_login_time": time.Now().Unix()},
		"$inc": bson.M{"num_of_logins": 1},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*us.timeout)
	defer cancel()

	var ud model.User
	if err := us.coll.FindOneAndUpdate(ctx, bson.M{"_id": hexID.Hex()}, update).Decode(&ud); err != nil {
		log.Printf("Cannot update login metadata of user %s: %s\n", userID, err)
	}
}

// Close is a no-op.
func (us *UserStorage) Close() {}
