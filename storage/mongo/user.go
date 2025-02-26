package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const usersCollectionName = "Users"

// NewUserStorage creates and inits MongoDB user storage.
func NewUserStorage(
	logger *slog.Logger,
	settings model.MongoDatabaseSettings) (model.UserStorage, error) {
	if len(settings.ConnectionString) == 0 || len(settings.DatabaseName) == 0 {
		return nil, ErrorEmptyConnectionStringDatabase
	}

	// create database
	db, err := NewDB(logger, settings.ConnectionString, settings.DatabaseName)
	if err != nil {
		return nil, err
	}

	coll := db.database.Collection(usersCollectionName)
	us := &UserStorage{
		logger:  logger,
		coll:    coll,
		timeout: 30 * time.Second,
	}

	userNameIndexOptions := &options.IndexOptions{}
	userNameIndexOptions.SetUnique(false)
	userNameIndexOptions.SetSparse(true)
	userNameIndexOptions.SetCollation(&options.Collation{Locale: "en", Strength: 1})

	userNameIndex := &mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: userNameIndexOptions,
	}

	emailIndexOptions := &options.IndexOptions{}
	emailIndexOptions.SetUnique(false)
	emailIndexOptions.SetSparse(true)

	emailIndex := &mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: emailIndexOptions,
	}

	phoneIndexOptions := &options.IndexOptions{}
	phoneIndexOptions.SetUnique(false)
	phoneIndexOptions.SetSparse(true)

	phoneIndex := &mongo.IndexModel{
		Keys:    bson.D{{Key: "phone", Value: 1}},
		Options: phoneIndexOptions,
	}

	err = db.EnsureCollectionIndices(usersCollectionName, []mongo.IndexModel{*userNameIndex, *emailIndex, *phoneIndex})
	return us, err
}

// UserStorage implements user storage interface.
type UserStorage struct {
	logger  *slog.Logger
	coll    *mongo.Collection
	timeout time.Duration
}

// UserByID returns user by its ID.
func (us *UserStorage) UserByID(id string) (model.User, error) {
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.User{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	var u model.User
	if err := us.coll.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&u); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, model.ErrUserNotFound
		}

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
	if err := us.coll.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&u); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, model.ErrUserNotFound
		}

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
	if err := us.coll.FindOne(ctx, bson.D{{Key: "federated_ids", Value: sid}}).Decode(&u); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, model.ErrUserNotFound
		}

		return model.User{}, err
	}
	// clear password hash
	u.Pswd = ""
	return u, nil
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

// TODO: implement get all device tokens logic
func (us *UserStorage) AllDeviceTokens(userID string) ([]string, error) {
	return nil, nil
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
	if err := us.coll.FindOne(ctx, bson.D{{Key: "phone", Value: phone}}).Decode(&u); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, model.ErrUserNotFound
		}

		return model.User{}, err
	}
	u.Pswd = ""
	return u, nil
}

// UserByUsername returns user by name.
func (us *UserStorage) UserByUsername(username string) (model.User, error) {
	strictPattern := "^" + strings.ReplaceAll(username, "+", "\\+") + "$"
	q := bson.D{{Key: "username", Value: primitive.Regex{Pattern: strictPattern, Options: "i"}}}

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	var u model.User
	if err := us.coll.FindOne(ctx, q).Decode(&u); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, model.ErrUserNotFound
		}

		return model.User{}, err
	}

	// clear password hash
	u.Pswd = ""
	return u, nil
}

