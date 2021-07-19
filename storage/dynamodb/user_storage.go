package dynamodb

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/madappgang/identifo/model"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

const (
	usersTableName             = "Users"              // usersTableName is a table where to store users.
	usersFederatedIDTableName  = "UsersByFederatedID" // usersFederatedIDTableName is a table to store federated ids.
	userTableUsernameIndexName = "username-index"     // userTableUsernameIndexName is a user table global index name to access by users by username.
	usersPhoneNumbersIndexName = "phone-index"        // usersPhoneNumbersIndexName is a table global index to access users by phone numbers.
)

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
		return model.User{}, model.ErrorWrongDataFormat
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
		return model.User{}, ErrorInternalError
	}
	if result.Item == nil {
		return model.User{}, model.ErrUserNotFound
	}

	userdata := model.User{}
	if err = dynamodbattribute.UnmarshalMap(result.Item, &userdata); err != nil {
		log.Println("error while unmarshal item: ", err)
		return model.User{}, ErrorInternalError
	}
	return userdata, nil
}

// UserByEmail returns user by its email.
func (us *UserStorage) UserByEmail(email string) (model.User, error) {
	// TODO: implement dynamodb UserByEmail
	// clear password hash
	// u.Pswd = ""
	return model.User{}, errors.New("not implemented")
}

func (us *UserStorage) userIDByFederatedID(provider model.FederatedIdentityProvider, id string) (string, error) {
	fid := string(provider) + ":" + id
	result, err := us.db.C.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(usersFederatedIDTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"federated_id": {
				S: aws.String(fid),
			},
		},
	})
	if err != nil {
		log.Println("error getting item from DynamoDB: ", err)
		return "", ErrorInternalError
	}
	if result.Item == nil {
		return "", model.ErrUserNotFound
	}

	fedData := federatedUserID{}
	if err = dynamodbattribute.UnmarshalMap(result.Item, &fedData); err != nil || len(fedData.UserID) == 0 {
		log.Println("error while unmarshal item: ", err)
		return "", ErrorInternalError
	}

	return fedData.UserID, nil
}

// UserByFederatedID returns user by federated ID.
func (us *UserStorage) UserByFederatedID(provider model.FederatedIdentityProvider, id string) (model.User, error) {
	userID, err := us.userIDByFederatedID(provider, id)
	if err != nil {
		return model.User{}, err
	}
	if len(userID) == 0 {
		return model.User{}, model.ErrorWrongDataFormat
	}
	u, err := us.UserByID(userID)
	// clear password hash
	u.Pswd = ""
	return u, err
}

// UserExists checks if user with provided name exists.
func (us *UserStorage) UserExists(name string) bool {
	_, err := us.userIdxByName(name)
	return err == nil
}

// AttachDeviceToken do nothing here
// TODO: implement device storage
func (us *UserStorage) AttachDeviceToken(id, token string) error {
	// we are not supporting devices for users here
	return model.ErrorNotImplemented
}

// DetachDeviceToken do nothing here yet
// TODO: implement
func (us *UserStorage) DetachDeviceToken(token string) error {
	return model.ErrorNotImplemented
}

// RequestScopes for now returns requested scope
// TODO: implement scope logic
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
		return nil, ErrorInternalError
	}
	if len(result.Items) == 0 {
		return nil, model.ErrUserNotFound
	}

	item := result.Items[0]
	userdata := new(userIndexByNameData)
	if err = dynamodbattribute.UnmarshalMap(item, userdata); err != nil {
		log.Println("error while unmarshal item:", err)
		return nil, ErrorInternalError
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
		log.Println("error querying for user by phone number: ", err)
		return nil, ErrorInternalError
	}
	if len(result.Items) == 0 {
		return nil, model.ErrUserNotFound
	}

	item := result.Items[0]
	userdata := new(userIndexByPhoneData)
	if err = dynamodbattribute.UnmarshalMap(item, userdata); err != nil {
		log.Println("error while unmarshal user: ", err)
		return nil, ErrorInternalError
	}
	return userdata, nil
}

// UserByUsername returns user by name
func (us *UserStorage) UserByUsername(username string) (model.User, error) {
	username = strings.ToLower(username)
	userIdx, err := us.userIdxByName(username)
	if err != nil {
		log.Println("error getting user by name: ", err)
		return model.User{}, err
	}

	user, err := us.UserByID(userIdx.ID)
	if err != nil {
		log.Println("error querying user by id: ", err)
		return model.User{}, ErrorInternalError
	}
	// clear password hash
	user.Pswd = ""
	return user, nil
}

// UserByPhone fetches user by the phone number.
func (us *UserStorage) UserByPhone(phone string) (model.User, error) {
	userIdx, err := us.userIdxByPhone(phone)
	if err != nil {
		log.Println("error getting user by phone: ", err)
		return model.User{}, err
	}

	user, err := us.UserByID(userIdx.ID)
	if err != nil {
		log.Println("error querying user by id: ", err)
		return model.User{}, ErrorInternalError
	}
	// clear password hash
	user.Pswd = ""
	return user, nil
}

// AddNewUser adds new user.
func (us *UserStorage) AddNewUser(usr model.User, password string) (model.User, error) {
	preparedUser, err := us.prepareUserForSaving(usr)
	if err != nil {
		return model.User{}, err
	}

	if len(password) > 0 {
		preparedUser.Pswd = model.PasswordHash(password)
	}

	updatedUser, err := us.addNewUser(preparedUser)
	return updatedUser, err
}

func (us *UserStorage) prepareUserForSaving(usr model.User) (model.User, error) {
	// Generate new ID if it's not set.
	if _, err := xid.FromString(usr.ID); err != nil {
		usr.ID = xid.New().String()
	}
	usr.Username = strings.ToLower(usr.Username)

	return usr, nil
}

