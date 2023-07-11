package model

// we support OIDC scopes, not claims
// https://openid.net/specs/openid-connect-core-1_0.html#ScopeClaims
const (
	// OfflineScope is a scope value to request refresh token.
	OfflineScope           = "offline"
	OIDCScope              = "openid"
	EmailScope             = "email"
	PhoneScope             = "phone"
	ProfileScope           = "profile"
	AddressScope           = "address"
	CustomScopePrefix      = "custom:"
	IDTokenScope           = "id"
	TenantScope            = "tenant"
	AccessTokenScopePrefix = "access:"
	TenantScopePrefix      = "tenant:"    // tenant:123 - request tenant data only for tenant 123
	TenantScopeAll         = "tenant:all" // "tenant:all" - return all scopes for all ten
)

func FieldsetForScopes(scopes []string) []string {
	fieldset := []string{}
	for _, scope := range scopes {
		switch scope {
		case OIDCScope:
			fieldset = append(fieldset, UserFieldsetMap[UserFieldsetScopeOIDC]...)
		case EmailScope:
			fieldset = append(fieldset, UserFieldsetMap[UserFieldsetScopeEmail]...)
		case PhoneScope:
			fieldset = append(fieldset, UserFieldsetMap[UserFieldsetScopePhone]...)
		case AddressScope:
			fieldset = append(fieldset, UserFieldsetMap[UserFieldsetScopeAddress]...)
		case ProfileScope:
			fieldset = append(fieldset, UserFieldsetMap[UserFieldsetScopeProfile]...)
		}
	}
	return removeDuplicate(fieldset)
}

// removeDuplicate
func removeDuplicate[T comparable](slice []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range slice {
		if _, ok := allKeys[item]; !ok {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
