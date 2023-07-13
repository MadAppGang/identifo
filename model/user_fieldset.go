package model

type UserFieldset string

const (
	UserFieldsetBasic             UserFieldset = "basic"
	UserFieldsetAll               UserFieldset = "all"
	UserFieldsetBlockStatus       UserFieldset = "block_status"
	UserFieldsetPassword          UserFieldset = "password"
	UserFieldsetSecondaryIdentity UserFieldset = "secondary_identity"
	UserFieldsetUpdatableByUser   UserFieldset = "updatable_by_user"
	UserFieldsetScopeOIDC         UserFieldset = OIDCScope
	UserFieldsetScopeEmail        UserFieldset = EmailScope
	UserFieldsetScopePhone        UserFieldset = PhoneScope
	UserFieldsetScopeProfile      UserFieldset = ProfileScope
	UserFieldsetScopeAddress      UserFieldset = AddressScope
	UserFieldsetInviteToken       UserFieldset = "invite_token"

	// TODO: Add fieldset cases for other cases.

	UserFieldEmail    = "Email"
	UserFieldID       = "ID"
	UserFieldUsername = "Username"
	UserFieldPhone    = "PhoneNumber"
)

// TODO: Add more fieldset for a map.
var UserFieldsetMap = map[UserFieldset][]string{
	UserFieldsetBasic: {
		"ID",
		"Name",
		"Username",
		"Email",
		"GivenName",
		"FamilyName",
		"MiddleName",
		"Nickname",
		"PreferredUsername",
		"PhoneNumber",
		"Locale",
	},
	UserFieldsetPassword: {
		"ID",
		"PasswordHash",
		"PasswordResetRequired",
		"PasswordChangeForced",
		"LastPasswordResetAt",
	},
	UserFieldsetUpdatableByUser: {
		"Username",
		"Email",
		"GivenName",
		"FamilyName",
		"MiddleName",
		"Nickname",
		"PreferredUsername",
		"PhoneNumber",
		"Profile",
		"Picture",
		"Website",
		"Gender",
		"Birthday",
		"Timezone",
		"Locale",
		"Address",
	},
	UserFieldsetSecondaryIdentity: {
		"Username",
		"PhoneNumber",
	},
	UserFieldsetScopeOIDC: {
		"Profile",
		"Picture",
		"Website",
		"Gender",
		"Birthday",
		"Timezone",
		"Locale",
	},
	UserFieldsetScopeEmail: {
		"Email",
		"EmailVerificationDetails",
	},
	UserFieldsetScopePhone: {
		"PhoneNumber",
		"PhoneVerificationDetails",
	},
	UserFieldsetScopeProfile: {
		"Username",
		"GivenName",
		"FamilyName",
		"MiddleName",
		"Nickname",
		"PreferredUsername",
	},
	UserFieldsetScopeAddress: {
		"Address",
	},
	UserFieldsetInviteToken: {
		"Email",
	},
}

func (f UserFieldset) Fields() []string {
	return UserFieldsetMap[f]
}

var ImmutableFields = []string{"ID", "CreatedAt"}

// the list for update for specific fieldset
// usually it includes UpdateAt and excludes ID from the list
func (f UserFieldset) UpdateFields() []string {
	r := UserFieldsetMap[f]
	r = subtract(r, ImmutableFields)
	r = append(r, "UpdatedAt")
	return r
}

// subtract sl2 from sl1
func subtract[T comparable](sl1, sl2 []T) []T {
	var diff []T

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range sl1 {
			found := false
			for _, s2 := range sl2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			sl1, sl2 = sl2, sl1
		}
	}

	return diff
}
