package boltdb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/madappgang/identifo/model"
	"github.com/rs/xid"
)

const (
	// InviteBucket is a name for bucket with invites.
	InviteBucket = "Invites"
)

// InviteStorage is a BoltDB invite storage.
type InviteStorage struct {
	db *bolt.DB
}

// NewInviteStorage creates a BoltDB invites storage.
func NewInviteStorage(db *bolt.DB) (model.InviteStorage, error) {
	is := &InviteStorage{db: db}
	// Ensure that we have needed bucket in the database.
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(InviteBucket)); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return is, nil
}

// Save creates and saves new invite to a database.
func (is *InviteStorage) Save(email, inviteToken, role, appID, createdBy string, expiresAt time.Time) error {
	return is.db.Update(func(tx *bolt.Tx) error {
		ib := tx.Bucket([]byte(InviteBucket))

		invite := model.Invite{
			ID:        xid.New().String(),
			AppID:     appID,
			Token:     inviteToken,
			Valid:     true,
			Email:     email,
			Role:      role,
			CreatedBy: createdBy,
			CreatedAt: time.Now(),
			ExpiresAt: expiresAt,
		}

		if err := invite.Validate(); err != nil {
			return err
		}

		data, err := json.Marshal(invite)
		if err != nil {
			return err
		}

		return ib.Put([]byte(invite.ID), data)
	})
}

// GetByEmail returns valid and not expired invite by email.
func (is *InviteStorage) GetByEmail(email string) (model.Invite, error) {
	var invite model.Invite

	err := is.db.View(func(tx *bolt.Tx) error {
		ib := tx.Bucket([]byte(InviteBucket))

		iterErr := ib.ForEach(func(k, v []byte) error {
			var res model.Invite
			if err := json.Unmarshal(v, &res); err != nil {
				return err
			}

			if res.Email == email && res.Valid == true && res.ExpiresAt.After(time.Now()) {
				invite = res
				return nil
			}

			return nil
		})
		if iterErr != nil {
			return iterErr
		}

		return nil
	})

	return invite, err
}

// GetByID returns invite by its ID.
func (is *InviteStorage) GetByID(id string) (model.Invite, error) {
	var invite model.Invite

	err := is.db.View(func(tx *bolt.Tx) error {
		ib := tx.Bucket([]byte(InviteBucket))

		res := ib.Get([]byte(id))
		return json.Unmarshal(res, &invite)
	})
	if err != nil {
		return model.Invite{}, err
	}

	return model.Invite{}, nil
}

// GetAll returns all active invites by default.
// To get an invalid invites need to set withInvalid argument to true.
func (is *InviteStorage) GetAll(withInvalid bool, skip, limit int) ([]model.Invite, int, error) {
	var (
		invites []model.Invite
		total   int
	)

	err := is.db.View(func(tx *bolt.Tx) error {
		ib := tx.Bucket([]byte(InviteBucket))

		return ib.ForEach(func(k, v []byte) error {
			var invite model.Invite
			if err := json.Unmarshal(v, &invite); err != nil {
				return err
			}

			if withInvalid == false && invite.Valid == false {
				return nil
			}

			total++
			skip--
			if skip > -1 || (limit != 0 && len(invites) == limit) {
				return nil
			}
			invites = append(invites, invite)
			return nil
		})
	})
	if err != nil {
		return []model.Invite{}, 0, err
	}
	return invites, total, nil
}

// InvalidateAllByEmail invalidates all invites by email.
func (is *InviteStorage) InvalidateAllByEmail(email string) error {
	return is.db.Update(func(tx *bolt.Tx) error {
		ib := tx.Bucket([]byte(InviteBucket))

		iterErr := ib.ForEach(func(k, v []byte) error {
			var invite model.Invite
			if err := json.Unmarshal(v, &invite); err != nil {
				return err
			}

			if invite.Email == email {
				invite.Valid = false

				data, err := json.Marshal(invite)
				if err != nil {
					return err
				}

				return ib.Put([]byte(invite.ID), data)
			}

			return nil
		})
		if iterErr != nil {
			return iterErr
		}

		return nil
	})
}

// InvalidateByID invalidates specific invite by its ID.
func (is *InviteStorage) InvalidateByID(id string) error {
	return is.db.Update(func(tx *bolt.Tx) error {
		ib := tx.Bucket([]byte(InviteBucket))

		iterErr := ib.ForEach(func(k, v []byte) error {
			var invite model.Invite
			if err := json.Unmarshal(v, &invite); err != nil {
				return err
			}

			if invite.ID == id {
				invite.Valid = false
				data, err := json.Marshal(invite)
				if err != nil {
					return err
				}

				return ib.Put(k, data)
			}

			return nil
		})
		if iterErr != nil {
			return iterErr
		}
		return nil
	})
}
