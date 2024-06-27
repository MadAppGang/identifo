package api_test

import (
	"os"
	"testing"

	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFetchUser(t *testing.T) {
	if os.Getenv("IDENTIFO_STORAGE_MONGO_TEST_INTEGRATION") == "" {
		t.SkipNow()
	}

	connStr := os.Getenv("IDENTIFO_STORAGE_MONGO_CONN")

	s, err := mongo.NewUserStorage(
		logging.DefaultLogger,
		model.MongoDatabaseSettings{
			ConnectionString: connStr,
			DatabaseName:     "test_users",
		})
	require.NoError(t, err)

	s.(*mongo.UserStorage).ClearAllUserData()
	testUser := model.User{
		Username:   "test",
		Email:      "test@examplle.com",
		Phone:      "+71111111111",
		FullName:   "test",
		Scopes:     []string{"test"},
		AccessRole: "test",
		Active:     true,
	}

	u, err := s.AddUserWithPassword(testUser, "test", "test", false)
	require.NoError(t, err)

	assert.NotEmpty(t, u.ID)
	assert.NotEmpty(t, u.Pswd)

	testUser.ID = u.ID
	testUser.Pswd = u.Pswd

	assert.Equal(t, testUser, u)

	t.Run("fetch user by id", func(t *testing.T) {
		tu, err := s.UserByID(u.ID)
		require.NoError(t, err)
		assert.Equal(t, u, tu)
	})

	t.Run("fetch user by email", func(t *testing.T) {
		tu, err := s.UserByEmail(testUser.Email)
		require.NoError(t, err)

		assert.Empty(t, tu.Pswd)

		tu.Pswd = u.Pswd
		assert.Equal(t, u, tu)
	})

	t.Run("fetch user by phone", func(t *testing.T) {
		tu, err := s.UserByPhone(testUser.Phone)
		require.NoError(t, err)

		assert.Empty(t, tu.Pswd)

		tu.Pswd = u.Pswd
		assert.Equal(t, u, tu)
	})

	t.Run("user already exists", func(t *testing.T) {
		_, err := s.AddUserWithPassword(testUser, "test", "test", false)
		require.ErrorIs(t, err, model.ErrorUserExists)
	})

	t.Run("user not found", func(t *testing.T) {
		_, err := s.UserByID(primitive.NewObjectID().Hex())
		require.ErrorIs(t, err, model.ErrUserNotFound)

		_, err = s.UserByPhone("+71111111112")
		require.ErrorIs(t, err, model.ErrUserNotFound)

		_, err = s.UserByEmail("noemail@example.com")
		require.ErrorIs(t, err, model.ErrUserNotFound)
	})
}
