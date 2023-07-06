package mongo

// import (
// 	"context"
// 	"time"

// 	"github.com/madappgang/identifo/v2/model"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// 	"go.mongodb.org/mongo-driver/x/bsonx"
// )

// const (
// 	verificationCodesCollectionName = "VerificationCodes"
// 	// verificationCodesExpirationTime specifies time before deleting records.
// 	verificationCodesExpirationTime = 5 * time.Minute
// )

// // NewVerificationCodeStorage creates and inits MongoDB verification code storage.
// func NewVerificationCodeStorage(settings model.MongoDatabaseSettings) (model.VerificationCodeStorage, error) {
// 	if len(settings.ConnectionString) == 0 || len(settings.DatabaseName) == 0 {
// 		return nil, ErrorEmptyConnectionStringDatabase
// 	}

// 	// create database
// 	db, err := NewDB(settings.ConnectionString, settings.DatabaseName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	coll := db.Database.Collection(verificationCodesCollectionName)
// 	vcs := &VerificationCodeStorage{coll: coll, timeout: 30 * time.Second}

// 	phoneIndexOptions := &options.IndexOptions{}
// 	phoneIndexOptions.SetUnique(true)

// 	phoneIndex := &mongo.IndexModel{
// 		Keys:    bsonx.Doc{{Key: "phone", Value: bsonx.Int32(int32(1))}},
// 		Options: phoneIndexOptions,
// 	}

// 	codeIndexOptions := &options.IndexOptions{}
// 	codeIndexOptions.SetUnique(true)

// 	codeIndex := &mongo.IndexModel{
// 		Keys:    bsonx.Doc{{Key: "code", Value: bsonx.Int32(int32(1))}},
// 		Options: codeIndexOptions,
// 	}

// 	createdAtOptions := &options.IndexOptions{}
// 	createdAtOptions.SetUnique(true)
// 	createdAtOptions.SetExpireAfterSeconds(int32(verificationCodesExpirationTime.Seconds()))

// 	createdAtIndex := &mongo.IndexModel{
// 		Keys:    bsonx.Doc{{Key: "createdAt", Value: bsonx.Int32(int32(1))}},
// 		Options: createdAtOptions,
// 	}

// 	err = db.EnsureCollectionIndices(verificationCodesCollectionName, []mongo.IndexModel{*phoneIndex, *codeIndex, *createdAtIndex})
// 	return vcs, err
// }

// // VerificationCodeStorage implements verification code storage interface.
// type VerificationCodeStorage struct {
// 	coll    *mongo.Collection
// 	timeout time.Duration
// }

// // IsVerificationCodeFound checks whether verification code can be found.
// func (vcs *VerificationCodeStorage) IsVerificationCodeFound(phone, code string) (bool, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	var c interface{}
// 	if err := vcs.coll.FindOneAndDelete(ctx, bson.M{"phone": phone, "code": code}).Decode(&c); err != nil {
// 		if isErrNotFound(err) {
// 			return false, nil
// 		}
// 		return false, err
// 	}
// 	return true, nil
// }

// // CreateVerificationCode inserts new verification code to the database.
// func (vcs *VerificationCodeStorage) CreateVerificationCode(phone, code string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	if _, err := vcs.coll.DeleteMany(ctx, bson.M{"phone": phone}); err != nil {
// 		return err
// 	}

// 	_, err := vcs.coll.InsertOne(ctx, bson.M{"phone": phone, "code": code, "createdAt": time.Now()})
// 	return err
// }

// // Close is a no-op here.
// func (vcs *VerificationCodeStorage) Close() {}
