package model

import "context"

// Storage represents any storage type.
type Storage interface {
	// Ready returns nil if the storage is ready and connected.
	Ready(ctx context.Context) error
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}

// ImportableStorage is a storage that could import raw data in JSON format.
type ImportableStorage interface {
	ImportJSON(data []byte, clearOldData bool) error
}

type UserImportData struct {
	Users []User     `json:"users"`
	Data  []UserData `json:"data"`
}
