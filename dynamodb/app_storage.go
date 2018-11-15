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
	//AppsTable table name to store apps data
	AppsTable = "Applications"
)

//NewAppStorage creates new dynamoDB AppStorage implementation
func NewAppStorage(db *DB) (model.AppStorage, error) {
	as := &AppStorage{db: db}
	err := as.ensureTable()
	return as, err
}

//AppStorage is fully functional app storage
type AppStorage struct {
	db *DB
}

//ensureTable ensures app storage table is exists in database
func (as *AppStorage) ensureTable() error {
	exists, err := as.db.isTableExists(AppsTable)
	if err != nil {
		log.Println(err)
		return err
	}
	if !exists {
		//create table, AWS DynamoDB table creation is overcomplicated for sure
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

//AppByID returns app from dynamodb by ID
//ID is generated with https://github.com/rs/xid
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
		log.Println(err)
		return nil, ErrorInternalError
	}
	//empty result
	if result.Item == nil {
		return nil, model.ErrorNotFound
	}
	appdata := appData{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &appdata)
	if err != nil {
		log.Println(err)
		return nil, ErrorInternalError
	}
	return AppData{appData: appdata}, nil
}

//AddNewApp add new app to dynamodb storage
func (as *AppStorage) AddNewApp(app model.AppData) (model.AppData, error) {
	a, ok := app.(AppData)
	if !ok {
		return nil, model.ErrorWrongDataFormat
	}
	//generate new ID if it's not set
	if len(a.ID()) == 0 {
		a.appData.ID = xid.New().String()
	}

	av, err := dynamodbattribute.MarshalMap(a)
	if err != nil {
		log.Println(err)
		return nil, ErrorInternalError
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(AppsTable),
	}

	if _, err = as.db.C.PutItem(input); err != nil {
		log.Println(err)
		return nil, ErrorInternalError
	}
	return a, nil
}

//DisableApp disables app in dynamodb storage
func (as *AppStorage) DisableApp(app model.AppData) error {
	_, err := xid.FromString(app.ID())
	if err != nil {
		log.Println("wrong AppID: ", app.ID())
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
	_, err = as.db.C.UpdateItem(input)
	//this approach was to delete record for dearivated app, leave it here for a while
	// input := &dynamodb.DeleteItemInput{
	// 	Key: map[string]*dynamodb.AttributeValue{
	// 		"id": {
	// 			S: aws.String(app.ID()),
	// 		},
	// 	},
	// 	TableName: aws.String(AppsTable),
	// }
	// _, err = as.db.C.DeleteItem(input)

	if err != nil {
		log.Println(err)
		return ErrorInternalError
	}
	return nil
}

//UpdateApp updates app in dynamodb storage
func (as *AppStorage) UpdateApp(oldAppID string, newApp model.AppData) error {
	if _, err := xid.FromString(oldAppID); err != nil {
		log.Println("wrong oldAppID: ", oldAppID)
		return model.ErrorWrongDataFormat
	}

	ad := AppData{}
	ad.appData.ID = oldAppID
	if err := as.DisableApp(ad); err != nil {
		log.Println(err)
		return err
	}
	_, err := as.AddNewApp(newApp)
	return err
}

type appData struct {
	ID                   string   `json:"id,omitempty"`
	Secret               string   `json:"secret,omitempty"`
	Active               bool     `json:"active,omitempty"`
	Description          string   `json:"description,omitempty"`
	Scopes               []string `json:"scopes,omitempty"`
	Offline              bool     `json:"offline,omitempty"`
	RedirectURL          string   `json:"redirect_url,omitempty"`
	RefreshTokenLifespan int64    `json:"refresh_token_lifespan,omitempty"`
	TokenLifespan        int64    `json:"token_lifespan,omitempty"`
}

//AppData is mongo model for model.AppData
type AppData struct {
	appData
}

//NewAppData instantiate app data mongo model from general one
func NewAppData(data model.AppData) (AppData, error) {
	_, err := xid.FromString(data.ID())
	if err != nil {
		log.Println("wrong AppID: ", data.ID())
		return AppData{}, model.ErrorWrongDataFormat
	}
	return AppData{appData: appData{
		ID:                   data.ID(),
		Secret:               data.Secret(),
		Active:               data.Active(),
		Description:          data.Description(),
		Scopes:               data.Scopes(),
		Offline:              data.Offline(),
		RedirectURL:          data.RedirectURL(),
		RefreshTokenLifespan: data.RefreshTokenLifespan(),
		TokenLifespan:        data.TokenLifespan(),
	}}, nil
}

//AppDataFromJSON deserializes data from JSON
func AppDataFromJSON(d []byte) (AppData, error) {
	apd := appData{}
	if err := json.Unmarshal(d, &apd); err != nil {
		log.Println(err)
		return AppData{}, err
	}
	return AppData{appData: apd}, nil
}

//Marshal serialize data to byte array
func (ad AppData) Marshal() ([]byte, error) {
	return json.Marshal(ad.appData)
}

//MakeAppData creates new mongo app data instance
func MakeAppData(id, secret string, active bool, description string, scopes []string, offline bool, redirectURL string, refreshTokenLifespan, tokenLifespan int64) (AppData, error) {
	if _, err := xid.FromString(id); err != nil {
		log.Println(err)
		return AppData{}, model.ErrorWrongDataFormat
	}
	return AppData{appData: appData{
		ID:                   id,
		Secret:               secret,
		Active:               active,
		Description:          description,
		Scopes:               scopes,
		Offline:              offline,
		RedirectURL:          redirectURL,
		RefreshTokenLifespan: refreshTokenLifespan,
		TokenLifespan:        tokenLifespan,
	}}, nil
}

func (ad AppData) ID() string                  { return ad.appData.ID }
func (ad AppData) Secret() string              { return ad.appData.Secret }
func (ad AppData) Active() bool                { return ad.appData.Active }
func (ad AppData) Description() string         { return ad.appData.Description }
func (ad AppData) Scopes() []string            { return ad.appData.Scopes }
func (ad AppData) Offline() bool               { return ad.appData.Offline }
func (ad AppData) RedirectURL() string         { return ad.appData.RedirectURL }
func (ad AppData) RefreshTokenLifespan() int64 { return ad.appData.RefreshTokenLifespan }
func (ad AppData) TokenLifespan() int64        { return ad.appData.TokenLifespan }
