package s3

import (
	"io/fs"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jszwec/s3fs"
	"github.com/madappgang/identifo/model"
)

func NewFS(settings model.FileStorageS3) (fs.FS, error) {
	session, err := NewSession(settings.Region)
	if err != nil {
		return nil, err
	}
	return s3fs.New(s3.New(session), settings.Bucket), nil
}

type fsWithPath struct{}

// name := filepath.Join(c.Root, filepath.Clean("/"+p)) // "/"+ for security. TODO: Jack: it is not clear why adding leading slash providing extra security
