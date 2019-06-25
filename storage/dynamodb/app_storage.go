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

// NewAppData returns pointer to newly created app data.
func (as *AppStorage) NewAppData() model.AppData {
	return &AppData{appData: appData{}}
}

// ensureTable ensures that app table exists in database.
func (as *AppStorage) ensureTable() error {
	exists, err := as.db.IsTableExists(AppsTable)
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

// CreateApp creates new app in DynamoDB.
func (as *AppStorage) CreateApp(app model.AppData) (model.AppData, error) {
	res, ok := app.(*AppData)
	if !ok || app == nil {
		return nil, model.ErrorWrongDataFormat
	}
	result, err := as.addNewApp(*res)
	return result, err
}

// addNewApp adds new app to DynamoDB storage.
func (as *AppStorage) addNewApp(app model.AppData) (model.AppData, error) {
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

	oldAppData := AppData{appData: appData{ID: appID}}
	if err := as.DisableApp(oldAppData); err != nil {
		log.Println("Error disabling old app:", err)
		return nil, err
	}

	updatedApp, err := as.addNewApp(*res)
	return updatedApp, err
}

// FetchApps fetches apps which name satisfies provided filterString.
// Supports pagination. Search is case-senstive for now.
func (as *AppStorage) FetchApps(filterString string, skip, limit int) ([]model.AppData, int, error) {
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(AppsTable),
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
		apps[i] = AppData{appData: appData}
	}
	return apps, len(result.Items), nil
}

// DeleteApp deletes app by id.
func (as *AppStorage) DeleteApp(id string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
		TableName: aws.String(AppsTable),
	}
	_, err := as.db.C.DeleteItem(input)
	return err
}

// TestDatabaseConnection checks whether we can fetch the first document in the applications table.
func (as *AppStorage) TestDatabaseConnection() error {
	_, err := as.db.C.Scan(&dynamodb.ScanInput{
		TableName: aws.String(AppsTable),
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
		if _, err := as.addNewApp(AppData{appData: a}); err != nil {
			return err
		}
	}
	return nil
}

type appData struct {
	ID                    string           `json:"id,omitempty"`
	Secret                string           `json:"secret,omitempty"`
	Active                bool             `json:"active,omitempty"`
	Name                  string           `json:"name,omitempty"`
	Description           string           `json:"description,omitempty"`
	Scopes                []string         `json:"scopes,omitempty"`
	Offline               bool             `json:"offline,omitempty"`
	Type                  model.AppType    `json:"type,omitempty"`
	RedirectURL           string           `json:"redirect_url,omitempty"`
	RefreshTokenLifespan  int64            `json:"refresh_token_lifespan,omitempty"`
	InviteTokenLifespan   int64            `json:"invite_token_lifespan,omitempty"`
	TokenLifespan         int64            `json:"token_lifespan,omitempty"`
	TokenPayload          []string         `bson:"token_payload,omitempty" json:"token_payload,omitempty"`
	RegistrationForbidden bool             `json:"registration_forbidden,omitempty"`
	AppleInfo             *model.AppleInfo `json:"apple_info,omitempty"`
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
		ID:                    data.ID(),
		Secret:                data.Secret(),
		Active:                data.Active(),
		Name:                  data.Name(),
		Description:           data.Description(),
		Scopes:                data.Scopes(),
		Offline:               data.Offline(),
		RedirectURL:           data.RedirectURL(),
		RefreshTokenLifespan:  data.RefreshTokenLifespan(),
		InviteTokenLifespan:   data.InviteTokenLifespan(),
		TokenLifespan:         data.TokenLifespan(),
		TokenPayload:          data.TokenPayload(),
		RegistrationForbidden: data.RegistrationForbidden(),
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
func MakeAppData(id, secret string, active bool, name, description string, scopes []string, offline bool, redirectURL string, refreshTokenLifespan, inviteTokenLifespan, tokenLifespan int64, registrationForbidden bool) (AppData, error) {
	if _, err := xid.FromString(id); err != nil {
		log.Println("Cannot create ID from the string representation:", err)
		return AppData{}, model.ErrorWrongDataFormat
	}
	return AppData{appData: appData{
		ID:                    id,
		Secret:                secret,
		Active:                active,
		Name:                  name,
		Description:           description,
		Scopes:                scopes,
		Offline:               offline,
		RedirectURL:           redirectURL,
		RefreshTokenLifespan:  refreshTokenLifespan,
		InviteTokenLifespan:   inviteTokenLifespan,
		TokenLifespan:         tokenLifespan,
		RegistrationForbidden: registrationForbidden,
	}}, nil
}

// Sanitize removes all sensitive data.
func (ad AppData) Sanitize() model.AppData {
	ad.appData.Secret = ""
	if ad.appData.AppleInfo != nil {
		ad.appData.AppleInfo.ClientSecret = ""
	}
	return ad
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

// Type implements model.AppData interface.
func (ad AppData) Type() model.AppType { return ad.appData.Type }

// RedirectURL implements model.AppData interface.
func (ad AppData) RedirectURL() string { return ad.appData.RedirectURL }

// RefreshTokenLifespan implements model.AppData interface.
func (ad AppData) RefreshTokenLifespan() int64 { return ad.appData.RefreshTokenLifespan }

// InviteTokenLifespan a inviteToken lifespan in seconds, if 0 - default one is used.
func (ad AppData) InviteTokenLifespan() int64 { return ad.appData.InviteTokenLifespan }

// TokenLifespan implements model.AppData interface.
func (ad AppData) TokenLifespan() int64 { return ad.appData.TokenLifespan }

// TokenPayload implements model.AppData interface.
func (ad AppData) TokenPayload() []string { return ad.appData.TokenPayload }

// RegistrationForbidden implements model.AppData interface.
func (ad AppData) RegistrationForbidden() bool { return ad.appData.RegistrationForbidden }

// AppleInfo implements model.AppData interface.
func (ad AppData) AppleInfo() *model.AppleInfo { return ad.appData.AppleInfo }
