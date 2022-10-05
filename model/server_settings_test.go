package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigStorageSettingsFromStringS3HappyDay(t *testing.T) {
	flag := "s3://ap-southwest-2@my-bucket1t/dev/config.yaml"
	config, err := ConfigStorageSettingsFromStringS3(flag)

	require.NoError(t, err)
	assert.Equal(t, config.Type, ConfigStorageTypeS3)
	assert.Equal(t, config.S3.Region, "ap-southwest-2")
	assert.Equal(t, config.S3.Bucket, "my-bucket1t")
	assert.Equal(t, config.S3.Key, "/dev/config.yaml")
	assert.Empty(t, config.S3.Endpoint)
	assert.Equal(t, config.RawString, flag)
}

func TestConfigStorageSettingsFromStringS3EndpointEmpty(t *testing.T) {
	flag := "s3://ap-southwest-2@my-bucket1t/dev/config.yaml|"
	config, err := ConfigStorageSettingsFromStringS3(flag)

	require.NoError(t, err)
	assert.Equal(t, config.Type, ConfigStorageTypeS3)
	assert.Equal(t, config.S3.Region, "ap-southwest-2")
	assert.Equal(t, config.S3.Bucket, "my-bucket1t")
	assert.Equal(t, config.S3.Key, "/dev/config.yaml")
	assert.Empty(t, config.S3.Endpoint)
	assert.Equal(t, config.RawString, flag)
}

func TestConfigStorageSettingsFromStringS3Endpoint(t *testing.T) {
	flag := "s3://ap-southwest-2@my-bucket1t/dev/config.yaml|http://localhost:2020"
	config, err := ConfigStorageSettingsFromStringS3(flag)

	require.NoError(t, err)
	assert.Equal(t, config.Type, ConfigStorageTypeS3)
	assert.Equal(t, config.S3.Region, "ap-southwest-2")
	assert.Equal(t, config.S3.Bucket, "my-bucket1t")
	assert.Equal(t, config.S3.Key, "/dev/config.yaml")
	assert.Equal(t, config.S3.Endpoint, "http://localhost:2020")
	assert.Equal(t, config.RawString, flag)
}

func TestConfigStorageSettingsFromStringS3Wrong(t *testing.T) {
	_, err := ConfigStorageSettingsFromStringS3("s33://ap-southwest-2@my-bucket1t/dev/config.yaml|http://localhost:2020")
	require.Error(t, err)

	_, err = ConfigStorageSettingsFromStringS3("ap-southwest-2@my-bucket1t/dev/config.yaml|http://localhost:2020")
	require.Error(t, err)

	_, err = ConfigStorageSettingsFromStringS3("s3://ap-southwest-2@m$y-bucket1t/dev/config.yaml|http://localhost:2020")
	require.Error(t, err)
}

func TestConfigStorageSettingsFromStringFileHappyDay(t *testing.T) {
	flag := "myfolder/subfolder/config.yaml"
	config, err := ConfigStorageSettingsFromStringFile(flag)

	require.NoError(t, err)
	assert.Equal(t, config.Type, ConfigStorageTypeFile)
	assert.Equal(t, config.RawString, flag)
	assert.Empty(t, config.S3)
	assert.Equal(t, config.File.FileName, flag)
}

func TestConfigStorageSettingsFromStringFileHappyDayPrefix(t *testing.T) {
	flag := "file://myfolder/subfolder/config.yaml"
	config, err := ConfigStorageSettingsFromStringFile(flag)

	require.NoError(t, err)
	assert.Equal(t, config.Type, ConfigStorageTypeFile)
	assert.Equal(t, config.RawString, flag)
	assert.Empty(t, config.S3)
	assert.Equal(t, config.File.FileName, "myfolder/subfolder/config.yaml")
}

func TestConfigStorageSettingsFromStringFileHappyDayPrefixCase(t *testing.T) {
	flag := "FiLe://myfolder/subfolder/config.yaml"
	config, err := ConfigStorageSettingsFromStringFile(flag)

	require.NoError(t, err)
	assert.Equal(t, config.Type, ConfigStorageTypeFile)
	assert.Equal(t, config.RawString, flag)
	assert.Empty(t, config.S3)
	assert.Equal(t, config.File.FileName, "myfolder/subfolder/config.yaml")
}

func TestConfigStorageSettingsFromStringFileHappyDayRoot(t *testing.T) {
	flag := "/myfolder/subfolder/config.yaml"
	config, err := ConfigStorageSettingsFromStringFile(flag)

	require.NoError(t, err)
	assert.Equal(t, config.Type, ConfigStorageTypeFile)
	assert.Equal(t, config.RawString, flag)
	assert.Empty(t, config.S3)
	assert.Equal(t, config.File.FileName, "/myfolder/subfolder/config.yaml")
}

func TestConfigStorageSettingsFromString(t *testing.T) {
	flag := "s3://ap-southwest-2@my-bucket1t/dev/config.yaml|http://localhost:2020"
	config, err := ConfigStorageSettingsFromString(flag)

	require.NoError(t, err)
	assert.Equal(t, config.Type, ConfigStorageTypeS3)
	assert.Equal(t, config.S3.Region, "ap-southwest-2")
	assert.Equal(t, config.S3.Bucket, "my-bucket1t")
	assert.Equal(t, config.S3.Key, "/dev/config.yaml")
	assert.Equal(t, config.S3.Endpoint, "http://localhost:2020")
	assert.Equal(t, config.RawString, flag)
}

func TestConfigStorageSettingsFromStringFile(t *testing.T) {
	flag := "FiLe://myfolder/subfolder/config.yaml"
	config, err := ConfigStorageSettingsFromString(flag)

	require.NoError(t, err)
	assert.Equal(t, config.Type, ConfigStorageTypeFile)
	assert.Equal(t, config.RawString, flag)
	assert.Empty(t, config.S3)
	assert.Equal(t, config.File.FileName, "myfolder/subfolder/config.yaml")
}
