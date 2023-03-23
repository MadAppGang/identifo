package model

import "time"

//Add metadata
// https://auth0.com/docs/manage-users/user-accounts/metadata/metadata-fields-data

//2FA connectors
// https://community.auth0.com/t/when-will-api-v2-users-userid-authenticators-be-documented/52722

//Guardian authenticators
//https://auth0.com/docs/secure/multi-factor-authentication/auth0-guardian

// User is a user identitty data
// do not mix it with login strategy
//https://auth0.com/docs/manage-users/user-accounts/user-profiles
type User struct {
	ID              string   `json:"id" bson:"_id"`
  ExternalID string //external system user ID
	Username        string   `json:"username" bson:"username"`
	Email           string   `json:"email" bson:"email"`
  EmailVerified string
  FamilyName string
  GivenName string
  PhoneNumber string
  PhoneVerificationDetailes struct {
    VerifiedAt time.Time
    VerifiedDetails string
  }
  AvatarURL string

  Authenticators []UserAuthenticator
  Identities []UserIdentity
  LastIP string
  LastLoginAt time.Time
  LastPasswordResetAt time.Time
  LoginsCount integer
  MFAEnrollments []UserMFA


  PasswordHash string
  Blocked bool
  BlockedDetails struct {
    BlockedReason string
    BlockedAt time.Time
    BlockedBy string
  }
  EmailVerificationDetailes struct {
    VerifiedAt time.Time
    VerifiedDetails string
  }
  CreatedAt time.Time
  UpdatedAt time.Time
}

//https://auth0.com/docs/manage-users/user-migration/bulk-user-import-database-schema-and-examples
type UserMFA struct {
  
}

func maskLeft(s string, hideFraction int) string {
	rs := []rune(s)
	for i := 0; i < len(rs)-len(rs)/hideFraction; i++ {
		rs[i] = '*'
	}
	return string(rs)
}

// Sanitized returns data structure without sensitive information
func (u User) Sanitized() User {
	u.Pswd = ""
	u.TFAInfo.Secret = ""
	u.TFAInfo.HOTPCounter = 0
	u.TFAInfo.HOTPExpiredAt = time.Time{}
	return u
}

func (u *User) AddFederatedId(provider string, id string) string {
	fid := provider + ":" + id
	for _, ele := range u.FederatedIDs {
		if ele == fid {
			return fid
		}
	}
	u.FederatedIDs = append(u.FederatedIDs, fid)
	return fid
}

// SanitizedTFA returns data structure with masked sensitive data
func (u User) SanitizedTFA() User {
	u.Sanitized()
	if len(u.Email) > 0 {
		emailParts := strings.Split(u.Email, "@")
		u.Email = maskLeft(emailParts[0], 2) + "@" + maskLeft(emailParts[1], 2)
	}

	if len(u.Phone) > 0 {
		u.Phone = maskLeft(u.Phone, 3)
	}
	return u
}

// Deanonimized returns model with all fields set for deanonimized user
func (u User) Deanonimized() User {
	u.Anonymous = false
	return u
}

// TFAInfo encapsulates two-factor authentication user info.
type TFAInfo struct {
	IsEnabled     bool      `json:"is_enabled" bson:"is_enabled"`
	HOTPCounter   int       `json:"hotp_counter" bson:"hotp_counter"`
	HOTPExpiredAt time.Time `json:"hotp_expired_at" bson:"hotp_expired_at"`
	Email         string    `json:"email" bson:"email"`
	Phone         string    `json:"phone" bson:"phone"`
	Secret        string    `json:"secret" bson:"secret"`
}

// UserFromJSON deserialize user data from JSON.
func UserFromJSON(d []byte) (User, error) {
	var user User
	if err := json.Unmarshal(d, &user); err != nil {
		log.Println("Error while unmarshal user:", err)
		return User{}, err
	}
	return user, nil
}

// PasswordHash creates hash with salt for password.
func PasswordHash(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash)
}

// RandomPassword creates random password
func RandomPassword(length int) string {
	rand.Seed(time.Now().UnixNano())
	return randSeq(length)
}

var rndPassLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890?!@#")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = rndPassLetters[rand.Intn(len(rndPassLetters))]
	}
	return string(b)
}

