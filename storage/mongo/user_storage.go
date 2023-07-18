package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	usersCollectionName     = "Users"
	usersDataCollectionName = "UsersData"
)

func NewUserStorage(settings model.MongoDatabaseSettings) (*UserStorage, error) {
	if len(settings.ConnectionString) == 0 || len(settings.DatabaseName) == 0 {
		return nil, ErrorEmptyConnectionStringDatabase
	}

	db, err := NewDB(settings.ConnectionString, settings.DatabaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	usersColl := db.Database.Collection(usersCollectionName)

	userNameIndexOptions := &options.IndexOptions{}
	userNameIndexOptions.SetUnique(false)
	userNameIndexOptions.SetSparse(true)
	userNameIndexOptions.SetCollation(&options.Collation{Locale: "en", Strength: 1})

	userNameIndex := &mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "username", Value: bsonx.Int32(int32(1))}},
		Options: userNameIndexOptions,
	}

	emailIndexOptions := &options.IndexOptions{}
	emailIndexOptions.SetUnique(false)
	emailIndexOptions.SetSparse(true)

	emailIndex := &mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "email", Value: bsonx.Int32(int32(1))}},
		Options: emailIndexOptions,
	}

	phoneIndexOptions := &options.IndexOptions{}
	phoneIndexOptions.SetUnique(false)
	phoneIndexOptions.SetSparse(true)

	phoneIndex := &mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "phone", Value: bsonx.Int32(int32(1))}},
		Options: phoneIndexOptions,
	}

	if err := db.EnsureCollectionIndices(usersCollectionName, []mongo.IndexModel{*userNameIndex, *emailIndex, *phoneIndex}); err != nil {
		return &UserStorage{}, fmt.Errorf("failed to create indices for Users collection: %w", err)
	}

	userIDIndex := &mongo.IndexModel{
		Keys: bsonx.Doc{{Key: "user_id", Value: bsonx.Int32(int32(1))}},
	}

	if err := db.EnsureCollectionIndices(usersDataCollectionName, []mongo.IndexModel{*userIDIndex}); err != nil {
		return &UserStorage{}, fmt.Errorf("failed to create indices for UsersData collection: %w", err)
	}

	return &UserStorage{
		coll:     usersColl,
		dataColl: db.Database.Collection(usersDataCollectionName),
		timeout:  5 * time.Second,
	}, nil
}

type UserStorage struct {
	coll     *mongo.Collection
	dataColl *mongo.Collection
	timeout  time.Duration
}

var _ model.UserStorage = &UserStorage{}

func (us *UserStorage) Ready(ctx context.Context) error {
	return us.coll.Database().Client().Ping(ctx, readpref.Primary())
}

func (us *UserStorage) Connect(ctx context.Context) error {
	return us.coll.Database().Client().Connect(ctx)
}

func (us *UserStorage) Close(ctx context.Context) error {
	return us.coll.Database().Client().Disconnect(ctx)
}

func (us *UserStorage) ImportJSON(data []byte, clearOldData bool) error {
	if clearOldData {
		us.ClearAllUserData()
	}
	ud := []model.User{}
	if err := json.Unmarshal(data, &ud); err != nil {
		return err
	}

	// TODO: implement adding new user logic

	return nil
}

func (us *UserStorage) ClearAllUserData() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*us.timeout)
	defer cancel()

	if _, err := us.coll.DeleteMany(ctx, bson.M{}); err != nil {
		log.Printf("Error cleaning all user data: %s\n", err)
	}
}

func (us *UserStorage) UserByID(ctx context.Context, id string) (model.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to convert id to object id: %w", err)
	}

	res := us.coll.FindOne(ctx, bson.D{{Key: "_id", Value: oid}})
	if res.Err() != nil {
		return model.User{}, fmt.Errorf("failed to find user: %w", err)
	}

	var u model.User

	err = res.Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.User{}, l.ErrorUserNotFound
	}

	if err != nil {
		return model.User{}, fmt.Errorf("failed to decode result: %w", err)
	}

	return u, nil
}

