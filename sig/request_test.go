package sig_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/madappgang/identifo/sig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifySignature(t *testing.T) {
	secret := []byte("the most secret secret")
	body := "my request"
	request, _ := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		"http://google.com/whatever",
		strings.NewReader(body),
	)
	request.Header.Add(sig.ContentTypeHeaderKey, "application/json")
	bodyMD5 := sig.GetMD5([]byte(body))
	err := sig.AddHeadersAndSignRequest(request, secret, bodyMD5)

	require.NoError(t, err)

	assert.Equal(t, bodyMD5, request.Header["Content-Md5"][0])
	assert.NotEmpty(t, request.Header["Expires"][0])
	assert.NotEmpty(t, request.Header["Date"][0])
	assert.NotEmpty(t, request.Header["Digest"][0])

	err = sig.VerifySignature(request, secret)
	assert.NoError(t, err)
}
