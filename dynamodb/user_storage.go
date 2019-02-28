package dynamodb

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/madappgang/identifo/model"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

const (
	// UsersTableName is a table where to store users.
	UsersTableName = "Users"
	// UserTableEmailIndexName is a user table global index name to access by users by email.
	UserTableEmailIndexName = "EmailIndex"
	// UsersFederatedIDTableName is a table to store federated ids.
	UsersFederatedIDTableName = "UsersByFederatedID"
)

// NewUserStorage creates and provisions new user storage instance.
func NewUserStorage(db *DB) (model.UserStorage, error) {
	us := &UserStorage{db: db}
	err := us.ensureTable()
	return us, err
}

// UserStorage stores and manages data in DynamoDB storage.
type UserStorage struct {
	db *DB
}

// UserByID returns user by its ID.
func (us *UserStorage) UserByID(id string) (model.User, error) {
	idx, err := xid.FromString(id)
	if err != nil {
		log.Println("Incorrect user ID: ", id)
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
		log.Println("Error getting item from DynamoDB:", err)
		return nil, ErrorInternalError
	}
	if result.Item == nil {
		return nil, model.ErrorNotFound
	}

	userdata := userData{}
	if err = dynamodbattribute.UnmarshalMap(result.Item, &userdata); err != nil {
		log.Println("Error unmarshalling item:", err)
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
		log.Println("Error getting item from DynamoDB:", err)
		return "", ErrorInternalError
	}
	if result.Item == nil {
		return "", model.ErrorNotFound
	}

	fedData := federatedUserID{}
	if err = dynamodbattribute.UnmarshalMap(result.Item, &fedData); err != nil || len(fedData.UserID) == 0 {
		log.Println("Error unmarshalling item:", err)
		return "", ErrorInternalError
	}
	return fedData.UserID, nil
}

// UserByFederatedID returns user by federated ID.
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

// UserExists checks if user with provided name exists.
func (us *UserStorage) UserExists(name string) bool {
	_, err := us.userIdxByName(name)
	return err == nil
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

// Scopes returns supported scopes, could be static data of database.
func (us *UserStorage) Scopes() []string {
	// we allow all scopes for embedded database, you could implement your own logic in external service
	return []string{"offline", "user"}
}

// userIdxByName returns user data projected on the email index.
func (us *UserStorage) userIdxByName(name string) (*userIndexByNameData, error) {
	name = strings.ToLower(name)
	result, err := us.db.C.Query(&dynamodb.QueryInput{
		TableName:              aws.String(UsersTableName),
		IndexName:              aws.String(UserTableEmailIndexName),
		KeyConditionExpression: aws.String("username = :n"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":n": {S: aws.String(name)},
		},
		Select: aws.String("ALL_PROJECTED_ATTRIBUTES"), // retrieve all attributes, because we need to make local check.
	})
	if err != nil {
		log.Println("Error querying for items:", err)
		return nil, ErrorInternalError
	}
	if len(result.Items) == 0 {
		return nil, model.ErrorNotFound
	}

	item := result.Items[0]
	userdata := new(userIndexByNameData)
	if err = dynamodbattribute.UnmarshalMap(item, userdata); err != nil {
		log.Println("Error unmarshalling item:", err)
		return nil, ErrorInternalError
	}
	return userdata, nil
}

// UserByNamePassword returns user by name and password.
func (us *UserStorage) UserByNamePassword(name, password string) (model.User, error) {
	name = strings.ToLower(name)
	userIdx, err := us.userIdxByName(name)
	if err != nil {
		log.Println("Error getting user by name:", err)
		return nil, err
	}
	// if password is incorrect, return 'not found' error for security reasons.
	if bcrypt.CompareHashAndPassword([]byte(userIdx.Pswd), []byte(password)) != nil {
		return nil, model.ErrorNotFound
	}

	user, err := us.UserByID(userIdx.ID)
	if err != nil {
		log.Println("Error querying user by id:", err)
		return nil, ErrorInternalError
	}
	user.Sanitize()
	return user, nil
}

// AddNewUser adds new user.
func (us *UserStorage) AddNewUser(usr model.User, password string) (model.User, error) {
	u, ok := usr.(*User)
	if !ok {
		return nil, model.ErrorWrongDataFormat
	}
	// generate new ID if it's not set
	if _, err := xid.FromString(u.ID()); err != nil {
		u.userData.ID = xid.New().String()
	}
	if len(password) > 0 {
		u.userData.Pswd = PasswordHash(password)
	}

	u.userData.Name = strings.ToLower(u.userData.Name)
	uv, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		log.Println("Error marshalling user:", err)
		return nil, ErrorInternalError
	}

	input := &dynamodb.PutItemInput{
		Item:      uv,
		TableName: aws.String(UsersTableName),
	}
	if _, err = us.db.C.PutItem(input); err != nil {
		log.Println("Error putting item:", err)
		return nil, ErrorInternalError
	}
	u.Sanitize()
	return u, nil
}