func (us *UserStorage) UserBySecondaryID(ctx context.Context, idt model.AuthIdentityType, id string) (model.User, error) {
	switch idt {
	case model.AuthIdentityTypePhone:
		u, err := us.userByPhone(ctx, id)
		if err != nil {
			return model.User{}, fmt.Errorf("failed to get user by phone: %w", err)
		}

		return u, nil
	case model.AuthIdentityTypeEmail:
		u, err := us.userByEmail(ctx, id)
		if err != nil {
			return model.User{}, fmt.Errorf("failed to get user by email: %w", err)
		}

		return u, nil
	case model.AuthIdentityTypeUsername:
		u, err := us.userByUsername(ctx, id)
		if err != nil {
			return model.User{}, fmt.Errorf("failed to get user by username: %w", err)
		}

		return u, nil
	}

	return model.User{}, fmt.Errorf("invalid id type")
}

func (us *UserStorage) userByPhone(ctx context.Context, phone string) (model.User, error) {
	res := us.coll.FindOne(ctx, bson.D{{Key: "phone", Value: phone}})
	if res.Err() != nil {
		return model.User{}, fmt.Errorf("failed to find user by phone: %w", res.Err())
	}

	var u model.User

	err := res.Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.User{}, l.ErrorUserNotFound
	}

	if err != nil {
		return model.User{}, fmt.Errorf("failed to decode result: %w", err)
	}

	return u, nil
}

func (us *UserStorage) userByEmail(ctx context.Context, email string) (model.User, error) {
	res := us.coll.FindOne(ctx, bson.D{{Key: "email", Value: email}})
	if res.Err() != nil {
		return model.User{}, fmt.Errorf("failed to find user by email: %w", res.Err())
	}

	var u model.User

	err := res.Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.User{}, l.ErrorUserNotFound
	}

	if err != nil {
		return model.User{}, fmt.Errorf("failed to decode result: %w", err)
	}

	return u, nil
}

func (us *UserStorage) userByUsername(ctx context.Context, username string) (model.User, error) {
	filter := bson.D{{
		Key: "username",
		Value: primitive.Regex{
			Pattern: fmt.Sprintf("^%s$", regexp.QuoteMeta(username)),
			Options: "i",
		},
	}}

	var u model.User

	res := us.coll.FindOne(ctx, filter)
	if res.Err() != nil {
		return model.User{}, fmt.Errorf("failed to find user by username: %w", res.Err())
	}

	err := res.Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.User{}, l.ErrorUserNotFound
	}

	if err != nil {
		return model.User{}, fmt.Errorf("failed to decode result: %w", err)
	}

	return u, nil
}

func (us *UserStorage) UserByFederatedID(ctx context.Context, idType model.UserFederatedType, userIdentityTypeOther, externalID string) (model.User, error) {
	sid := string(idType) + ":" + userIdentityTypeOther + ":" + externalID

	res := us.coll.FindOne(ctx, bson.D{{Key: "federated_ids", Value: sid}})
	if res.Err() != nil {
		return model.User{}, fmt.Errorf("failed to find user by federated id: %w", res.Err())
	}

	var u model.User

	if err := res.Decode(&u); err != nil {
		return model.User{}, fmt.Errorf("failed to decode result: %w", err)
	}

	return u, nil
}

func (us *UserStorage) UserData(ctx context.Context, userID string, fields ...model.UserDataField) (model.UserData, error) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return model.UserData{}, fmt.Errorf("failed to convert id to object id: %w", err)
	}

	res := us.dataColl.FindOne(ctx, bson.D{{Key: "user_id", Value: oid}})
	if res.Err() != nil {
		return model.UserData{}, fmt.Errorf("failed to find user: %w", err)
	}

	var u model.UserData

	err = res.Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.UserData{}, l.ErrorUserNotFound
	}

	if err != nil {
		return model.UserData{}, fmt.Errorf("failed to decode result: %w", err)
	}

	return model.FilterUserDataFields(u, fields...), nil
}
