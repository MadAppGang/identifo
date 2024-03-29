package s3

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// NewS3Client creates and returns new S3 client.
func NewS3Client(region, endpoint string) (*s3.S3, error) {
	cfg := getConfig(region, endpoint)

	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating new s3 session: %s", err)
	}

	return s3.New(sess, cfg), nil
}

func NewSession(region, endpoint string) (*session.Session, error) {
	cfg := getConfig(region, endpoint)
	sess, err := session.NewSession(cfg.WithCredentialsChainVerboseErrors(true))
	if err != nil {
		return nil, fmt.Errorf("error creating new s3 session: %s", err)
	}
	return sess, err
}

func getConfig(region, endpoint string) *aws.Config {
	cfg := aws.NewConfig().
		WithHTTPClient(&http.Client{
			Timeout: 10 * time.Second,
		}).
		WithCredentialsChainVerboseErrors(true)

	// critically important for local tests, as we could not create localhost subdomains
	// https://docs.aws.amazon.com/AmazonS3/latest/userguide/VirtualHosting.html
	if len(os.Getenv("IDENTIFO_FORCE_S3_PATH_STYLE")) > 0 {
		cfg.WithS3ForcePathStyle(true)
	}

	if len(endpoint) > 0 {
		cfg.WithEndpoint(endpoint)
	}

	if len(region) > 0 {
		cfg.WithRegion(region)
	}

	return cfg
}
