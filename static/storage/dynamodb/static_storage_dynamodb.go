package dynamodb

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/madappgang/identifo/model"
	staticStoreLocal "github.com/madappgang/identifo/static/storage/local"
	idynamodb "github.com/madappgang/identifo/storage/dynamodb"
)

const staticFilesTableName = "StaticFiles"

// StaticFilesStorage is a storage of static files in DynamoDB.
type StaticFilesStorage struct {
	db           *idynamodb.DB
	localStorage *staticStoreLocal.StaticFilesStorage
}

// NewStaticFilesStorage creates and returns new local static files storage.
func NewStaticFilesStorage(settings model.StaticFilesStorageSettings, localStorage *staticStoreLocal.StaticFilesStorage) (*StaticFilesStorage, error) {
	db, err := idynamodb.NewDB(settings.Endpoint, settings.Region)
	if err != nil {
		return nil, err
	}
	sfs := &StaticFilesStorage{
		db:           db,
		localStorage: localStorage,
	}
	if err = sfs.ensureTable(); err != nil {
		return nil, err
	}
	return sfs, nil
}

// ParseTemplate parses the html template.
func (sfs *StaticFilesStorage) ParseTemplate(templateName string) (*template.Template, error) {
	fileStr, err := sfs.getStaticFile(templateName)
	if err != nil {
		return nil, err
	}
	return template.New(templateName).Parse(fileStr)
}

// UploadTemplate is for html template uploads.
func (sfs *StaticFilesStorage) UploadTemplate(templateName string, contents []byte) error {
	return sfs.putStaticFile(templateName, contents)
}

// ReadAppleFile is for reading Apple-related static files.
func (sfs *StaticFilesStorage) ReadAppleFile(filename string) ([]byte, error) {
	// Check if file exists. If not - return nil error and nil slice.
	fileStr, err := sfs.getStaticFile(filename)
	if err != nil {
		if err == model.ErrorNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while checking filename '%s' existence. %s", filename, err)
	}
	return []byte(fileStr), nil
}

// UploadAppleFile is for Apple-related file uploads.
func (sfs *StaticFilesStorage) UploadAppleFile(filename string, contents []byte) error {
	return sfs.putStaticFile(filename, contents)
}

// AssetHandlers returns handlers for assets.
func (sfs *StaticFilesStorage) AssetHandlers() *model.AssetHandlers {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		split := strings.Split(r.URL.Path, "/")
		if len(split) == 0 || len(split[0]) == 0 {
			err := fmt.Errorf("Empty file name")
			writeError(w, err, http.StatusNotFound, err.Error())
			return
		}
		name := split[0]

		fileStr, err := sfs.getStaticFile(name)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(name)))
		if _, err = w.Write([]byte(fileStr)); err != nil {
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

type fileData struct {
	Name string `json:"name"`
	File string `json:"file"`
}

func (sfs *StaticFilesStorage) getStaticFile(name string) (string, error) {
	if len(name) == 0 {
		return "", model.ErrorWrongDataFormat
	}

	result, err := sfs.db.C.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(staticFilesTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(name),
			},
		},
	})

	if err != nil {
		err = fmt.Errorf("Error getting static file from db: %s", err)
		return "", err
	}
	if result.Item == nil {
		return "", fmt.Errorf("%s not found in %s table", name, staticFilesTableName)
	}

	fd := new(fileData)
	if err = dynamodbattribute.UnmarshalMap(result.Item, fd); err != nil {
		err = fmt.Errorf("Error unmarshalling static file data: %s", err)
		return "", err
	}
	return fd.File, nil
}

func (sfs *StaticFilesStorage) putStaticFile(name string, contents []byte) error {
	if len(name) == 0 {
		return model.ErrorWrongDataFormat
	}

	f := &fileData{Name: name, File: string(contents)}
	marshalled, err := dynamodbattribute.MarshalMap(f)
	if err != nil {
		return fmt.Errorf("Error marshalling static file: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      marshalled,
		TableName: aws.String(staticFilesTableName),
	}

	if _, err = sfs.db.C.PutItem(input); err != nil {
		return fmt.Errorf("Error putting static file to storage: %s", err)
	}
	return nil
}

func (sfs *StaticFilesStorage) ensureTable() error {
	exists, err := sfs.db.IsTableExists(staticFilesTableName)
	if err != nil {
		return fmt.Errorf("Error checking static files table existence: %s", err)
	}
	if exists {
		return nil
	}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
		TableName:   aws.String(staticFilesTableName),
	}

	_, err = sfs.db.C.CreateTable(input)
	return err
}

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
