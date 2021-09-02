package model

type ConnectionTester interface {
	Connect() error
}

type TestConnection struct {
	Type       TestType            `json:"type"`
	Database   *DatabaseSettings   `json:"database,omitempty"`
	KeyStorage *KeyStorageSettings `json:"key_storage,omitempty"`
}

// TestType is a test type
type TestType string

const (
	TTDatabase   TestType = "database"
	TTKeyStorage TestType = "key_sotrage"
)
