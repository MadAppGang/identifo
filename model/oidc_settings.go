package model

import "fmt"

type OIDCSettings struct {
	ProviderName     string   `bson:"provider_name,omitempty" json:"provider_name,omitempty"`
	ProviderURL      string   `bson:"provider_url,omitempty" json:"provider_url,omitempty"`
	Issuer           string   `bson:"issuer,omitempty" json:"issuer,omitempty"`
	ClientID         string   `bson:"client_id,omitempty" json:"client_id,omitempty"`
	ClientSecret     string   `bson:"client_secret,omitempty" json:"client_secret,omitempty"`
	EmailClaimField  string   `bson:"email_claim_field,omitempty" json:"email_claim_field,omitempty"`
	UserIDClaimField string   `bson:"user_id_claim_field,omitempty" json:"user_id_claim_field,omitempty"`
	Scopes           []string `bson:"scopes,omitempty" json:"scopes,omitempty"`
	InitURL          string   `bson:"init_url,omitempty" json:"init_url,omitempty"`
	// ScopeMapping maps OIDC scopes to Identifo scopes.
	ScopeMapping map[string]string `bson:"scope_mapping,omitempty" json:"scope_mapping,omitempty"`
}

func (s OIDCSettings) IsValid() error {
	if s.ProviderURL == "" {
		return fmt.Errorf("provider_url not specified for oidc")
	}
	if s.ClientID == "" {
		return fmt.Errorf("client_id not specified for oidc")
	}

	return nil
}
