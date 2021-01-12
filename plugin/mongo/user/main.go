package mongo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/plugin/shared"
	"github.com/madappgang/identifo/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	conn := os.Getenv("DB_CONN")
	if conn == "" {
		panic("Empty DB_CONN")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		panic("Empty DB_NAME")
	}

	db, err := NewDB(conn, dbName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	us, err := NewUserStorage(db)
	if err != nil {
		panic(err)
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"user_storage": &shared.UserStorageGRPCPlugin{
				Impl: us,
			},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

// NewDB creates new database connection.
func NewDB(conn string, dbName string) (*DB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(conn))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		return nil, err
	}
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	db := &DB{
		Client:   client,
		Database: client.Database(dbName),
	}
	return db, nil
}

// DB is database connection structure.
type DB struct {
	Database *mongo.Database
	Client   *mongo.Client
}

// Close closes database connection.
func (db *DB) Close() error {
	return db.Client.Disconnect(context.TODO())
}

// EnsureCollectionIndices creates indices on a collection.
func (db *DB) EnsureCollectionIndices(collectionName string, newIndices []mongo.IndexModel) error {
	coll := db.Database.Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	for _, newIndex := range newIndices {
		if _, err := coll.Indexes().CreateOne(ctx, newIndex); err != nil && strings.Contains(err.Error(), "already exists with different options") {
			name, err := generateIndexName(newIndex)
			if err != nil {
				return err
			}
			if _, err = coll.Indexes().DropOne(ctx, name); err != nil {
				return err
			}
			if _, err = coll.Indexes().CreateOne(ctx, newIndex); err != nil {
				return err
			}
		}
	}
	return nil
}

func generateIndexName(index mongo.IndexModel) (string, error) {
	if index.Options != nil && index.Options.Name != nil {
		return *index.Options.Name, nil
	}

	name := bytes.NewBufferString("")
	first := true

	keys, ok := index.Keys.(bsonx.Doc)
	if !ok {
		return "", fmt.Errorf("Incorrect index keys type - expecting bsonx.Doc")
	}
	for _, elem := range keys {
		if !first {
			if _, err := name.WriteRune('_'); err != nil {
				return "", err
			}
		}
		if _, err := name.WriteString(elem.Key); err != nil {
			return "", err
		}
		if _, err := name.WriteRune('_'); err != nil {
			return "", err
		}

		var value string

		switch elem.Value.Type() {
		case bsontype.Int32:
			value = fmt.Sprintf("%d", elem.Value.Int32())
		case bsontype.Int64:
			value = fmt.Sprintf("%d", elem.Value.Int64())
		case bsontype.String:
			value = elem.Value.StringValue()
		default:
			return "", fmt.Errorf("Incorrect index value type %s", elem.Value.Type())
		}

		if _, err := name.WriteString(value); err != nil {
			return "", err
		}

		first = false
	}
	return name.String(), nil
}

func isErrNotFound(err error) bool {
	return strings.Contains(err.Error(), "no documents in result")
}

func isErrDuplication(err error) bool {
	return strings.Contains(err.Error(), "duplicate key")
}

const usersCollectionName = "Users"

// NewUserStorage creates and inits MongoDB user storage.
func NewUserStorage(db *DB) (shared.UserStorage, error) {
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

	err := db.EnsureCollectionIndices(usersCollectionName, []mongo.IndexModel{*userNameIndex, *emailIndex, *phoneIndex})
	return us, err
}

// UserStorage implements user storage interface.
type UserStorage struct {
	coll    *mongo.Collection
	timeout time.Duration
}