// AddNewUser adds new user to the database.
func (us *UserStorage) AddNewUser(user model.User, password string) (model.User, error) {
	user.Email = strings.ToLower(user.Email)

	if len(user.ID) == 0 {
		user.ID = primitive.NewObjectID().Hex()
	}
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
	if len(user.Username) > 0 {
		if _, err := us.UserByUsername(user.Username); err == nil {
			return model.User{}, model.ErrorUserExists
		}
	}
	if len(user.Email) > 0 {
		if _, err := us.UserByEmail(user.Email); err == nil {
			return model.User{}, model.ErrorUserExists
		}
	}
	if len(user.Phone) > 0 {
		if _, err := us.UserByPhone(user.Phone); err == nil {
			return model.User{}, model.ErrorUserExists
		}
	}

	u := model.User{
		ID:         primitive.NewObjectID().Hex(),
		Active:     true,
		Username:   user.Username,
		Phone:      user.Phone,
		FullName:   user.FullName,
		Scopes:     user.Scopes,
		Email:      user.Email,
		AccessRole: role,
		Anonymous:  isAnonymous,
	}

	return us.AddNewUser(u, password)
}

// AddUserWithFederatedID adds new user with social ID.
func (us *UserStorage) AddUserWithFederatedID(user model.User, provider string, federatedID, role string) (model.User, error) {
	// If there is no error, it means user already exists.
	_, err := us.UserByFederatedID(provider, federatedID)
	if err == nil {
		return model.User{}, model.ErrorUserExists
	}

	// unknown error during user existence check
	if !errors.Is(err, model.ErrUserNotFound) {
		return model.User{}, err
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
	oldUser, err := us.UserByID(userID)
	if err != nil {
		return model.User{}, err
	}

	newUser.Email = strings.ToLower(newUser.Email)
	newUser.Username = strings.ToLower(newUser.Username)
	// use ID from the request
	newUser.ID = userID

	if newUser.Pswd == "" {
		newUser.Pswd = oldUser.Pswd
	}

	if newUser.TFAInfo.Secret == "" {
		newUser.TFAInfo.Secret = oldUser.TFAInfo.Secret
	}

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
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": bson.M{"pswd": model.PasswordHash(password)}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	var ud model.User
	err = us.coll.FindOneAndUpdate(ctx, bson.D{{Key: "_id", Value: id}}, update, opts).Decode(&ud)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.ErrUserNotFound
		}

		return err
	}

	return nil
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
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	update := bson.M{"$set": bson.M{"username": username}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var ud model.User
	err = us.coll.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&ud)
	return err
}

// DeleteUser deletes user by id.
func (us *UserStorage) DeleteUser(id string) error {
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	_, err = us.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// FetchUsers fetches users which name satisfies provided filterString.
// Supports pagination.
func (us *UserStorage) FetchUsers(filterString string, skip, limit int) ([]model.User, int, error) {
	q := bson.D{}
	if len(filterString) > 0 {
		q = bson.D{
			primitive.E{Key: "$or", Value: bson.A{
				bson.D{primitive.E{Key: "username", Value: primitive.Regex{Pattern: filterString, Options: "i"}}},
				bson.D{primitive.E{Key: "email", Value: primitive.Regex{Pattern: filterString, Options: "i"}}},
				bson.D{primitive.E{Key: "phone", Value: primitive.Regex{Pattern: filterString, Options: "i"}}},
			}},
		}
	}

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
func (us *UserStorage) ImportJSON(data []byte, clearOldData bool) error {
	if clearOldData {
		us.ClearAllUserData()
	}
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
func (us *UserStorage) UpdateLoginMetadata(operation, app, userID string, scopes []string, payload map[string]any) {
	hexID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		us.logger.Error("Cannot update login metadata of user",
			logging.FieldUserID, userID,
			logging.FieldError, err)
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
		us.logger.Error("Cannot update login metadata of user",
			logging.FieldUserID, userID,
			logging.FieldError, err)
	}
}

// Close is a no-op.
func (us *UserStorage) Close() {}

// ClearAllUserData - clears the database, used for integration testing
func (us *UserStorage) ClearAllUserData() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*us.timeout)
	defer cancel()

	if _, err := us.coll.DeleteMany(ctx, bson.M{}); err != nil {
		us.logger.Error("Error cleaning all user data",
			logging.FieldError, err)
	}
}
