package api_test

import (
	"os"
	"testing"
	"time"

	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionStore(t *testing.T) {
	if os.Getenv("IDENTIFO_STORAGE_MONGO_TEST_INTEGRATION") == "" {
		t.SkipNow()
	}

	rh := os.Getenv("IDENTIFO_REDIS_HOST")

	s := model.RedisDatabaseSettings{
		Address: rh,
		Prefix:  "identifo",
	}
	storage, err := redis.NewSessionStorage(logging.DefaultLogger, s)
	require.NoError(t, err)

	expDate := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
	testSession := model.Session{
		ID:             "abc",
		ExpirationTime: expDate.Unix(),
	}

	err = storage.InsertSession(testSession)
	require.NoError(t, err)

	ses, err := storage.GetSession("abc")
	require.NoError(t, err)

	assert.Equal(t, testSession, ses)
}