// UserByID returns user by its ID.
func (us *UserStorage) UserByID(id string) (*proto.User, error) {
	hexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	u := proto.User{}
	if err := us.coll.FindOne(ctx, bson.M{"_id": hexID}).Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

// UserByEmail returns user by their email.
func (us *UserStorage) UserByEmail(email string) (*proto.User, error) {
	if email == "" {
		return nil, model.ErrorWrongDataFormat
	}
	email = strings.ToLower(email)

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	u := proto.User{}
	if err := us.coll.FindOne(ctx, bson.M{"email": email}).Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

// UserByFederatedID returns user by federated ID.
func (us *UserStorage) UserByFederatedID(provider proto.FederatedIdentityProvider, id string) (*proto.User, error) {
	sid := provider.String() + ":" + id

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	u := proto.User{}
	if err := us.coll.FindOne(ctx, bson.M{"federated_ids": sid}).Decode(&u); err != nil {
		return nil, errors.New("User not found")
	}
	//clear password hash
	u.PasswordHash = ""
	return &u, nil
}

// UserExists checks if user with provided name exists.
func (us *UserStorage) UserExists(name string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	strictPattern := "^" + name + "$"
	q := bson.D{primitive.E{Key: "username", Value: primitive.Regex{Pattern: strictPattern, Options: "i"}}}

	u := proto.User{}
	err := us.coll.FindOne(ctx, q).Decode(&u)
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

// UserByPhone fetches user by phone number.
func (us *UserStorage) UserByPhone(phone string) (*proto.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	u := proto.User{}
	if err := us.coll.FindOne(ctx, bson.M{"phone": phone}).Decode(&u); err != nil {
		return nil, err
	}
	u.PasswordHash = ""
	return &u, nil
}

// UserByNamePassword returns user by name and password.
func (us *UserStorage) UserByNamePassword(name, password string) (*proto.User, error) {
	strictPattern := "^" + strings.ReplaceAll(name, "+", "\\+") + "$"
	q := bson.D{primitive.E{Key: "username", Value: primitive.Regex{Pattern: strictPattern, Options: "i"}}}

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	u := proto.User{}
	if err := us.coll.FindOne(ctx, q).Decode(&u); err != nil {
		return nil, errors.New("User not found")
	}

	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("User not found")
	}
	//clear password hash
	u.PasswordHash = ""
	return &u, nil
}

// AddNewUser adds new user to the database.
func (us *UserStorage) AddNewUser(u *proto.User, password string) (*proto.User, error) {
	u.Email = strings.ToLower(u.Email)

	u.Id = primitive.NewObjectID().String()
	if len(password) > 0 {
		u.PasswordHash = PasswordHash(password)
	}
	u.NumOfLogins = 0

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	if _, err := us.coll.InsertOne(ctx, u); err != nil {
		if isErrDuplication(err) {
			return nil, model.ErrorUserExists
		}
		return nil, err
	}
	return u, nil
}

// AddUserByPhone registers new user with phone number.
func (us *UserStorage) AddUserByPhone(phone, role string) (*proto.User, error) {
	u := proto.User{
		Id:          primitive.NewObjectID().String(),
		Username:    phone,
		IsActive:    true,
		Phone:       phone,
		AccessRole:  role,
		NumOfLogins: 0,
	}

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	if _, err := us.coll.InsertOne(ctx, &u); err != nil {
		if isErrDuplication(err) {
			return nil, model.ErrorUserExists
		}
		return nil, err
	}
	return &u, nil
}

// AddUserByNameAndPassword registers new user.
func (us *UserStorage) AddUserByNameAndPassword(username, password, role string, isAnonymous bool) (*proto.User, error) {
	u := proto.User{
		Id:          primitive.NewObjectID().String(),
		IsActive:    true,
		Username:    username,
		AccessRole:  role,
		IsAnonymous: isAnonymous,
	}
	if shared.EmailRegexp.MatchString(u.Username) {
		u.Email = u.Username
	}
	if shared.PhoneRegexp.MatchString(u.Username) {
		u.Phone = u.Username
	}
	return us.AddNewUser(&u, password)
}

// AddUserWithFederatedID adds new user with social ID.
func (us *UserStorage) AddUserWithFederatedID(provider proto.FederatedIdentityProvider, federatedID, role string) (*proto.User, error) {
	// If there is no error, it means user already exists.
	if _, err := us.UserByFederatedID(provider, federatedID); err == nil {
		return nil, errors.New("User exists")
	}

	sid := provider.String() + ":" + federatedID
	u := proto.User{
		Id:           primitive.NewObjectID().String(),
		IsActive:     true,
		Username:     sid,
		AccessRole:   role,
		FederatedIds: []string{sid},
	}
	return us.AddNewUser(&u, "")
}

// UpdateUser updates user in MongoDB storage.
func (us *UserStorage) UpdateUser(userID string, newUser *proto.User) (*proto.User, error) {
	hexID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newUser.Email = strings.ToLower(newUser.Email)

	// use ID from the request
	newUser.Id = hexID.String()

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	update := bson.M{"$set": newUser}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	u := proto.User{}
	if err := us.coll.FindOneAndUpdate(ctx, bson.M{"_id": hexID}, update, opts).Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

// ResetPassword sets new user's password.
func (us *UserStorage) ResetPassword(id, password string) error {
	hexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": bson.M{"pswd": PasswordHash(password)}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	u := proto.User{}
	err = us.coll.FindOneAndUpdate(ctx, bson.M{"_id": hexID}, update, opts).Decode(&u)
	return err
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

	u := proto.User{}
	err = us.coll.FindOneAndUpdate(ctx, bson.M{"_id": hexID.String()}, update, opts).Decode(&u)
	return err
}

// IDByName returns userID by name.
func (us *UserStorage) IDByName(name string) (string, error) {
	strictPattern := "^" + name + "$"
	q := bson.D{primitive.E{Key: "username", Value: primitive.Regex{Pattern: strictPattern, Options: "i"}}}

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	u := proto.User{}
	if err := us.coll.FindOne(ctx, q).Decode(&u); err != nil {
		return "", model.ErrorNotFound
	}

	user := proto.User{}
	if !user.IsActive {
		return "", errors.New("Inactive user")
	}
	return user.Id, nil
}

// DeleteUser deletes user by id.
func (us *UserStorage) DeleteUser(id string) error {
	hexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), us.timeout)
	defer cancel()

	_, err = us.coll.DeleteOne(ctx, bson.M{"_id": hexID})
	return err
}

// FetchUsers fetches users which name satisfies provided filterString.
// Supports pagination.
func (us *UserStorage) FetchUsers(filterString string, skip, limit int) ([]*proto.User, int, error) {
	q := bson.D{primitive.E{Key: "username", Value: primitive.Regex{Pattern: filterString, Options: "i"}}}

	ctx, cancel := context.WithTimeout(context.Background(), 2*us.timeout)
	defer cancel()

	total, err := us.coll.CountDocuments(ctx, q)
	if err != nil {
		return []*proto.User{}, 0, err
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{primitive.E{Key: "username", Value: 1}})
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))

	curr, err := us.coll.Find(ctx, q, findOptions)
	if err != nil {
		return nil, 0, err
	}

	users := []*proto.User{}
	if err = curr.All(ctx, &users); err != nil {
		return nil, 0, err
	}
	return users, int(total), err
}

// ImportJSON imports data from JSON.
func (us *UserStorage) ImportJSON(data []byte) error {
	ud := []*proto.User{}
	if err := json.Unmarshal(data, &ud); err != nil {
		return err
	}
	for _, u := range ud {
		pswd := u.PasswordHash
		u.PasswordHash = ""
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

	u := proto.User{}
	if err := us.coll.FindOneAndUpdate(ctx, bson.M{"_id": hexID.String()}, update).Decode(&u); err != nil {
		log.Printf("Cannot update login metadata of user %s: %s\n", userID, err)
	}
}

// Close is a no-op.
func (us *UserStorage) Close() {}

// PasswordHash creates hash with salt for password.
func PasswordHash(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash)
}
