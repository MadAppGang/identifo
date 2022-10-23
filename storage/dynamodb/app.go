package dynamodb

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/madappgang/identifo/v2/model"
	"github.com/rs/xid"
)

var appsTableName = "Applications"

// NewAppStorage creates new DynamoDB AppStorage implementation.
func NewAppStorage(settings model.DynamoDatabaseSettings) (model.AppStorage, error) {
	if len(settings.Endpoint) == 0 || len(settings.Region) == 0 {
		return nil, ErrorEmptyEndpointRegion
	}

	// create database
	db, err := NewDB(settings.Endpoint, settings.Region)
	if err != nil {
		return nil, err
	}

	as := &AppStorage{db: db}
	err = as.ensureTable()
	return as, err
}

// AppStorage a is fully functional app storage.
type AppStorage struct {
	db *DB
}

// ensureTable ensures that app table exists in database.
func (as *AppStorage) ensureTable() error {
	exists, err := as.db.IsTableExists(appsTableName)
	if err != nil {
		log.Println("Error checking Applications table existence:", err)
		return err
	}
	if exists {
		return nil
	}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
		TableName:   aws.String(appsTableName),
	}

	_, err = as.db.C.CreateTable(input)
	return err
}

// AppByID returns app from DynamoDB by ID. IDs are generated with https://github.com/rs/xid.
func (as *AppStorage) AppByID(id string) (model.AppData, error) {
	if len(id) == 0 {
		return model.AppData{}, model.ErrorWrongDataFormat
	}

	result, err := as.db.C.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(appsTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		log.Println("Error getting application:", err)
		return model.AppData{}, ErrorInternalError
	}

	if result.Item == nil {
		return model.AppData{}, model.ErrorNotFound
	}

	appdata := model.AppData{}
	if err = dynamodbattribute.UnmarshalMap(result.Item, &appdata); err != nil {
		log.Println("Error unmarshalling app data:", err)
		return model.AppData{}, ErrorInternalError
	}
	return appdata, nil
}

// ActiveAppByID returns app by id only if it's active.
func (as *AppStorage) ActiveAppByID(appID string) (model.AppData, error) {
	if appID == "" {
		return model.AppData{}, ErrorEmptyAppID
	}

	app, err := as.AppByID(appID)
	if err != nil {
		return model.AppData{}, err
	}

	if !app.Active {
		return model.AppData{}, ErrorInactiveApp
	}

	return app, nil
}

// CreateApp creates new app in DynamoDB.
func (as *AppStorage) CreateApp(app model.AppData) (model.AppData, error) {
	// generate new ID if it's not set
	if len(app.ID) == 0 {
		app.ID = xid.New().String()
	}

	av, err := dynamodbattribute.MarshalMap(app)
	if err != nil {
		log.Println("error marshalling app: ", err)
		return model.AppData{}, ErrorInternalError
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(appsTableName),
	}

	if _, err = as.db.C.PutItem(input); err != nil {
		log.Println("error putting app to storage: ", err)
		return model.AppData{}, ErrorInternalError
	}
	return app, nil
}

// DisableApp disables app in DynamoDB storage.
func (as *AppStorage) DisableApp(app model.AppData) error {
	if _, err := xid.FromString(app.ID); err != nil {
		log.Println("incorrect AppID: ", app.ID)
		return model.ErrorWrongDataFormat
	}
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {
				BOOL: aws.Bool(false),
			},
		},
		TableName: aws.String(appsTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(app.ID),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set active = :a"),
	}

	if _, err := as.db.C.UpdateItem(input); err != nil {
		log.Println("Error updating app:", err)
		return ErrorInternalError
	}
	return nil
}

// UpdateApp updates app in DynamoDB storage.
func (as *AppStorage) UpdateApp(appID string, app model.AppData) (model.AppData, error) {
	if _, err := xid.FromString(appID); err != nil {
		log.Println("incorrect appID: ", appID)
		return model.AppData{}, model.ErrorWrongDataFormat
	}

	// use ID from the request if it's not set
	if len(app.ID) == 0 {
		app.ID = appID
	}

	oldAppData := model.AppData{ID: appID}
	if err := as.DisableApp(oldAppData); err != nil {
		log.Println("Error disabling old app:", err)
		return model.AppData{}, err
	}

	updatedApp, err := as.CreateApp(app)
	return updatedApp, err
}

const maxAppsLimit = 20

// FetchApps fetches apps which name satisfies provided filterString.
// Supports pagination. Search is case-sensitive for now.
func (as *AppStorage) FetchApps(filterString string) ([]model.AppData, error) {
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(appsTableName),
	}

	if len(filterString) != 0 {
		scanInput.FilterExpression = aws.String("contains(#name, :filterStr)")
		scanInput.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{
			":filterStr": {S: aws.String(filterString)},
		}
		scanInput.ExpressionAttributeNames = map[string]*string{
			"#name": aws.String("name"),
		}
	}

	result, err := as.db.C.Scan(scanInput)
	if err != nil {
		log.Println("Error querying for apps:", err)
		return []model.AppData{}, ErrorInternalError
	}

	apps := make([]model.AppData, len(result.Items))
	for i := 0; i < len(result.Items); i++ {
		appData := model.AppData{}
		if err = dynamodbattribute.UnmarshalMap(result.Items[i], &appData); err != nil {
			log.Println("error while unmarshal app: ", err)
			return []model.AppData{}, ErrorInternalError
		}
		apps[i] = appData
	}
	return apps, nil
}

// DeleteApp deletes app by id.
func (as *AppStorage) DeleteApp(id string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
		TableName: aws.String(appsTableName),
	}
	_, err := as.db.C.DeleteItem(input)
	return err
}

// TestDatabaseConnection checks whether we can fetch the first document in the applications table.
func (as *AppStorage) TestDatabaseConnection() error {
	_, err := as.db.C.Scan(&dynamodb.ScanInput{
		TableName: aws.String(appsTableName),
		Limit:     aws.Int64(1),
	})
	return err
}

// ImportJSON imports data from JSON.
func (as *AppStorage) ImportJSON(data []byte, cleanOldData bool) error {
	if cleanOldData {
		as.db.DeleteTable(appsTableName)
		as.ensureTable()
	}
	apd := []model.AppData{}
	if err := json.Unmarshal(data, &apd); err != nil {
		log.Println("error while unmarshal app data: ", err)
		return err
	}
	for _, a := range apd {
		if _, err := as.CreateApp(a); err != nil {
			return err
		}
	}
	return nil
}

// Close does nothing here.
func (as *AppStorage) Close() {}