// DeleteUser deletes user by id.
func (us *UserStorage) DeleteUser(id string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
		TableName: aws.String(UsersTableName),
	}
	_, err := us.db.C.DeleteItem(input)
	return err
}

//AddUserByNameAndPassword register new user
func (us *UserStorage) AddUserByNameAndPassword(name, password string, profile map[string]interface{}) (model.User, error) {
	name = strings.ToLower(name)
	_, err := us.userIdxByName(name)
	if err != nil && err != model.ErrorNotFound {
		log.Println(err)
		return nil, err
	} else if err == nil {
		return nil, model.ErrorUserExists
	}
	u := userData{Active: true, Name: name, Profile: profile}
	return us.AddNewUser(&User{userData: u}, password)
}

// AddUserWithFederatedID adds new user with social ID.
func (us *UserStorage) AddUserWithFederatedID(provider model.FederatedIdentityProvider, federatedID string) (model.User, error) {
	_, err := us.userIDByFederatedID(provider, federatedID)
	if err != nil && err != model.ErrorNotFound {
		log.Println("Error getting user by name:", err)
		return nil, err
	} else if err == nil {
		return nil, model.ErrorUserExists
	}

	fid := string(provider) + ":" + federatedID

	user, err := us.userIdxByName(fid)
	if err != nil && err != model.ErrorNotFound {
		log.Println("Error getting user by name:", err)
		return nil, err
	} else if err == model.ErrorNotFound {
		// no such user, let's create it
		uData := userData{Name: fid, Active: true}
		u, creationErr := us.AddNewUser(&User{userData: uData}, "")
		if creationErr != nil {
			log.Println("Error adding new user:", creationErr)
			return nil, creationErr
		}
		user = &userIndexByNameData{ID: u.ID(), Name: u.Name()}
		// user = &(u.(*User).userData) //yep, looks like old C :-), payment for interfaces
	}

	// Nil error means that there already is a user with this federated id.
	// The only possible way for that is faulty creation of the federated accout before.

	fedData := federatedUserID{FederatedID: fid, UserID: user.ID}
	fedInputData, err := dynamodbattribute.MarshalMap(fedData)
	if err != nil {
		log.Println("Error marshalling federated data:", err)
		return nil, ErrorInternalError
	}

	input := &dynamodb.PutItemInput{
		Item:      fedInputData,
		TableName: aws.String(UsersFederatedIDTableName),
	}
	if _, err = us.db.C.PutItem(input); err != nil {
		log.Println("Error putting item:", err)
		return nil, ErrorInternalError
	}
	// just in case
	if user == nil {
		return nil, ErrorInternalError
	}

	udata := userData{ID: user.ID, Name: user.Name, Active: true}
	return &User{userData: udata}, nil
}

// ResetPassword sets new user password.
func (us *UserStorage) ResetPassword(id, password string) error {
	idx, err := xid.FromString(id)
	if err != nil {
		log.Println("Incorrect user ID: ", id)
		return model.ErrorWrongDataFormat
	}

	hash := PasswordHash(password)
	_, err = us.db.C.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(UsersTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(idx.String())},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":p": {S: aws.String(hash)},
		},
		UpdateExpression: aws.String("set pswd = :p"),
		ReturnValues:     aws.String("NONE"),
	})

	return err
}

// IDByName returns userID by name.
func (us *UserStorage) IDByName(name string) (string, error) {
	userIndex, err := us.userIdxByName(name)
	if err != nil {
		return "", err
	}

	user, err := us.UserByID(userIndex.ID)
	if err != nil {
		return "", err
	}

	if !user.Active() {
		return "", ErrorInactiveUser
	}

	return user.ID(), nil
}

