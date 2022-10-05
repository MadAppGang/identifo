package config

import (
	"testing"

	"github.com/madappgang/identifo/v2/model"
	"github.com/stretchr/testify/require"
)

func TestInitConfigurationStorage(t *testing.T) {
	config := model.ConfigStorageSettings{
		Type:      model.ConfigStorageTypeS3,
		RawString: "s3://region@bucket/folder/file.ext",
		S3: &model.S3StorageSettings{
			Region: "region",
			Bucket: "bucket",
			Key:    "folder/file.ext",
		},
	}

	_, err := InitConfigurationStorage(config)
	require.NoError(t, err)
}
