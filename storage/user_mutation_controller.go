package storage

import (
	"context"
	"errors"
	"time"

	"github.com/madappgang/identifo/v2/jwt"
	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
)

var _umc model.UserMutationController = NewUserStorageController(nil, model.ServerSettings{})

// ====================================
// Data mutation
// ====================================

// CreateUserWithPassword validates password policy, creates password hash and creates new user
// it also responsible to call pre-create and post-create callbacks
func (c *UserStorageController) CreateUserWithPassword(ctx context.Context, u model.User, password string) (model.User, error) {
	// TODO: implement check for isCompromised property
	isCompromised := false
	valid, vr := c.s.PasswordPolicy.Validate(password, isCompromised)
	if !valid {
		ee := []error{}
		for _, r := range vr {
			if !r.Valid {
				ee = append(ee, r.Error())
			}
		}
		// return all violated rules
		return model.User{}, errors.Join(ee...)
	}

	hash, err := jwt.PasswordHash(password, c.s.PasswordHash, []byte(c.s.PasswordHash.Pepper))
	if err != nil {
		return model.User{}, err
	}
	u.PasswordHash = hash

	// TODO: Call pre-create callbacks

	nu, err := c.ums.AddUser(ctx, u)
	if err != nil {
		return model.User{}, err
	}

	// TODO: Call post-create callbacks

	return nu, nil
}

func (c *UserStorageController) UpdateUserPassword(ctx context.Context, userID, password string) error {
	// TODO: implement check for isCompromised property
	isCompromised := false
	valid, vr := c.s.PasswordPolicy.Validate(password, isCompromised)
	if !valid {
		ee := []error{}
		for _, r := range vr {
			if !r.Valid {
				ee = append(ee, r.Error())
			}
		}
		// return all violated rules
		return errors.Join(ee...)
	}

	hash, err := jwt.PasswordHash(password, c.s.PasswordHash, []byte(c.s.PasswordHash.Pepper))
	if err != nil {
		return err
	}

	// TODO: Call pre-change password callback

	user := model.User{
		ID:                  userID,
		PasswordHash:        hash,
		UpdatedAt:           time.Now(),
		LastPasswordResetAt: time.Now(),
	}
	_, err = c.ums.UpdateUser(ctx, user, model.UserFieldsetPassword.UpdateFields()...)
	if err != nil {
		return err
	}

	// TODO: Call  post-change password callback

	return nil
}

func (c *UserStorageController) ChangeBlockStatus(ctx context.Context, userID, reason, whoName, whoID string, blocked bool) error {
	user := model.User{
		ID:        userID,
		UpdatedAt: time.Now(),
		Blocked:   blocked,
	}
	if blocked {
		user.BlockedDetails = &model.UserBlockedDetails{
			Reason:        reason,
			BlockedByName: whoName,
			BlockedById:   whoID,
			BlockedAt:     time.Now(),
		}
	}

	// TODO: Call  pre-block callback

	_, err := c.ums.UpdateUser(ctx, user, model.UserFieldsetBlockStatus.Fields()...)
	if err != nil {
		return err
	}

	// TODO: Call  post-block callback

	return nil
}

func (c *UserStorageController) DeleteUser(ctx context.Context, userID string) error {
	// TODO: Call  pre-update callback

	err := c.ums.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	// TODO: Call  post-update callback
	return nil
}

