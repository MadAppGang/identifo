package config

import (
	"testing"

	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/stretchr/testify/require"
)

func TestInitConfigurationStorage(t *testing.T) {
	//"s3://region@bucket/folder/file.ext"
	config := model.FileStorageSettings{
		Type: model.FileStorageTypeS3,
		S3: model.FileStorageS3{
			Region: "region",
			Bucket: "bucket",
			Key:    "folder/file.ext",
		},
	}

	_, err := InitConfigurationStorage(logging.DefaultLogger, config)
	require.NoError(t, err)
}
