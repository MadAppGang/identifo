package s3

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	s3Storage "github.com/madappgang/identifo/external_services/storage/s3"
	"github.com/madappgang/identifo/model"
	staticStoreLocal "github.com/madappgang/identifo/static/storage/local"
)

// StaticFilesStorage is a storage of static files in S3.
type StaticFilesStorage struct {
	client       *s3.S3
	bucket       string
	folder       string
	localStorage *staticStoreLocal.StaticFilesStorage
}

// NewStaticFilesStorage creates and returns new static files storage in S3.
func NewStaticFilesStorage(settings model.StaticFilesStorageSettings, localStorage *staticStoreLocal.StaticFilesStorage) (*StaticFilesStorage, error) {
	s3Client, err := s3Storage.NewS3Client(settings.Region)
	if err != nil {
		return nil, err
	}

	return &StaticFilesStorage{
		client:       s3Client,
		bucket:       settings.Bucket,
		folder:       settings.Folder,
		localStorage: localStorage,
	}, nil
}

// GetFile is for fetching a file by name from the S3 bucket.
// It is a wrapper over the private method getFile.
// It provides fallback behavior via using eponymous local storage method.
func (sfs *StaticFilesStorage) GetFile(name string) ([]byte, error) {
	file, err := sfs.getFile(name)
	if err == nil {
		return file, nil
	}

	log.Printf("Error getting %s from DynamoDB: %s. Using local storage.\n", name, err)
	return sfs.localStorage.GetFile(name)
}

func (sfs *StaticFilesStorage) getFile(name string) ([]byte, error) {
	key, err := model.GetStaticFilePathByFilename(name, sfs.folder)
	if err != nil {
		return nil, fmt.Errorf("Cannot get file %s. %s", key, err)
	}

	getFileInput := &s3.GetObjectInput{
		Bucket: aws.String(sfs.bucket),
		Key:    aws.String(key),
	}

	resp, err := sfs.client.GetObject(getFileInput)
	if err != nil {
		return nil, fmt.Errorf("Cannot get %s from S3: %s", key, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Cannot read response from S3: %s", err)
	}
	if len(body) == 0 {
		return nil, model.ErrorNotFound
	}
	return body, nil
}

// UploadFile is a generic file uploader.
func (sfs *StaticFilesStorage) UploadFile(name string, contents []byte) error {
	filepath, err := model.GetStaticFilePathByFilename(name, sfs.folder)
	if err != nil {
		return fmt.Errorf("Cannot compose filepath. %s", err)
	}
	_, err = sfs.client.PutObject(&s3.PutObjectInput{
		Bucket:       aws.String(sfs.bucket),
		Key:          aws.String(filepath),
		ACL:          aws.String("private"),
		StorageClass: aws.String(s3.ObjectStorageClassStandard),
		Body:         bytes.NewReader(contents),
		ContentType:  aws.String(mime.TypeByExtension(path.Ext(name))),
	})
	if err == nil {
		log.Printf("Successfully put %s to S3\n", filepath)
	}
	return nil
}

// ParseTemplate parses the html template.
func (sfs *StaticFilesStorage) ParseTemplate(templateName string) (*template.Template, error) {
	tmplBytes, err := sfs.GetFile(templateName)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(templateName).Parse(string(tmplBytes))
	if err != nil {
		return nil, fmt.Errorf("Cannot parse template '%s'. %s", templateName, err)
	}
	return tmpl, nil
}

// GetAppleFile is for reading Apple-related static files.
// Unlike generic GetFile, it does not treat model.ErrorNotFound as error.
func (sfs *StaticFilesStorage) GetAppleFile(filename string) ([]byte, error) {
	// Call private method since we don't want to retry fetching file from the local storage.
	// If error is not nil and not model.ErrorNotFound, we'll retry the whole GetAppleFile.
	file, err := sfs.getFile(filename)
	if err == nil {
		return file, nil
	}

	if err == model.ErrorNotFound {
		return nil, nil
	}

	log.Printf("Error getting %s from S3: %s. Using local storage.\n", filename, err)
	return sfs.localStorage.GetAppleFile(filename)
}

// AssetHandlers returns handlers for assets.
func (sfs *StaticFilesStorage) AssetHandlers() *model.AssetHandlers {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		split := strings.Split(r.URL.Path, "/")

		lensplit := len(split)
		if lensplit == 0 || len(split[lensplit-1]) == 0 {
			err := fmt.Errorf("Empty file name")
			writeError(w, err, http.StatusNotFound, err.Error())
			return
		}
		name := split[lensplit-1]

		file, err := sfs.GetFile(name)
		if err != nil {
			if err == model.ErrorNotFound {
				writeError(w, fmt.Errorf("%s not found", name), http.StatusNotFound, "")
				return
			}
			writeError(w, err, http.StatusInternalServerError, "")
			return
		}

		w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(name)))
		if _, err = w.Write(file); err != nil {
			log.Printf("Error writing body to the response: %s\n", err)
			return
		}
	})

	return &model.AssetHandlers{
		StylesHandler:  handler,
		ScriptsHandler: handler,
		ImagesHandler:  handler,
		FontsHandler:   handler,
	}
}

// AdminPanelHandlers returns handlers for the admin panel.
// Adminpanel build is always being stored locally, despite the static storage type.
func (sfs *StaticFilesStorage) AdminPanelHandlers() *model.AdminPanelHandlers {
	return sfs.localStorage.AdminPanelHandlers()
}

// Close is to satisfy the interface.
func (sfs *StaticFilesStorage) Close() {}

// writeError writes an error message to the response and logger.
func writeError(w http.ResponseWriter, err error, code int, userInfo string) {
	log.Printf("http error: %s (code=%d)\n", err, code)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	responseString := `
	<!DOCTYPE html>
	<html>
	<head>
	  <title>Home Network</title>
	</head>
	<body>
	<h2>Error</h2></br>
	<h3>
	` +
		fmt.Sprintf("Error: %s, code: %d, userInfo: %s", err.Error(), code, userInfo) +
		`
	</h3>
	</body>
	</html>
	`
	w.WriteHeader(code)
	if _, wrErr := io.WriteString(w, responseString); wrErr != nil {
		log.Println("Error writing response string:", wrErr)
	}
}
