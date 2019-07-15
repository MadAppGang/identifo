package s3

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	awsAccessKeyEnvName       = "AWS_ACCESS_KEY"
	awsSecretAccessKeyEnvName = "AWS_SECRET_ACCESS_KEY"
)

// NewS3Client creates and returns new S3 client.
func NewS3Client(region string) (*s3.S3, error) {
	if len(region) == 0 {
		return nil, fmt.Errorf("No S3 region for configuration specified")
	}

	cfg := getConfig(region)
	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, fmt.Errorf("Cannot create new session: %s", err)
	}
	return s3.New(sess, cfg), nil
}

func getConfig(region string) *aws.Config {
	cfg := aws.NewConfig().WithRegion(region)

	cfg.HTTPClient = http.DefaultClient
	cfg.HTTPClient.Timeout = 10 * time.Second

	return cfg
}
