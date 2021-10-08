package s3

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// NewS3Client creates and returns new S3 client.
func NewS3Client(region string) (*s3.S3, error) {
	cfg := getConfig(region)
	sess, err := NewSession(region)
	if err != nil {
		return nil, err
	}
	return s3.New(sess, cfg), nil
}

func NewSession(region string) (*session.Session, error) {
	cfg := getConfig(region)
	sess, err := session.NewSession(cfg.WithCredentialsChainVerboseErrors(true))
	if err != nil {
		return nil, fmt.Errorf("error creating new s3 session: %s", err)
	}
	return sess, err
}

func getConfig(region string) *aws.Config {
	cfg := aws.NewConfig()
	if len(region) > 0 {
		cfg = cfg.WithRegion(region)
	}

	cfg.HTTPClient = http.DefaultClient
	cfg.HTTPClient.Timeout = 10 * time.Second

	return cfg
}
