package dynamodb

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/madappgang/identifo/model"
	"github.com/rs/xid"
)

var appsTableName = "Applications"

// NewAppStorage creates new DynamoDB AppStorage implementation.
func NewAppStorage(db *DB) (model.AppStorage, error) {
	as := &AppStorage{db: db}
	err := as.ensureTable()
	return as, err
}

// AppStorage a is fully functional app storage.
type AppStorage struct {
	db *DB
}

// NewAppData returns pointer to newly created app data.
func (as *AppStorage) NewAppData() model.AppData {
	return &AppData{appData: appData{}}
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
		return nil, model.ErrorWrongDataFormat
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
		return nil, ErrorInternalError
	}

	if result.Item == nil {
		return nil, model.ErrorNotFound
	}

	appdata := appData{}
	if err = dynamodbattribute.UnmarshalMap(result.Item, &appdata); err != nil {
		log.Println("Error unmarshalling app data:", err)
		return nil, ErrorInternalError
	}
	return &AppData{appData: appdata}, nil
}

// ActiveAppByID returns app by id only if it's active.
func (as *AppStorage) ActiveAppByID(appID string) (model.AppData, error) {
	if appID == "" {
		return nil, ErrorEmptyAppID
	}

	app, err := as.AppByID(appID)
	if err != nil {
		return nil, err
	}

	if !app.Active() {
		return nil, ErrorInactiveApp
	}

	return app, nil
}

// CreateApp creates new app in DynamoDB.
func (as *AppStorage) CreateApp(app model.AppData) (model.AppData, error) {
	res, ok := app.(*AppData)
	if !ok || res == nil {
		return nil, model.ErrorWrongDataFormat
	}
	result, err := as.addNewApp(res)
	return result, err
}

// addNewApp adds new app to DynamoDB storage.
func (as *AppStorage) addNewApp(app model.AppData) (model.AppData, error) {
	a, ok := app.(*AppData)
	if !ok || a == nil {
		return nil, model.ErrorWrongDataFormat
	}
	// generate new ID if it's not set
	if len(a.ID()) == 0 {
		a.appData.ID = xid.New().String()
	}

	av, err := dynamodbattribute.MarshalMap(a)
	if err != nil {
		log.Println("Error marshalling app:", err)
		return nil, ErrorInternalError
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(appsTableName),
	}

	if _, err = as.db.C.PutItem(input); err != nil {
		log.Println("Error putting app to storage:", err)
		return nil, ErrorInternalError
	}
	return a, nil
}

// DisableApp disables app in DynamoDB storage.
func (as *AppStorage) DisableApp(app model.AppData) error {
	if _, err := xid.FromString(app.ID()); err != nil {
		log.Println("Incorrect AppID: ", app.ID())
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
				S: aws.String(app.ID()),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set active = :a"),
	}

	// Old approach was to delete the record for deactivated app, leave it here for a while.
	// input := &dynamodb.DeleteItemInput{
	// 	Key: map[string]*dynamodb.AttributeValue{
	// 		"id": {
	// 			S: aws.String(app.ID()),
	// 		},
	// 	},
	// 	TableName: aws.String(AppsTable),
	// }
	// _, err = as.db.C.DeleteItem(input)

	if _, err := as.db.C.UpdateItem(input); err != nil {
		log.Println("Error updating app:", err)
		return ErrorInternalError
	}
	return nil
}

// UpdateApp updates app in DynamoDB storage.
func (as *AppStorage) UpdateApp(appID string, newApp model.AppData) (model.AppData, error) {
	if _, err := xid.FromString(appID); err != nil {
		log.Println("Incorrect appID: ", appID)
		return nil, model.ErrorWrongDataFormat
	}

	res, ok := newApp.(*AppData)
	if !ok || res == nil {
		return nil, model.ErrorWrongDataFormat
	}

	// use ID from the request if it's not set
	if len(newApp.ID()) == 0 {
		res.appData.ID = appID
	}

	oldAppData := &AppData{appData: appData{ID: appID}}
	if err := as.DisableApp(oldAppData); err != nil {
		log.Println("Error disabling old app:", err)
		return nil, err
	}

	updatedApp, err := as.addNewApp(res)
	return updatedApp, err
}

// FetchApps fetches apps which name satisfies provided filterString.
// Supports pagination. Search is case-senstive for now.
func (as *AppStorage) FetchApps(filterString string, skip, limit int) ([]model.AppData, int, error) {
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(appsTableName),
		Limit:     aws.Int64(int64(limit)),
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
		return []model.AppData{}, 0, ErrorInternalError
	}

	apps := make([]model.AppData, len(result.Items))
	for i := 0; i < len(result.Items); i++ {
		if i < skip {
			continue // TODO: use internal pagination mechanism
		}
		appData := appData{}
		if err = dynamodbattribute.UnmarshalMap(result.Items[i], &appData); err != nil {
			log.Println("Error unmarshalling app:", err)
			return []model.AppData{}, 0, ErrorInternalError
		}
		apps[i] = &AppData{appData: appData}
	}
	return apps, len(result.Items), nil
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
func (as *AppStorage) ImportJSON(data []byte) error {
	apd := []appData{}
	if err := json.Unmarshal(data, &apd); err != nil {
		log.Println("Error unmarshalling app data:", err)
		return err
	}
	for _, a := range apd {
		if _, err := as.addNewApp(&AppData{appData: a}); err != nil {
			return err
		}
	}
	return nil
}

// Close does nothing here.
func (as *AppStorage) Close() {}