// FetchUsers fetches users which name satisfies provided filterString.
// Supports pagination. Search is case-senstive for now.
func (us *UserStorage) FetchUsers(filterString string, skip, limit int) ([]model.User, error) {
	result, err := us.db.C.Query(&dynamodb.QueryInput{
		TableName:              aws.String(UsersTableName),
		IndexName:              aws.String(UserTableEmailIndexName),
		KeyConditionExpression: aws.String("contains(username, :filterStr)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":filterStr": {S: aws.String(filterString)},
		},
		Select: aws.String("ALL_PROJECTED_ATTRIBUTES"),
	})
	if err != nil {
		log.Println("Error querying for users:", err)
		return nil, ErrorInternalError
	}

	users := make([]model.User, len(result.Items))
	for i := 0; i < len(result.Items); i++ {
		user := new(User)
		if err = dynamodbattribute.UnmarshalMap(result.Items[i], user); err != nil {
			log.Println("Error unmarshalling user:", err)
			return nil, ErrorInternalError
		}
		users[i] = user
	}
	return users, nil
}

// ImportJSON imports data from JSON.
func (us *UserStorage) ImportJSON(data []byte) error {
	ud := []userData{}
	if err := json.Unmarshal(data, &ud); err != nil {
		return err
	}
	for _, u := range ud {
		pswd := u.Pswd
		u.Pswd = ""
		if _, err := us.AddNewUser(&User{userData: u}, pswd); err != nil {
			return err
		}
	}
	return nil
}

// userIndexByNameData represents index projected user data.
type userIndexByNameData struct {
	ID   string `json:"id,omitempty"`
	Pswd string `json:"pswd,omitempty"`
	Name string `json:"username,omitempty"`
}

// User data implementation.
type userData struct {
	ID      string                 `json:"id,omitempty"`
	Name    string                 `json:"username,omitempty"`
	Pswd    string                 `json:"pswd,omitempty"`
	Profile map[string]interface{} `json:"profile,omitempty"`
	Active  bool                   `json:"active,omitempty"`
}

// federatedUserID is a struct for mapping federated id to user id.
type federatedUserID struct {
	FederatedID string `json:"federated_id,omitempty"`
	UserID      string `json:"user_id,omitempty"`
}

// User is a user data structure for DynamoDB storage.
type User struct {
	userData
}

// Sanitize removes sensitive data.
func (u *User) Sanitize() {
	u.userData.Pswd = ""
}

// UserFromJSON deserializes user data from JSON.
func UserFromJSON(d []byte) (*User, error) {
	user := userData{}
	if err := json.Unmarshal(d, &user); err != nil {
		log.Println("Error unmarshalling user:", err)
		return &User{}, err
	}
	return &User{userData: user}, nil
}

// ID implements model.User interface.
func (u *User) ID() string { return u.userData.ID }

// Name implements model.User interface.
func (u *User) Name() string { return u.userData.Name }

// PasswordHash implements model.User interface.
func (u *User) PasswordHash() string { return u.userData.Pswd }

// Profile implements model.User interface.
func (u *User) Profile() map[string]interface{} { return u.userData.Profile }

// Active implements model.User interface.
func (u *User) Active() bool { return u.userData.Active }

// PasswordHash creates hash with salt for password.
func PasswordHash(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash)
}

// ensureTable ensures that app storage table exists in the database.
// I'm hiding it in the end of the file, because AWS devs, you are killing me with this API.
func (us *UserStorage) ensureTable() error {
	exists, err := us.db.isTableExists(UsersTableName)
	if err != nil {
		log.Println("Error checking for table existence:", err)
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
					// we are doing local password check.
					Projection: &dynamodb.Projection{
						NonKeyAttributes: []*string{aws.String("pswd"), aws.String("id")},
						ProjectionType:   aws.String("INCLUDE"),
					},
					// Projection: &dynamodb.Projection{
					// 	ProjectionType: aws.String("KEYS_ONLY"),
					// },
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
		if _, err = us.db.C.CreateTable(input); err != nil {
			log.Println("Error creating table:", err)
			return err
		}
	}

	// create table to handle federated ID's
	exists, err = us.db.isTableExists(UsersFederatedIDTableName)
	if err != nil {
		log.Println("Error checking for table existence:", err)
		return err
	}
	if !exists {
		// create table, AWS DynamoDB table creation is overcomplicated for sure
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
		if _, err = us.db.C.CreateTable(input); err != nil {
			log.Println("Error creating table:", err)
			return err
		}
	}
	return nil
}
