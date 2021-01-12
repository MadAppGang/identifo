package dynamodb

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/hashicorp/go-plugin"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/plugin/shared"
	"github.com/madappgang/identifo/proto"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

const (
	// usersTableName is a table where to store users.
	usersTableName = "Users"
	// usersFederatedIDTableName is a table to store federated ids.
	usersFederatedIDTableName = "UsersByFederatedID"
	// userTableUsernameIndexName is a user table global index name to access by users by username.
	userTableUsernameIndexName = "username-index"
	// usersPhoneNumbersIndexName is a table global index to access users by phone numbers.
	usersPhoneNumbersIndexName = "phone-index"
)

func main() {
	endpoint := os.Getenv("DB_ENDPOINT")
	if endpoint == "" {
		panic("Empty DB_ENDPOINT")
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		panic("Empty AWS_REGION")
	}

	db, err := NewDB(endpoint, region)
	if err != nil {
		panic(err)
	}

	us, err := NewUserStorage(db)
	if err != nil {
		panic(err)
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"user_storage": &shared.UserStorageGRPCPlugin{
				Impl: us,
			},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

// userIndexByNameData represents username index projected user data.
type userIndexByNameData struct {
	ID       string `json:"id,omitempty"`
	Pswd     string `json:"pswd,omitempty"`
	Username string `json:"username,omitempty"`
}

// userIndexByPhoneData represents phone index projected user data.
type userIndexByPhoneData struct {
	ID    string `json:"id,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// federatedUserID is a struct for mapping federated id to user id.
type federatedUserID struct {
	FederatedID string `json:"federated_id,omitempty"`
	UserID      string `json:"user_id,omitempty"`
}

// NewDB creates new database connection.
func NewDB(endpoint string, region string) (*DB, error) {
	config := &aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	}
	sess, err := session.NewSession(config)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &DB{C: dynamodb.New(sess)}, nil
}

// DB represents connection to AWS DynamoDB service or local instance.
type DB struct {
	C *dynamodb.DynamoDB
}

// IsTableExists checks if table exists.
func (db *DB) IsTableExists(table string) (bool, error) {
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.C.DescribeTableWithContext(timeoutCtx, input)
	if AwsErrorErrorNotFound(err) {
		return false, nil
		//if table not exists - create table
	}
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

// AwsErrorErrorNotFound checks if error has type dynamodb.ErrCodeResourceNotFoundException.
func AwsErrorErrorNotFound(err error) bool {
	if err == nil {
		return false
	}
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
			return true
		}
	}
	return false
}

// NewUserStorage creates and provisions new user storage instance.
func NewUserStorage(db *DB) (shared.UserStorage, error) {
	us := &UserStorage{db: db}
	err := us.ensureTable()
	return us, err
}

// UserStorage stores and manages data in DynamoDB storage.
type UserStorage struct {
	db *DB
}

// UserByID returns user by its ID.
func (us *UserStorage) UserByID(id string) (*proto.User, error) {
	idx, err := xid.FromString(id)
	if err != nil {
		log.Println("Incorrect user ID: ", id)
		return nil, model.ErrorWrongDataFormat
	}

	result, err := us.db.C.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(usersTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(idx.String()),
			},
		},
	})
	if err != nil {
		log.Println("Error getting item from DynamoDB:", err)
		return nil, shared.ErrorInternalError
	}
	if result.Item == nil {
		return nil, shared.ErrUserNotFound
	}

	u := proto.User{}
	if err = dynamodbattribute.UnmarshalMap(result.Item, &u); err != nil {
		log.Println("Error unmarshalling item:", err)
		return nil, shared.ErrorInternalError
	}
	return &u, nil
}

// UserByEmail returns user by its email.
func (us *UserStorage) UserByEmail(email string) (*proto.User, error) {
	// TODO: implement dynamodb UserByEmail
	return nil, errors.New("Not implemented. ")
}

func (us *UserStorage) userIDByFederatedID(provider proto.FederatedIdentityProvider, id string) (string, error) {
	fid := provider.String() + ":" + id
	result, err := us.db.C.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(usersFederatedIDTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"federated_id": {
				S: aws.String(fid),
			},
		},
	})
	if err != nil {
		log.Println("Error getting item from DynamoDB:", err)
		return "", shared.ErrorInternalError
	}
	if result.Item == nil {
		return "", shared.ErrUserNotFound
	}

	fedData := federatedUserID{}
	if err = dynamodbattribute.UnmarshalMap(result.Item, &fedData); err != nil || len(fedData.UserID) == 0 {
		log.Println("Error unmarshalling item:", err)
		return "", shared.ErrorInternalError
	}
	return fedData.UserID, nil
}

// UserByFederatedID returns user by federated ID.
func (us *UserStorage) UserByFederatedID(provider proto.FederatedIdentityProvider, id string) (*proto.User, error) {
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
		TableName:              aws.String(usersTableName),
		IndexName:              aws.String(userTableUsernameIndexName),
		KeyConditionExpression: aws.String("username = :n"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":n": {S: aws.String(name)},
		},
		Select: aws.String("ALL_PROJECTED_ATTRIBUTES"), // retrieve all attributes, because we need to make local check.
	})
	if err != nil {
		log.Println("Error querying for items:", err)
		return nil, shared.ErrorInternalError
	}
	if len(result.Items) == 0 {
		return nil, shared.ErrUserNotFound
	}

	item := result.Items[0]
	userdata := new(userIndexByNameData)
	if err = dynamodbattribute.UnmarshalMap(item, userdata); err != nil {
		log.Println("Error unmarshalling item:", err)
		return nil, shared.ErrorInternalError
	}
	return userdata, nil
}

// userIdxByPhone returns user data projected on the phone index.
func (us *UserStorage) userIdxByPhone(phone string) (*userIndexByPhoneData, error) {
	result, err := us.db.C.Query(&dynamodb.QueryInput{
		TableName:              aws.String(usersTableName),
		IndexName:              aws.String(usersPhoneNumbersIndexName),
		KeyConditionExpression: aws.String("phone = :n"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":n": {S: aws.String(phone)},
		},
		Select: aws.String("ALL_PROJECTED_ATTRIBUTES"),
	})
	if err != nil {
		log.Println("Error querying for user by phone number:", err)
		return nil, shared.ErrorInternalError
	}
	if len(result.Items) == 0 {
		return nil, shared.ErrUserNotFound
	}

	item := result.Items[0]
	userdata := new(userIndexByPhoneData)
	if err = dynamodbattribute.UnmarshalMap(item, userdata); err != nil {
		log.Println("Error unmarshalling user:", err)
		return nil, shared.ErrorInternalError
	}
	return userdata, nil
}

// UserByNamePassword returns user by name and password.
func (us *UserStorage) UserByNamePassword(name, password string) (*proto.User, error) {
	name = strings.ToLower(name)
	userIdx, err := us.userIdxByName(name)
	if err != nil {
		log.Println("Error getting user by name:", err)
		return nil, err
	}
	// if password is incorrect, return 'not found' error for security reasons.
	if bcrypt.CompareHashAndPassword([]byte(userIdx.Pswd), []byte(password)) != nil {
		return nil, shared.ErrUserNotFound
	}

	user, err := us.UserByID(userIdx.ID)
	if err != nil {
		log.Println("Error querying user by id:", err)
		return nil, shared.ErrorInternalError
	}
	return user, nil
}

// UserByPhone fetches user by the phone number.
func (us *UserStorage) UserByPhone(phone string) (*proto.User, error) {
	userIdx, err := us.userIdxByPhone(phone)
	if err != nil {
		log.Println("Error getting user by phone:", err)
		return nil, err
	}

	user, err := us.UserByID(userIdx.ID)
	if err != nil {
		log.Println("Error querying user by id:", err)
		return nil, shared.ErrorInternalError
	}

	return user, nil
}

// AddNewUser adds new user.
func (us *UserStorage) AddNewUser(usr *proto.User, password string) (*proto.User, error) {
	preparedUser, err := us.prepareUserForSaving(usr)
	if err != nil {
		return nil, err
	}

	if len(password) > 0 {
		preparedUser.PasswordHash = PasswordHash(password)
	}

	updatedUser, err := us.addNewUser(preparedUser)
	return updatedUser, err
}

func (us *UserStorage) prepareUserForSaving(u *proto.User) (*proto.User, error) {
	// Generate new ID if it's not set.
	if _, err := xid.FromString(u.Id); err != nil {
		u.Id = xid.New().String()
	}
	u.Username = strings.ToLower(u.Username)

	return u, nil
}

func (us *UserStorage) addNewUser(u *proto.User) (*proto.User, error) {
	u.NumOfLogins = 0
	uv, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		log.Println("Error marshalling user:", err)
		return nil, shared.ErrorInternalError
	}

	input := &dynamodb.PutItemInput{
		Item:      uv,
		TableName: aws.String(usersTableName),
	}
	if _, err = us.db.C.PutItem(input); err != nil {
		log.Println("Error putting item:", err)
		return nil, shared.ErrorInternalError
	}
	return u, err
}

// DeleteUser deletes user by id.
func (us *UserStorage) DeleteUser(id string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
		TableName: aws.String(usersTableName),
	}
	_, err := us.db.C.DeleteItem(input)
	return err
}

// AddUserByNameAndPassword registers new user.
func (us *UserStorage) AddUserByNameAndPassword(username, password, role string, isAnonymous bool) (*proto.User, error) {
	username = strings.ToLower(username)
	_, err := us.userIdxByName(username)
	if err != nil && err != shared.ErrUserNotFound {
		log.Println(err)
		return nil, err
	} else if err == nil {
		return nil, model.ErrorUserExists
	}

	u := proto.User{
		IsActive:    true,
		Username:    username,
		AccessRole:  role,
		IsAnonymous: isAnonymous,
	}
	if shared.EmailRegexp.MatchString(u.Username) {
		u.Email = u.Username
	}
	if shared.PhoneRegexp.MatchString(u.Username) {
		u.Phone = u.Username
	}

	return us.AddNewUser(&u, password)
}

// AddUserWithFederatedID adds new user with social ID.
func (us *UserStorage) AddUserWithFederatedID(provider proto.FederatedIdentityProvider, federatedID, role string) (*proto.User, error) {
	_, err := us.userIDByFederatedID(provider, federatedID)
	if err != nil && err != shared.ErrUserNotFound {
		log.Println("Error getting user by name:", err)
		return nil, err
	} else if err == nil {
		return nil, model.ErrorUserExists
	}

	fid := provider.String() + ":" + federatedID

	user, err := us.userIdxByName(fid)
	if err != nil && err != shared.ErrUserNotFound {
		log.Println("Error getting user by name:", err)
		return nil, err
	} else if err == shared.ErrUserNotFound {
		// no such user, let's create it
		u := &proto.User{Username: fid, AccessRole: role, IsActive: true}
		u, creationErr := us.AddNewUser(u, "")
		if creationErr != nil {
			log.Println("Error adding new user:", creationErr)
			return nil, creationErr
		}
		user = &userIndexByNameData{ID: u.Id, Username: u.Username}
		// user = &(u.(*User).userData) //yep, looks like old C :-), payment for interfaces
	}

	// Nil error means that there already is a user with this federated id.
	// The only possible way for that is faulty creation of the federated accout before.

	fedData := federatedUserID{FederatedID: fid, UserID: user.ID}
	fedInputData, err := dynamodbattribute.MarshalMap(fedData)
	if err != nil {
		log.Println("Error marshalling federated data:", err)
		return nil, shared.ErrorInternalError
	}

	input := &dynamodb.PutItemInput{
		Item:      fedInputData,
		TableName: aws.String(usersFederatedIDTableName),
	}
	if _, err = us.db.C.PutItem(input); err != nil {
		log.Println("Error putting item:", err)
		return nil, shared.ErrorInternalError
	}
	// just in case
	if user == nil {
		return nil, shared.ErrorInternalError
	}
	return &proto.User{Id: user.ID, Username: user.Username, IsActive: true}, nil
}

// AddUserByPhone registers new user with phone number.
func (us *UserStorage) AddUserByPhone(phone, role string) (*proto.User, error) {
	_, err := us.userIdxByPhone(phone)
	if err != nil && err != shared.ErrUserNotFound {
		log.Println(err)
		return nil, err
	} else if err == nil {
		return nil, model.ErrorUserExists
	}

	u := &proto.User{
		Id:          xid.New().String(),
		Username:    phone,
		IsActive:    true,
		Phone:       phone,
		AccessRole:  role,
		NumOfLogins: 0,
	}
	return us.AddNewUser(u, "")
}

// UpdateUser updates user in DynamoDB storage.
func (us *UserStorage) UpdateUser(userID string, newUser *proto.User) (*proto.User, error) {
	if _, err := xid.FromString(userID); err != nil {
		log.Println("Incorrect userID: ", userID)
		return nil, model.ErrorWrongDataFormat
	}

	// use ID from the request if it's not set
	if len(newUser.Id) == 0 {
		newUser.Id = userID
	}

	if err := us.DeleteUser(userID); err != nil {
		log.Println("Error deleting old user:", err)
		return nil, err
	}

	preparedUser, err := us.prepareUserForSaving(newUser)
	if err != nil {
		return nil, err
	}

	updatedUser, err := us.addNewUser(preparedUser)
	return updatedUser, err
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
		TableName: aws.String(usersTableName),
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

// ResetUsername sets user username.
func (us *UserStorage) ResetUsername(id, username string) error {
	idx, err := xid.FromString(id)
	if err != nil {
		log.Println("Incorrect user ID: ", id)
		return model.ErrorWrongDataFormat
	}

	_, err = us.db.C.UpdateItem(&dynamodb.UpdateItemInput{

		TableName: aws.String(usersTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(idx.String())},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":u": {S: aws.String(username)},
		},
		UpdateExpression: aws.String("set username = :u"),
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

	if !user.IsActive {
		return "", errors.New("User is inactive")
	}
	return user.Id, nil
}

// FetchUsers fetches users which name satisfies provided filterString.
// Supports pagination. Search is case-senstive for now.
func (us *UserStorage) FetchUsers(filterString string, skip, limit int) ([]*proto.User, int, error) {
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(usersTableName),
		Limit:     aws.Int64(int64(limit)),
	}

	if len(filterString) != 0 {
		scanInput.FilterExpression = aws.String("contains(username, :filterStr)")
		scanInput.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{
			":filterStr": {S: aws.String(filterString)},
		}
	}

	result, err := us.db.C.Scan(scanInput)
	if err != nil {
		log.Println("Error querying for users:", err)
		return []*proto.User{}, 0, shared.ErrorInternalError
	}

	users := make([]*proto.User, len(result.Items))
	for i := 0; i < len(result.Items); i++ {
		if i < skip {
			continue // TODO: use internal pagination mechanism
		}
		user := new(proto.User)
		if err = dynamodbattribute.UnmarshalMap(result.Items[i], user); err != nil {
			log.Println("Error unmarshalling user:", err)
			return []*proto.User{}, 0, shared.ErrorInternalError
		}
		users[i] = user
	}
	return users, len(result.Items), nil
}

// ImportJSON imports data from JSON.
func (us *UserStorage) ImportJSON(data []byte) error {
	ud := []*proto.User{}
	if err := json.Unmarshal(data, &ud); err != nil {
		return err
	}
	for _, u := range ud {
		pswd := u.PasswordHash
		u.PasswordHash = ""
		if _, err := us.AddNewUser(u, pswd); err != nil {
			return err
		}
	}
	return nil
}

// UpdateLoginMetadata updates user's login metadata.
func (us *UserStorage) UpdateLoginMetadata(userID string) {
	if _, err := xid.FromString(userID); err != nil {
		log.Println("Incorrect userID: ", userID)
		return
	}

	if _, err := us.UserByID(userID); err != nil {
		log.Println("Cannot get user by ID: ", userID)
		return
	}

	_, err := us.db.C.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(usersTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(userID)},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":now": {N: aws.String(strconv.Itoa(int(time.Now().Unix())))},
			":one": {N: aws.String("1")},
		},
		UpdateExpression: aws.String("set latest_login_time :now add num_of_logins :one "),
		ReturnValues:     aws.String("NONE"),
	})
	if err != nil {
		log.Printf("Cannot update login metadata of user %s: %s\n", userID, err)
		return
	}
}

// ensureTable ensures that user storage table exists in the database.
// I'm hiding it in the end of the file, because AWS devs, you are killing me with this API.
func (us *UserStorage) ensureTable() error {
	exists, err := us.db.IsTableExists(usersTableName)
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
				{
					AttributeName: aws.String("phone"),
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
					IndexName: aws.String(userTableUsernameIndexName),
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
				},
				{
					IndexName: aws.String(usersPhoneNumbersIndexName),
					KeySchema: []*dynamodb.KeySchemaElement{
						{
							AttributeName: aws.String("phone"),
							KeyType:       aws.String("HASH"),
						},
					},
					Projection: &dynamodb.Projection{
						NonKeyAttributes: []*string{aws.String("pswd"), aws.String("id")},
						ProjectionType:   aws.String("INCLUDE"),
					},
				},
			},
			BillingMode: aws.String("PAY_PER_REQUEST"),
			TableName:   aws.String(usersTableName),
		}
		if _, err = us.db.C.CreateTable(input); err != nil {
			log.Println("Error creating table:", err)
			return err
		}
	}

	// create table to handle federated ID's
	exists, err = us.db.IsTableExists(usersFederatedIDTableName)
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
			BillingMode: aws.String("PAY_PER_REQUEST"),
			TableName:   aws.String(usersFederatedIDTableName),
		}
		if _, err = us.db.C.CreateTable(input); err != nil {
			log.Println("Error creating table:", err)
			return err
		}
	}
	return nil
}

// Close does nothing here.
func (us *UserStorage) Close() {}

// PasswordHash creates hash with salt for password.
func PasswordHash(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash)
}
