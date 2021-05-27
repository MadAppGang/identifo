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

// GetFile is a generic method for fetching a file by name from DynamoDB.
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

type fileData struct {
	Name string `json:"name"`
	File string `json:"file"`
}

func (sfs *StaticFilesStorage) getFile(name string) ([]byte, error) {
	if len(name) == 0 {
		return nil, model.ErrorWrongDataFormat
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
		return nil, fmt.Errorf("Error getting static file from db: %s", err)
	}
	if result.Item == nil {
		return nil, model.ErrorNotFound
	}

	fd := new(fileData)
	if err = dynamodbattribute.UnmarshalMap(result.Item, fd); err != nil {
		return nil, fmt.Errorf("Error unmarshalling static file data: %s", err)
	}
	return []byte(fd.File), nil
}

// UploadFile is a generic file uploader.
func (sfs *StaticFilesStorage) UploadFile(name string, contents []byte) error {
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
func (sfs *StaticFilesStorage) GetAppleFile(name string) ([]byte, error) {
	// Call private method since we don't want to retry fetching file from the local storage.
	// If error is not nil and not model.ErrorNotFound, we'll retry the whole GetAppleFile.
	file, err := sfs.getFile(name)
	if err == nil {
		return file, nil
	}

	if err == model.ErrorNotFound {
		return nil, nil
	}

	log.Printf("Error getting %s from DynamoDB: %s. Using local storage.\n", name, err)
	return sfs.localStorage.GetAppleFile(name)
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

		// Call private method since local handlers, if needed, will be defined later.
		file, err := sfs.getFile(name)
		if err == nil {
			w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(name)))
			if _, err = w.Write(file); err != nil {
				log.Printf("Error writing body to the response: %s\n", err)
			}
			return
		}

		prefix := strings.TrimSuffix(r.URL.Path, name)
		localHandler := http.FileServer(http.Dir(path.Join(sfs.localStorage.Folder, prefix)))
		http.StripPrefix(prefix, localHandler).ServeHTTP(w, r)
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

// WebHandlers returns handlers for the web.
// Web build is always being stored locally, despite the static storage type.
func (sfs *StaticFilesStorage) WebHandlers() *model.WebHandlers {
	return sfs.localStorage.WebHandlers()
}

// Close is to satisfy the interface.
func (sfs *StaticFilesStorage) Close() {}

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
