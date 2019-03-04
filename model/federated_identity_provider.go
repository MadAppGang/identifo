package model

//FederatedIdentityProvider external federated identity provider type, if you missing the provider you need, please feel free to add it here
type FederatedIdentityProvider string

var (
	// FacebookIDProvider is a Facebook ID provider.
	FacebookIDProvider FederatedIdentityProvider = "FACEBOOK"
	// GoogleIDProvider is a Google ID provider.
	GoogleIDProvider FederatedIdentityProvider = "GOOGLE"
	// TwitterIDProvider is a Twitter ID provider.
	TwitterIDProvider FederatedIdentityProvider = "TWITTER"
)

// IsValid has to be called everywhere input happens, otherwise you risk to operate on bad data - no guarantees.
func (fid FederatedIdentityProvider) IsValid() bool {
	switch fid {
	case FacebookIDProvider, GoogleIDProvider, TwitterIDProvider:
		return true
	}
	return false
}
