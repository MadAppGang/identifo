package dynamodb

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/madappgang/identifo/model"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

const (
	//UsersTableName where to store users
	UsersTableName = "Users"
	//UserTableEmailIndexName user table global index name to access by email
	UserTableEmailIndexName = "EmailIndex"
	//UsersFederatedIDTableName table to store federatedId's for user
	//beacuse Dynamodb does not support
	UsersFederatedIDTableName = "UsersByFederatedID"
)

//NewUserStorage crates and provision new user storage instance
func NewUserStorage(db *DB) (model.UserStorage, error) {
	us := UserStorage{}
	us.db = db
	(&us).ensureTable()
	return &us, nil
}

//UserStorage stores and manages data in dynamodb sotrage
type UserStorage struct {
	db *DB
}

//UserByID returns user by it's ID
func (us *UserStorage) UserByID(id string) (model.User, error) {
	idx, err := xid.FromString(id)
	if err != nil {
		return nil, model.ErrorWrongDataFormat
	}

	result, err := us.db.C.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(UsersTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(idx.String()),
			},
		},
	})
	if err != nil {
		return nil, ErrorInternalError
	}
	//empty result
	if result.Item == nil {
		return nil, model.ErrorNotFound
	}
	userdata := userData{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &userdata)
	if err != nil {
		return nil, ErrorInternalError
	}
	return &User{userData: userdata}, nil
}

func (us *UserStorage) userIDByFederatedID(provider model.FederatedIdentityProvider, id string) (string, error) {
	fid := string(provider) + ":" + id
	result, err := us.db.C.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(UsersFederatedIDTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"federated_id": {
				S: aws.String(fid),
			},
		},
	})
	if err != nil {
		return "", ErrorInternalError
	}
	//empty result
	if result.Item == nil {
		return "", model.ErrorNotFound
	}
	fidd := federatedUserID{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &fidd)
	if err != nil || len(fidd.UserID) == 0 {
		return "", ErrorInternalError
	}
	return fidd.UserID, nil
}

//UserByFederatedID returns user by federated ID
func (us *UserStorage) UserByFederatedID(provider model.FederatedIdentityProvider, id string) (model.User, error) {
	userID, err := us.userIDByFederatedID(provider, id)
	if err != nil {
		return nil, err
	}
	if len(userID) == 0 {
		return nil, model.ErrorWrongDataFormat
	}
	return us.UserByID(userID)
}

//AttachDeviceToken do nothing here
//TODO: implement device storage
func (us *UserStorage) AttachDeviceToken(id, token string) error {
	//we are not supporting devices for users here
	return model.ErrorNotImplemented
}

//DetachDeviceToken do nothing here yet
//TODO: implement
func (us *UserStorage) DetachDeviceToken(token string) error {
	return model.ErrorNotImplemented
}

//RequestScopes for now returns requested scope
//TODO: implement scope logic
func (us *UserStorage) RequestScopes(userID string, scopes []string) ([]string, error) {
	return scopes, nil
}

func (us *UserStorage) userByName(name string) (*userData, error) {
	name = strings.ToLower(name)
	result, err := us.db.C.Query(&dynamodb.QueryInput{
		TableName:              aws.String(UsersTableName),
		IndexName:              aws.String(UserTableEmailIndexName),
		KeyConditionExpression: aws.String("username = :n"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":n": {S: aws.String(name)},
		},
		Select: aws.String("ALL_ATTRIBUTES"), //retrieve all attributes, because we need to make local check.
	})
	if err != nil {
		return nil, ErrorInternalError
	}
	//empty result
	if len(result.Items) == 0 {
		return nil, model.ErrorNotFound
	}
	item := result.Items[0]
	userdata := userData{}
	err = dynamodbattribute.UnmarshalMap(item, &userdata)
	if err != nil {
		return nil, ErrorInternalError
	}
	return &userdata, nil
}

//UserByNamePassword returns  user by name and password
func (us *UserStorage) UserByNamePassword(name, password string) (model.User, error) {
	name = strings.ToLower(name)
	userdata, err := us.userByName(name)
	if err != nil {
		return nil, err
	}
	//if password is incorrect, returning not found error for secure reason
	if bcrypt.CompareHashAndPassword([]byte(userdata.Pswd), []byte(password)) != nil {
		return nil, model.ErrorNotFound
	}
	u := &User{userData: *userdata}
	u.Sanitize()
	return u, nil
}

//AddNewUser adds new user
func (us *UserStorage) AddNewUser(usr model.User, password string) (model.User, error) {
	u, ok := usr.(*User)
	if !ok {
		return nil, model.ErrorWrongDataFormat
	}
	//generate new ID if it's not set
	if _, err := xid.FromString(u.ID()); err != nil {
		u.userData.ID = xid.New().String()
	}
	if len(password) > 0 {
		u.userData.Pswd = PasswordHash(password)
	}

	u.userData.Name = strings.ToLower(u.userData.Name)
	uv, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, ErrorInternalError
	}

	input := &dynamodb.PutItemInput{
		Item:      uv,
		TableName: aws.String(UsersTableName),
	}
	_, err = us.db.C.PutItem(input)
	if err != nil {
		return nil, ErrorInternalError
	}
	u.Sanitize()
	return u, nil
}

