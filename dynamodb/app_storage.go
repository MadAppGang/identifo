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

var (
	// AppsTable is table name for storing apps data.
	AppsTable = "Applications"
)

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

// ensureTable ensures that app table exists in database.
func (as *AppStorage) ensureTable() error {
	exists, err := as.db.isTableExists(AppsTable)
	if err != nil {
		log.Println("Error checking Applications table existence:", err)
		return err
	}
	if !exists {
		// create table, AWS DynamoDB table creation is overcomplicated for sure
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
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(10),
				WriteCapacityUnits: aws.Int64(10),
			},
			TableName: aws.String(AppsTable),
		}
		_, err = as.db.C.CreateTable(input)
	}
	return err
}

// AppByID returns app from DynamoDB by ID. IDs are generated with https://github.com/rs/xid.
func (as *AppStorage) AppByID(id string) (model.AppData, error) {
	if len(id) == 0 {
		return nil, model.ErrorWrongDataFormat
	}

	result, err := as.db.C.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(AppsTable),
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
	return AppData{appData: appdata}, nil
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

// AddNewApp add new app to dynamodb storage.
func (as *AppStorage) AddNewApp(app model.AppData) (model.AppData, error) {
	a, ok := app.(AppData)
	if !ok {
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
		TableName: aws.String(AppsTable),
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
		TableName: aws.String(AppsTable),
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
func (as *AppStorage) UpdateApp(oldAppID string, newApp model.AppData) error {
	if _, err := xid.FromString(oldAppID); err != nil {
		log.Println("Incorrect oldAppID: ", oldAppID)
		return model.ErrorWrongDataFormat
	}

	ad := AppData{appData: appData{ID: oldAppID}}
	if err := as.DisableApp(ad); err != nil {
		log.Println("Error disabling old app:", err)
		return err
	}
	_, err := as.AddNewApp(newApp)
	return err
}

// FetchApps fetches apps which name satisfies provided filterString.
// Supports pagination. Search is case-senstive for now.
func (as *AppStorage) FetchApps(filterString string, skip, limit int) ([]model.AppData, error) {
	result, err := as.db.C.Query(&dynamodb.QueryInput{
		TableName:              aws.String(AppsTable),
		KeyConditionExpression: aws.String("contains(name, :filterStr)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":filterStr": {S: aws.String(filterString)},
		},
		Select: aws.String("ALL_PROJECTED_ATTRIBUTES"),
	})
	if err != nil {
		log.Println("Error querying for apps:", err)
		return nil, ErrorInternalError
	}

	apps := make([]model.AppData, len(result.Items))
	for i := 0; i < len(result.Items); i++ {
		appData := appData{}
		if err = dynamodbattribute.UnmarshalMap(result.Items[i], &appData); err != nil {
			log.Println("Error unmarshalling app:", err)
			return nil, ErrorInternalError
		}
		apps[i] = AppData{appData: appData}
	}
	return apps, nil
}

// ImportJSON imports data from JSON.
func (as *AppStorage) ImportJSON(data []byte) error {
	apd := []appData{}
	if err := json.Unmarshal(data, &apd); err != nil {
		log.Println("Error unmarshalling app data:", err)
		return err
	}
	for _, a := range apd {
		if _, err := as.AddNewApp(AppData{appData: a}); err != nil {
			return err
		}
	}
	return nil
}

type appData struct {
	ID                   string   `json:"id,omitempty"`
	Secret               string   `json:"secret,omitempty"`
	Active               bool     `json:"active,omitempty"`
	Name                 string   `json:"name,omitempty"`
	Description          string   `json:"description,omitempty"`
	Scopes               []string `json:"scopes,omitempty"`
	Offline              bool     `json:"offline,omitempty"`
	RedirectURL          string   `json:"redirect_url,omitempty"`
	RefreshTokenLifespan int64    `json:"refresh_token_lifespan,omitempty"`
	TokenLifespan        int64    `json:"token_lifespan,omitempty"`
	TokenPayload         []string `bson:"token_payload,omitempty" json:"token_payload,omitempty"`
}

// AppData is DynamoDB model for model.AppData.
type AppData struct {
	appData
}

// NewAppData instantiates DynamoDB app data model from the general one.
func NewAppData(data model.AppData) (AppData, error) {
	if _, err := xid.FromString(data.ID()); err != nil {
		log.Println("Incorrect AppID: ", data.ID())
		return AppData{}, model.ErrorWrongDataFormat
	}
	return AppData{appData: appData{
		ID:                   data.ID(),
		Secret:               data.Secret(),
		Active:               data.Active(),
		Name:                 data.Name(),
		Description:          data.Description(),
		Scopes:               data.Scopes(),
		Offline:              data.Offline(),
		RedirectURL:          data.RedirectURL(),
		RefreshTokenLifespan: data.RefreshTokenLifespan(),
		TokenLifespan:        data.TokenLifespan(),
		TokenPayload:         data.TokenPayload(),
	}}, nil
}

// AppDataFromJSON deserializes data from JSON.
func AppDataFromJSON(d []byte) (AppData, error) {
	apd := appData{}
	if err := json.Unmarshal(d, &apd); err != nil {
		log.Println(err)
		return AppData{}, err
	}
	return AppData{appData: apd}, nil
}

// Marshal serializes data to byte array.
func (ad AppData) Marshal() ([]byte, error) {
	return json.Marshal(ad.appData)
}

// MakeAppData creates new DynamoDB app data instance.
func MakeAppData(id, secret string, active bool, name, description string, scopes []string, offline bool, redirectURL string, refreshTokenLifespan, tokenLifespan int64) (AppData, error) {
	if _, err := xid.FromString(id); err != nil {
		log.Println("Cannot create ID from the string representation:", err)
		return AppData{}, model.ErrorWrongDataFormat
	}
	return AppData{appData: appData{
		ID:                   id,
		Secret:               secret,
		Active:               active,
		Name:                 name,
		Description:          description,
		Scopes:               scopes,
		Offline:              offline,
		RedirectURL:          redirectURL,
		RefreshTokenLifespan: refreshTokenLifespan,
		TokenLifespan:        tokenLifespan,
	}}, nil
}

// ID implements model.AppData interface.
func (ad AppData) ID() string { return ad.appData.ID }

// Secret implements model.AppData interface.
func (ad AppData) Secret() string { return ad.appData.Secret }

// Active implements model.AppData interface.
func (ad AppData) Active() bool { return ad.appData.Active }

// Name implements model.AppData interface.
func (ad AppData) Name() string { return ad.appData.Name }

// Description implements model.AppData interface.
func (ad AppData) Description() string { return ad.appData.Description }

// Scopes implements model.AppData interface.
func (ad AppData) Scopes() []string { return ad.appData.Scopes }

// Offline implements model.AppData interface.
func (ad AppData) Offline() bool { return ad.appData.Offline }

// RedirectURL implements model.AppData interface.
func (ad AppData) RedirectURL() string { return ad.appData.RedirectURL }

// RefreshTokenLifespan implements model.AppData interface.
func (ad AppData) RefreshTokenLifespan() int64 { return ad.appData.RefreshTokenLifespan }

// TokenLifespan implements model.AppData interface.
func (ad AppData) TokenLifespan() int64 { return ad.appData.TokenLifespan }

// TokenPayload implements model.AppData interface.
func (ad AppData) TokenPayload() []string { return ad.appData.TokenPayload }
