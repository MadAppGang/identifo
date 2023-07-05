package storage

import "context"

// we need to create challenge for user
// challenge type is phone confirmation
// user should enter the code from SMS
// the challenge has a specific TTL
func (c *UserStorageController) SendPhoneConfirmation(ctx context.Context, userID string) error {
	// TODO: Create auth data storage to keep all challenges and use enrollments there
	// TODO: Add log events there as well
	// challenge, err := c.authDataStorage.CreateChallenge(userID, challengeType, challengeTTL)
	// err = c.smsSender.Send(ctx, phoneNumber, challenge.Code)
	return nil
}

// we need to create challenge for user
// challenge type is phone confirmation
// user should enter the code from EMAIL
// it could be done with
// the challenge has a specific TTL
// email confirmation is a link with a code in it
func (c *UserStorageController) SendEmailConfirmation(ctx context.Context, userID string) error {
	return nil
}
