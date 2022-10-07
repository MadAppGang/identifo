package boltdb_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/boltdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const dbpath = "../../data/boltdb-test.db"

var testApp = model.AppData{
	ID:           "aabbccddeeff",
	Name:         "Test App",
	Secret:       "secret",
	Active:       true,
	RedirectURLs: []string{"http://localhost:44000"},
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	os.MkdirAll(filepath.Dir(dbpath), os.ModePerm)
}

func shutdown() {
	os.Remove(dbpath)
}

func TestBoltDBAppCreateApp(t *testing.T) {
	sts := model.BoltDBDatabaseSettings{
		Path: dbpath,
	}
	apps, err := boltdb.NewAppStorage(sts)
	defer apps.Close()
	require.NoError(t, err)

	a, err := apps.CreateApp(testApp)
	require.NoError(t, err)

	assert.Equal(t, a.ID, testApp.ID)

	testApp2 := testApp
	testApp2.ID = ""

	a, err = apps.CreateApp(testApp2)
	require.NoError(t, err)

	assert.NotEqual(t, a.ID, testApp.ID)
	assert.NotEmpty(t, a.ID)
}

func TestBoltDBAppFindAppById(t *testing.T) {
	sts := model.BoltDBDatabaseSettings{
		Path: dbpath,
	}
	apps, err := boltdb.NewAppStorage(sts)
	defer apps.Close()
	require.NoError(t, err)

	testApp2 := testApp
	testApp2.ID = ""

	a, err := apps.CreateApp(testApp2)
	require.NoError(t, err)
	assert.NotEmpty(t, a.ID)

	a2, err := apps.ActiveAppByID(a.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, a2.ID)
}

func TestBoltDBAppFindAppFetchApps(t *testing.T) {
	shutdown()

	sts := model.BoltDBDatabaseSettings{
		Path: dbpath,
	}
	apps, err := boltdb.NewAppStorage(sts)
	defer apps.Close()
	require.NoError(t, err)

	a, err := apps.CreateApp(testApp)
	require.NoError(t, err)

	assert.Equal(t, a.ID, testApp.ID)

	testApp2 := testApp
	testApp2.ID = ""

	a, err = apps.CreateApp(testApp2)
	require.NoError(t, err)
	assert.NotEmpty(t, a.ID)

	a2, err := apps.ActiveAppByID(a.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, a2.ID)

	allApps, err := apps.FetchApps("")
	require.NoError(t, err)
	assert.Len(t, allApps, 2)
}