// Update User is the most complicated part here:
// 1. check if we are changing email and check if we can update it and should send confirmation email
// 2. check each secondary identity is in the list of change, if so - do the following
// 2.1 get the list of the ways user can authenticate to collect the list of identities that should be unique
// 2.2 get the list of unique identities from server config
// 2.3 get the list of immutable identities form server config
// 2.4 get the united list of unique identities
// 2.5 if secondary identity is immutable - fire error
// 2.6 if secondary identity should be unique - check if it already occupied
// 2.7 if it is occupied - return error
// 3. call pre-update callback
// 4. update the fields
// 5. send email confirmation email if needed
// 6. call post-update callback
func (c *UserStorageController) UpdateUser(ctx context.Context, user model.User, fields []string) (model.User, error) {
	// are we trying to change user identity field???
	// then extra care required
	fidToChange := intersect(fields, append(model.UserFieldsetSecondaryIdentity.Fields(), model.UserFieldEmail))
	if len(fidToChange) > 0 {
		// 2.1
		idsInUse, err := c.allIdentityTypesInUse()
		if err != nil {
			return user, err
		}
		// 2.2
		uniqueIDFields := authTypeToField(idsInUse)
		// 2.4
		uniqueIDFields = concatUnique(uniqueIDFields, c.uidfs)

		// let's check email
		if sliceContains(fields, model.UserFieldEmail) {
			// validate email first
			if user.Email != "" && !model.EmailRegexp.MatchString(user.Email) {
				return user, l.LocalizedError{ErrID: l.ErrorAPIRequestBodyEmailInvalid}
			}
			// is it should be unique
			if sliceContains(uniqueIDFields, model.UserFieldEmail) {
				// if it unique it should not be empty
				if len(user.Email) == 0 {
					return user, l.LocalizedError{ErrID: l.ErrorEmailEmpty}
				}

				_, err := c.u.UserBySecondaryID(ctx, model.AuthIdentityTypeEmail, user.Email)
				// 2.7
				if err == nil {
					return user, l.LocalizedError{ErrID: l.ErrorAPIEmailTaken}
				} else if !errors.Is(err, model.ErrUserNotFound) {
					// some other error
					return user, err
				}
			}
			// we need to wipe confirmation email data as we set new email
		}

		// let's check phone number
		if sliceContains(fields, model.UserFieldPhone) {
			// validate phone first
			if user.PhoneNumber != "" && !model.PhoneRegexp.MatchString(user.PhoneNumber) {
				return user, l.LocalizedError{ErrID: l.ErrorInvalidPhone}
			}
			// is it should be unique
			if sliceContains(uniqueIDFields, model.UserFieldPhone) {
				// if it unique it should not be empty
				if len(user.PhoneNumber) == 0 {
					return user, l.LocalizedError{ErrID: l.ErrorPhoneEmpty}
				}

				_, err := c.u.UserBySecondaryID(ctx, model.AuthIdentityTypePhone, user.PhoneNumber)
				// 2.7
				if err == nil {
					return user, l.LocalizedError{ErrID: l.ErrorAPIPhoneTaken}
				} else if !errors.Is(err, model.ErrUserNotFound) {
					// some other error
					return user, err
				}
			}
			// we need to wipe confirmation phone data as we set new email
		}

		// let's check username number
		if sliceContains(fields, model.UserFieldUsername) {
			// validate phone first
			if user.Username != "" && !model.UsernameRegexp.MatchString(user.Username) {
				return user, l.LocalizedError{ErrID: l.ErrorInvalidPhone}
			}
			// is it should be unique
			if sliceContains(uniqueIDFields, model.UserFieldPhone) {
				// if it unique it should not be empty
				if len(user.Nickname) == 0 {
					return user, l.LocalizedError{ErrID: l.ErrorUsernameEmpty}
				}

				_, err := c.u.UserBySecondaryID(ctx, model.AuthIdentityTypeUsername, user.Username)
				// 2.7
				if err == nil {
					return user, l.LocalizedError{ErrID: l.ErrorAPIUsernameTaken}
				} else if !errors.Is(err, model.ErrUserNotFound) {
					// some other error
					return user, err
				}
			}
		}
	}

	// 3.
	// TODO: Call  pre-update callback

	// 4.
	u, err := c.ums.UpdateUser(ctx, user, fields...)
	if err != nil {
		return model.User{}, err
	}

	// 5.
	if sliceContains(fidToChange, model.UserFieldEmail) {
		err = c.SendEmailConfirmation(ctx, u.ID)
		if err != nil {
			// TODO: log error, but no return error as it we have already updated the user
		}
	}
	if sliceContains(fidToChange, model.UserFieldPhone) {
		err = c.SendPhoneConfirmation(ctx, u.ID)
		if err != nil {
			// TODO: log error, but no return error as it we have already updated the user
		}
	}

	// 6.
	// TODO: Call  post-update callback
	return u, nil
}

// return three: shouldSendConfirmationEmail: bool, err: error
// if error nil - user can update his email
// if shouldSendConfirmationEmail is true - we need to send email confirmation to user
func (c *UserStorageController) canUpdateEmail(ctx context.Context, email string) (bool, error) {
	return false, nil
}

// allIdentityTypesInUse returns identity types used now by all apps.
func (c *UserStorageController) allIdentityTypesInUse() ([]model.AuthIdentityType, error) {
	if c.idts == nil {
		ffs, err := c.allActiveFirstFactorStrategies()
		if err != nil {
			return nil, err
		}

		// fill the cache
		c.idts = []model.AuthIdentityType{}
		for _, s := range ffs {
			// only unique local strategies, skipping federated ones
			if s.Type == model.FirstFactorTypeLocal && s.Local != nil && !sliceContains(c.idts, s.Local.Identity) {
				c.idts = append(c.idts, s.Local.Identity)
			}
		}
	}
	return c.idts, nil
}

// allActiveFirstFactorStrategies return all first factor strategies
func (c *UserStorageController) allActiveFirstFactorStrategies() ([]model.FirstFactorStrategy, error) {
	if c.ffs == nil {
		// get strategies first
		apps, err := c.as.FetchApps("")
		if err != nil {
			return nil, err
		}
		// fill the cache
		c.ffs = []model.FirstFactorStrategy{}
		for _, a := range apps {
			for _, s := range a.AuthStrategies {
				if s.Type == model.AuthStrategyFirstFactor && s.FirstFactor != nil {
					c.ffs = append(c.ffs, *s.FirstFactor)
				}
			}
		}
	}

	return c.ffs, nil
}

// authTypeToField converts auth type to associated field in user struct
func authTypeToField(a []model.AuthIdentityType) []string {
	res := []string{}
	for _, t := range a {
		f := t.Field()
		if len(f) > 0 {
			res = append(res, t.Field())
		}
	}
	return res
}

// // authTypeForField converts string to auth type
// func authTypeForField(f string) model.AuthIdentityType {
// 	switch f {
// 	case model.UserFieldEmail:
// 		return model.AuthIdentityTypeEmail
// 	case model.UserFieldUsername:
// 		return model.AuthIdentityTypeUsername
// 	case model.UserFieldPhone:
// 		return model.AuthIdentityTypePhone
// 	}

// 	return model.AuthIdentityTypeAnonymous
// }
