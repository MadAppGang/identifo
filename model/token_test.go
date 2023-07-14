package model_test

import (
	"encoding/json"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/madappgang/identifo/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClaimsMarshalJSON(t *testing.T) {
	cl := model.Claims{
		Payload: map[string]any{
			"key1": "value1",
			"key2": "value2",
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  "issuer",
			Subject: "subject",
		},
	}

	b, err := json.Marshal(cl)
	require.NoError(t, err)

	var rcm map[string]string
	err = json.Unmarshal(b, &rcm)
	require.NoError(t, err)

	assert.Len(t, rcm, 4)
	assert.Equal(t, "issuer", rcm["iss"])
	assert.Equal(t, "subject", rcm["sub"])
	assert.Equal(t, "value1", rcm["key1"])
	assert.Equal(t, "value2", rcm["key2"])

	cl.Type = model.TokenTypeActor.String()
	cl.ID = "12345"
	cl.KeyID = "12345"
	b, err = json.Marshal(cl)
	require.NoError(t, err)

	err = json.Unmarshal(b, &rcm)
	require.NoError(t, err)
	assert.Len(t, rcm, 7)
	assert.Equal(t, "12345", rcm["kid"])
	assert.Equal(t, "12345", rcm["jti"])
	assert.Equal(t, model.TokenTypeActor.String(), rcm["type"])
}

func TestClaimsUnmarshalJSON(t *testing.T) {
	data := `{"iss":"issuer","key1":"value1","key2":"value2","sub":"subject", "role:tenant1:group1": ["admin",  "guest"], "role:tenant1:group2": ["guest"] }`

	var cl model.Claims
	err := json.Unmarshal([]byte(data), &cl)
	require.NoError(t, err)

	assert.Len(t, cl.Payload, 4)
	assert.Equal(t, "issuer", cl.Issuer)
	assert.Equal(t, "subject", cl.Subject)
	assert.Equal(t, "value1", cl.Payload["key1"])
	assert.Equal(t, "value2", cl.Payload["key2"])
	assert.Contains(t, cl.Payload, "role:tenant1:group1")
	assert.Contains(t, cl.Payload["role:tenant1:group1"], "admin")
	assert.Contains(t, cl.Payload["role:tenant1:group1"], "guest")
	assert.Equal(t, cl.Payload["role:tenant1:group2"], []string{"guest"})
}

func TestTokenClaimMethod(t *testing.T) {
	exp := time.Now().Truncate(time.Second)
	cl := model.Claims{
		Payload: map[string]any{
			"key1": "value1",
			"key2": "value2",
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "issuer",
			Subject:   "subject",
			ExpiresAt: jwt.NewNumericDate(exp),
			ID:        "12345",
			IssuedAt:  jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(exp),
		},
		Type: model.TokenTypeID.String(),
	}
	token := model.TokenWithClaims(jwt.SigningMethodES256, "12345", cl)
	assert.Equal(t, "subject", token.UserID())
	assert.Equal(t, cl.Payload, token.Payload())
	assert.Equal(t, model.TokenTypeID, token.Type())
	assert.Equal(t, exp, token.ExpiresAt())
	assert.Equal(t, "12345", token.ID())
	assert.Equal(t, exp, token.IssuedAt())
	assert.Equal(t, "issuer", token.Issuer())
	assert.Equal(t, exp, token.NotBefore())
	assert.Equal(t, "subject", token.Subject())
}
