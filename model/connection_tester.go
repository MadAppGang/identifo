package model

type ConnectionTester interface {
	Connect() error
}

type TestConnection struct {
	Type        TestType                 `json:"type"`
	Database    *DatabaseSettings        `json:"database,omitempty"`
	KeyStorage  *KeyStorageSettings      `json:"key_storage,omitempty"`
	FileStorage *FileStorageTestSettings `json:"file_storage,omitempty"`
}

type FileStorageTestSettings struct {
	ExpectedFiles []string             `json:"expected_file,omitempty"`
	FileStorage   *FileStorageSettings `json:"file_storage,omitempty"`
}

// TestType is a test type
type TestType string

const (
	TTDatabase    TestType = "database"
	TTKeyStorage  TestType = "key_storage"
	TTFileStorage TestType = "file_storage"
)
