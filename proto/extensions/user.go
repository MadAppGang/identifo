package extensions

import "github.com/madappgang/identifo/proto"

// Sanitize removes all sensitive data.
func SanitizeUser(u *proto.User) {
	if u == nil {
		return
	}
	u.PasswordHash = ""

	if u.TfaInfo != nil {
		u.TfaInfo.Secret = ""
	}
}
