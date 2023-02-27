package mongo_test

import (
	"os"
	"testing"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppCRUD(t *testing.T) {
	if os.Getenv("IDENTIFO_STORAGE_MONGO_TEST_INTEGRATION") == "" {
		t.SkipNow()
	}

	connStr := os.Getenv("IDENTIFO_STORAGE_MONGO_CONN")

	s, err := mongo.NewAppStorage(model.MongoDatabaseSettings{
		ConnectionString: connStr,
		DatabaseName:     "test_users",
	})
	require.NoError(t, err)

	expectedApp := model.AppData{
		ID:     "test_app",
		Active: true,
		Secret: "test_secret",
		OIDCSettings: model.OIDCSettings{
			ScopeMapping: map[string]string{
				"offline_access": "offline",
			},
		},
	}

	app, err := s.CreateApp(expectedApp)
	require.NoError(t, err)

	defer s.DeleteApp(app.ID)

	_, err = s.AppByID(app.ID)
	require.NoError(t, err)

	assert.Equal(t, app.ID, expectedApp.ID)
	assert.Equal(t, app.Active, expectedApp.Active)
	assert.Equal(t, app.Secret, expectedApp.Secret)
	assert.Equal(t, app.OIDCSettings.ScopeMapping, expectedApp.OIDCSettings.ScopeMapping)
}
