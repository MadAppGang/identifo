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
)

// StaticFilesStorage is a storage of static files in S3.
type StaticFilesStorage struct {
	client             *s3.S3
	bucket             string
	region             string
	pagesPath          string
	emailTemplatesPath string
	appleFilesPath     string
}

// NewStaticFilesStorage creates and returns new static files storage in S3.
func NewStaticFilesStorage(settings model.StaticFilesStorageSettings) (*StaticFilesStorage, error) {
	s3Client, err := s3Storage.NewS3Client(settings.Region)
	if err != nil {
		return nil, err
	}

	return &StaticFilesStorage{
		client:             s3Client,
		bucket:             settings.StaticFilesLocation,
		region:             settings.Region,
		pagesPath:          settings.PagesPath,
		emailTemplatesPath: settings.EmailTemplatesPath,
		appleFilesPath:     settings.AppleFilesPath,
	}, nil
}

// ParseTemplate parses the html template.
func (sfs *StaticFilesStorage) ParseTemplate(templateName string) (*template.Template, error) {
	if strings.Contains(strings.ToLower(templateName), "email") {
		templateName = path.Join(sfs.emailTemplatesPath, templateName)
	} else {
		templateName = path.Join(sfs.pagesPath, templateName)
	}

	getTemplateInput := &s3.GetObjectInput{
		Bucket: aws.String(sfs.bucket),
		Key:    aws.String(templateName),
	}

	resp, err := sfs.client.GetObject(getTemplateInput)
	if err != nil {
		return nil, fmt.Errorf("Cannot get %s from S3: %s", templateName, err)
	}
	defer resp.Body.Close()

	tmpl, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Cannot decode S3 response: %s", err)
	}

	return template.New(templateName).Parse(string(tmpl)) // TODO: test
}

// UploadTemplate is for html template uploads.
func (sfs *StaticFilesStorage) UploadTemplate(templateName string, contents []byte) error {
	if strings.Contains(strings.ToLower(templateName), "email") {
		templateName = path.Join(sfs.emailTemplatesPath, templateName)
	} else {
		templateName = path.Join(sfs.pagesPath, templateName)
	}

	_, err := sfs.client.PutObject(&s3.PutObjectInput{
		Bucket:       aws.String(sfs.bucket),
		Key:          aws.String(templateName),
		ACL:          aws.String("private"),
		StorageClass: aws.String(s3.ObjectStorageClassStandard),
		Body:         bytes.NewReader(contents),
		ContentType:  aws.String("text/html"),
	})
	if err == nil {
		log.Printf("Successfully put %s to S3\n", templateName)
	}
	return nil
}

// ReadAppleFile is for reading Apple-related static files.
func (sfs *StaticFilesStorage) ReadAppleFile(filename string) ([]byte, error) {
	getFileInput := &s3.GetObjectInput{
		Bucket: aws.String(sfs.bucket),
		Key:    aws.String(path.Join(sfs.appleFilesPath, filename)),
	}

	resp, err := sfs.client.GetObject(getFileInput)
	if err != nil {
		return nil, fmt.Errorf("Cannot get %s from S3: %s", filename, err)
	}
	defer resp.Body.Close()

	file, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Cannot decode S3 response: %s", err)
	}
	if len(file) == 0 {
		return nil, nil
	}
	return file, nil
}

// UploadAppleFile is for Apple-related file uploads.
func (sfs *StaticFilesStorage) UploadAppleFile(filename string, contents []byte) error {
	filename = path.Join(sfs.appleFilesPath, filename)
	_, err := sfs.client.PutObject(&s3.PutObjectInput{
		Bucket:       aws.String(sfs.bucket),
		Key:          aws.String(filename),
		ACL:          aws.String("private"),
		StorageClass: aws.String(s3.ObjectStorageClassStandard),
		Body:         bytes.NewReader(contents),
	})
	if err == nil {
		log.Printf("Successfully put %s to S3\n", filename)
	}
	return nil
}

// AssetHandlers returns handlers for assets.
func (sfs *StaticFilesStorage) AssetHandlers() *model.AssetHandlers {
	folder := strings.TrimLeft(strings.TrimSuffix(sfs.pagesPath, "/html"), ".")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := strings.Join([]string{folder, r.URL.Path}, "")

		getStyleInput := &s3.GetObjectInput{
			Bucket: aws.String(sfs.bucket),
			Key:    aws.String(key),
		}

		resp, err := sfs.client.GetObject(getStyleInput)
		if err != nil {
			err = fmt.Errorf("Cannot get %s from S3: %s", r.URL.Path, err)
			writeError(w, err, http.StatusInternalServerError, "")
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			err = fmt.Errorf("Cannot read response from S3: %s", err)
			writeError(w, err, http.StatusInternalServerError, "")
			return
		}

		w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(key)))
		if _, err = w.Write(body); err != nil {
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
	srcHandler := http.StripPrefix("/src/", http.FileServer(http.Dir(path.Join(model.AdminPanelBuildPath, "/src"))))
	managementHandleFunc := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(model.AdminPanelBuildPath, "/index.html"))
	}
	buildHandler := http.FileServer(http.Dir(model.AdminPanelBuildPath))

	return &model.AdminPanelHandlers{
		SrcHandler:        srcHandler,
		ManagementHandler: http.HandlerFunc(managementHandleFunc),
		BuildHandler:      buildHandler,
	}
}

// Close is to satisfy the interface.
func (sfs *StaticFilesStorage) Close() {}

// writeError writes an error message to the response and logger.
func writeError(w http.ResponseWriter, err error, code int, userInfo string) {
	log.Printf("http error: %s (code=%d)\n", err, code)

	// Hide error from client if it is internal.
	if code == http.StatusInternalServerError {
		err = model.ErrorInternal
	}

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