func (us *UserStorage) addNewUser(u model.User) (model.User, error) {
	u.NumOfLogins = 0
	uv, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		log.Println("error marshalling user: ", err)
		return model.User{}, ErrorInternalError
	}

	input := &dynamodb.PutItemInput{
		Item:      uv,
		TableName: aws.String(usersTableName),
	}
	if _, err = us.db.C.PutItem(input); err != nil {
		log.Println("error putting item: ", err)
		return model.User{}, ErrorInternalError
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

// AddUserWithPassword creates new user and saves it in the database.
func (us *UserStorage) AddUserWithPassword(user model.User, password, role string, isAnonymous bool) (model.User, error) {
	if _, err := us.UserByUsername(user.Username); err == nil {
		return model.User{}, model.ErrorUserExists
	}
	if _, err := us.UserByEmail(user.Email); err == nil {
		return model.User{}, model.ErrorUserExists
	}
	if _, err := us.UserByPhone(user.Phone); err == nil {
		return model.User{}, model.ErrorUserExists
	}

	u := model.User{
		ID:         xid.New().String(),
		Active:     true,
		Username:   user.Username,
		Phone:      user.Phone,
		Email:      user.Email,
		AccessRole: role,
		Anonymous:  isAnonymous,
	}

	return us.AddNewUser(u, password)
}

// AddUserWithFederatedID adds new user with social ID.
func (us *UserStorage) AddUserWithFederatedID(provider model.FederatedIdentityProvider, federatedID, role string) (model.User, error) {
	_, err := us.userIDByFederatedID(provider, federatedID)
	if err != nil && err != model.ErrUserNotFound {
		log.Println("error getting user by name: ", err)
		return model.User{}, err
	} else if err == nil {
		return model.User{}, model.ErrorUserExists
	}

	fid := string(provider) + ":" + federatedID

	user, err := us.userIdxByName(fid)
	if err != nil && err != model.ErrUserNotFound {
		log.Println("error getting user by name: ", err)
		return model.User{}, err
	} else if err == model.ErrUserNotFound {
		// no such user, let's create it
		uData := model.User{Username: fid, AccessRole: role, Active: true}
		u, creationErr := us.AddNewUser(uData, "")
		if creationErr != nil {
			log.Println("error adding new user: ", creationErr)
			return model.User{}, creationErr
		}
		user = &userIndexByNameData{ID: u.ID, Username: u.Username}
	}

	// Nil error means that there already is a user with this federated id.
	// The only possible way for that is faulty creation of the federated accout before.

	fedData := federatedUserID{FederatedID: fid, UserID: user.ID}
	fedInputData, err := dynamodbattribute.MarshalMap(fedData)
	if err != nil {
		log.Println("error marshalling federated data: ", err)
		return model.User{}, ErrorInternalError
	}

	input := &dynamodb.PutItemInput{
		Item:      fedInputData,
		TableName: aws.String(usersFederatedIDTableName),
	}
	if _, err = us.db.C.PutItem(input); err != nil {
		log.Println("error putting item: ", err)
		return model.User{}, ErrorInternalError
	}
	// just in case
	if user == nil {
		return model.User{}, ErrorInternalError
	}

	udata := model.User{ID: user.ID, Username: user.Username, Active: true}
	return udata, nil
}

// UpdateUser updates user in DynamoDB storage.
func (us *UserStorage) UpdateUser(userID string, user model.User) (model.User, error) {
	if _, err := xid.FromString(userID); err != nil {
		log.Println("incorrect userID: ", userID)
		return model.User{}, model.ErrorWrongDataFormat
	}

	// use ID from the request if it's not set
	if len(user.ID) == 0 {
		user.ID = userID
	}

	if err := us.DeleteUser(userID); err != nil {
		log.Println("error deleting old user: ", err)
		return model.User{}, err
	}

	preparedUser, err := us.prepareUserForSaving(user)
	if err != nil {
		return model.User{}, err
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

	hash := model.PasswordHash(password)
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

// CheckPassword check that password is valid for user id.
func (us *UserStorage) CheckPassword(id, password string) error {
	user, err := us.UserByID(id)
	if err != nil {
		return model.ErrUserNotFound
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Pswd), []byte(password)); err != nil {
		// return this error to hide the existence of the user.
		return model.ErrUserNotFound
	}
	return nil
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

	if !user.Active {
		return "", ErrorInactiveUser
	}

	return user.ID, nil
}

// FetchUsers fetches users which name satisfies provided filterString.
// Supports pagination. Search is case-sensitive for now.
func (us *UserStorage) FetchUsers(filterString string, skip, limit int) ([]model.User, int, error) {
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
		log.Println("error querying for users: ", err)
		return []model.User{}, 0, ErrorInternalError
	}

	users := make([]model.User, len(result.Items))
	for i := 0; i < len(result.Items); i++ {
		if i < skip {
			continue // TODO: use internal pagination mechanism
		}
		user := model.User{}
		if err = dynamodbattribute.UnmarshalMap(result.Items[i], &user); err != nil {
			log.Println("error while unmarshal user: ", err)
			return []model.User{}, 0, ErrorInternalError
		}
		users[i] = user
	}
	return users, len(result.Items), nil
}

// ImportJSON imports data from JSON.
func (us *UserStorage) ImportJSON(data []byte) error {
	ud := []model.User{}
	if err := json.Unmarshal(data, &ud); err != nil {
		return err
	}
	for _, u := range ud {
		pswd := u.Pswd
		u.Pswd = ""
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
