package model

//FederatedIdentityProvider external federated identity provider type, if you missing the provider you need, please feel free to add it here
type FederatedIdentityProvider string

var (
	FacebookIDProvider FederatedIdentityProvider = "FACEBOOK"
	GoogleIDProvider   FederatedIdentityProvider = "GOOGLE"
	TwitterIDProvider  FederatedIdentityProvider = "TWITTER"
)

// IsValid has to be called everywhere input happens, or you risk bad data - no guarantees
func (fid FederatedIdentityProvider) IsValid() bool {
	switch fid {
	case FacebookIDProvider, GoogleIDProvider, TwitterIDProvider:
		return true
	}
	return false
}
