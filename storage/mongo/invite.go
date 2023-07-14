package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const invitesCollectionName = "Invites"

// InviteStorage is a MongoDB invite storage.
type InviteStorage struct {
	coll    *mongo.Collection
	timeout time.Duration
}

// NewInviteStorage creates a MongoDB invite storage.
func NewInviteStorage(settings model.MongoDatabaseSettings) (model.InviteStorage, error) {
	if len(settings.ConnectionString) == 0 || len(settings.DatabaseName) == 0 {
		return nil, ErrorEmptyConnectionStringDatabase
	}

	// create database
	db, err := NewDB(settings.ConnectionString, settings.DatabaseName)
	if err != nil {
		return nil, err
	}

	coll := db.Database.Collection(invitesCollectionName)
	return &InviteStorage{coll: coll, timeout: 30 * time.Second}, nil
}

// Save creates and saves new invite to a database.
func (is *InviteStorage) Save(email, inviteToken, role, appID, createdBy string, expiresAt time.Time) error {
	if len(inviteToken) == 0 {
		return l.ErrorAPIDataError
	}

	ctx, cancel := context.WithTimeout(context.Background(), is.timeout)
	defer cancel()

	i := model.Invite{
		ID:        primitive.NewObjectID().Hex(),
		AppID:     appID,
		Token:     inviteToken,
		Archived:  false,
		Email:     email,
		Role:      role,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	if err := i.Validate(); err != nil {
		return err
	}

	_, err := is.coll.InsertOne(ctx, i)
	return err
}

// GetByEmail returns valid and not expired invite by email.
func (is *InviteStorage) GetByEmail(email string) (model.Invite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), is.timeout)
	defer cancel()

	filter := bson.M{
		"email":     email,
		"archived":  false,
		"expiresAt": bson.M{"$gt": time.Now()},
	}

	var invite model.Invite
	if err := is.coll.FindOne(ctx, filter).Decode(&invite); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Invite{}, l.NewError(l.ErrorNotFound, "invite")
		}
		return model.Invite{}, err
	}
	return invite, nil
}

// GetByID returns invite by its ID.
func (is *InviteStorage) GetByID(id string) (model.Invite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), is.timeout)
	defer cancel()

	filter := bson.M{"_id": id}

	var invite model.Invite
	if err := is.coll.FindOne(ctx, filter).Decode(&invite); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Invite{}, l.NewError(l.ErrorNotFound, "invite")
		}
		return model.Invite{}, err
	}
	return invite, nil
}

// GetAll returns all active invites by default.
// To get an invalid invites need to set withInvalid argument to true.
func (is *InviteStorage) GetAll(withArchived bool, skip, limit int) ([]model.Invite, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), is.timeout)
	defer cancel()

	filter := bson.M{}
	if !withArchived {
		filter["archived"] = false
	}

	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	curr, err := is.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	var invites []model.Invite
	if err = curr.All(ctx, &invites); err != nil {
		return nil, 0, err
	}
	return invites, len(invites), nil
}

// ArchiveAllByEmail invalidates all invites by email.
func (is *InviteStorage) ArchiveAllByEmail(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), is.timeout)
	defer cancel()

	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"archived": true}}

	_, err := is.coll.UpdateMany(ctx, filter, update)
	return err
}

// ArchiveByID invalidates specific invite by its ID.
func (is *InviteStorage) ArchiveByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), is.timeout)
	defer cancel()

	hexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": hexID.Hex()}
	update := bson.M{"$set": bson.M{"archived": true}}

	_, err = is.coll.UpdateOne(ctx, filter, update)
	return err
}

// Close is a no-op.
func (is *InviteStorage) Close() {}