//AddUserByNameAndPassword register new user
func (us *UserStorage) AddUserByNameAndPassword(name, password string, profile map[string]interface{}) (model.User, error) {
	name = strings.ToLower(name)
	_, err := us.userByName(name)
	if err != nil && err != model.ErrorNotFound {
		return nil, err
	} else if err == nil {
		return nil, model.ErrorUserExists
	}
	u := userData{}
	u.Active = true
	u.Name = name
	u.Profile = profile
	return us.AddNewUser(&User{u}, password)
}

//AddUserWithFederatedID add new user with social ID
func (us *UserStorage) AddUserWithFederatedID(provider model.FederatedIdentityProvider, federatedID string) (model.User, error) {
	_, err := us.userIDByFederatedID(provider, federatedID)
	if err != nil && err != model.ErrorNotFound {
		return nil, err
	} else if err == nil {
		return nil, model.ErrorUserExists
	}

	fid := string(provider) + ":" + federatedID

	uu, err := us.userByName(fid)
	//error getting user
	if err != nil && err != model.ErrorNotFound {
		return nil, err
	} else if err == model.ErrorNotFound {
		//no such user, let's create it
		u := userData{}
		u.Name = fid
		u.Active = true
		var independentError error
		uuu, independentError := us.AddNewUser(&User{u}, "")
		if independentError != nil {
			return nil, independentError
		}
		uu = &(uuu.(*User).userData) //yep, looks like old C :-), payment for interfaces
	}
	//if no error it means there is already user for this federated id somehow,
	//the only possible way for that is faulty creation of the federated accout before

	fedData := federatedUserID{}
	fedData.FederatedID = fid
	fedData.UserID = uu.ID
	fedInputData, err := dynamodbattribute.MarshalMap(fedData)
	if err != nil {
		return nil, ErrorInternalError
	}

	input := &dynamodb.PutItemInput{
		Item:      fedInputData,
		TableName: aws.String(UsersFederatedIDTableName),
	}
	_, err = us.db.C.PutItem(input)
	if err != nil {
		return nil, ErrorInternalError
	}
	//just in case
	if uu == nil {
		return nil, ErrorInternalError
	}
	resultUser := &User{*uu}
	resultUser.Sanitize()
	return resultUser, nil
}

//data implementation
type userData struct {
	ID      string                 `json:"id,omitempty"`
	Name    string                 `json:"username,omitempty"`
	Pswd    string                 `json:"pswd,omitempty"`
	Profile map[string]interface{} `json:"profile,omitempty"`
	Active  bool                   `json:"active,omitempty"`
}

//federatedUserID is storage for federated if to user mapping implementation
type federatedUserID struct {
	FederatedID string `json:"federated_id,omitempty"`
	UserID      string `json:"user_id,omitempty"`
}

//User user data structure for dynamodb storage
type User struct {
	userData
}

//Sanitize removes sensitive data
func (u *User) Sanitize() {
	u.userData.Pswd = ""
}

//UserFromJSON deserializes data
func UserFromJSON(d []byte) (*User, error) {
	user := userData{}
	if err := json.Unmarshal(d, &user); err != nil {
		return &User{}, err
	}
	return &User{user}, nil
}

//model.User interface implementation
func (u *User) ID() string                      { return u.userData.ID }
func (u *User) Name() string                    { return u.userData.Name }
func (u *User) PasswordHash() string            { return u.userData.Pswd }
func (u *User) Profile() map[string]interface{} { return u.userData.Profile }
func (u *User) Active() bool                    { return u.userData.Active }

//PasswordHash creates hash with salt for password
func PasswordHash(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash)
}

//ensureTable ensures app storage table is exists in database
//I'm hiding it in the end of the file, because AWS devs, you are killing me with this API
func (us *UserStorage) ensureTable() error {
	exists, err := us.db.isTableExists(UsersTableName)
	if err != nil {
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
				{
					AttributeName: aws.String("username"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("id"),
					KeyType:       aws.String("HASH"),
				},
			},
			GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
				{
					IndexName: aws.String(UserTableEmailIndexName),
					KeySchema: []*dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("username"),
							KeyType:       aws.String("HASH"),
						},
					},
					//we are doing local password check. As a result,  we don't need this projections
					// Projection: &dynamodb.Projection{
					// 	NonKeyAttributes: []*string{aws.String("pswd"), aws.String("id")},
					// 	ProjectionType:   aws.String("INCLUDE"),
					// },
					Projection: &dynamodb.Projection{
						ProjectionType: aws.String("KEYS_ONLY"),
					},
					ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(10),
						WriteCapacityUnits: aws.Int64(10),
					},
				},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(10),
				WriteCapacityUnits: aws.Int64(10),
			},
			TableName: aws.String(UsersTableName),
		}
		_, err = us.db.C.CreateTable(input)
		if err != nil {
			return err
		}
	}

	//create table to handle federated ID's
	exists, err = us.db.isTableExists(UsersFederatedIDTableName)
	if err != nil {
		return err
	}
	if !exists {
		//create table, AWS DynamoDB table creation is overcomplicated for sure
		input := &dynamodb.CreateTableInput{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("federated_id"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("federated_id"),
					KeyType:       aws.String("HASH"),
				},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(10),
				WriteCapacityUnits: aws.Int64(10),
			},
			TableName: aws.String(UsersFederatedIDTableName),
		}
		_, err = us.db.C.CreateTable(input)
		if err != nil {
			return err
		}
	}
	return nil
}
