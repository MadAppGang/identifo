package model_test

import (
	"encoding/json"
	"testing"

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
}

func TestClaimsUnmarshalJSON(t *testing.T) {
	data := `{"iss":"issuer","key1":"value1","key2":"value2","sub":"subject"}`

	var cl model.Claims
	err := json.Unmarshal([]byte(data), &cl)
	require.NoError(t, err)

	assert.Len(t, cl.Payload, 2)
	assert.Equal(t, "issuer", cl.Issuer)
	assert.Equal(t, "subject", cl.Subject)
	assert.Equal(t, "value1", cl.Payload["key1"])
	assert.Equal(t, "value2", cl.Payload["key2"])
}
