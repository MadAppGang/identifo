package model

// FederatedIdentityProvider is an external federated identity provider type.
// If you are missing the provider you need, please feel free to add it here.
type FederatedIdentityProvider string

var (
	// FacebookIDProvider is a Facebook ID provider.
	FacebookIDProvider FederatedIdentityProvider = "FACEBOOK"
	// GoogleIDProvider is a Google ID provider.
	GoogleIDProvider FederatedIdentityProvider = "GOOGLE"
	// TwitterIDProvider is a Twitter ID provider.
	TwitterIDProvider FederatedIdentityProvider = "TWITTER"
	// AppleIDProvider is an Apple ID provider.
	AppleIDProvider FederatedIdentityProvider = "APPLE"
)

// IsValid has to be called everywhere input happens, otherwise you risk to operate on bad data - no guarantees.
func (fid FederatedIdentityProvider) IsValid() bool {
	switch fid {
	case FacebookIDProvider, GoogleIDProvider, TwitterIDProvider, AppleIDProvider:
		return true
	}
	return false
}

// AppleInfo represents the information needed for Sign In with Apple.
type AppleInfo struct {
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}
