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
	// UsersTableName is a table where to store users.
	UsersTableName = "Users"
	// UsersFederatedIDTableName is a table to store federated ids.
	UsersFederatedIDTableName = "UsersByFederatedID"
	// UserTableUsernameIndexName is a user table global index name to access by users by username.
	UserTableUsernameIndexName = "username-index"
	// UsersPhoneNumbersIndexName is a table global index to access users by phone numbers.
	UsersPhoneNumbersIndexName = "phone-index"
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

// NewUser returns pointer to newly created user.
func (us *UserStorage) NewUser() model.User {
	return &User{}
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
		return nil, model.ErrUserNotFound
	}

	userdata := userData{}
	if err = dynamodbattribute.UnmarshalMap(result.Item, &userdata); err != nil {
		log.Println("Error unmarshalling item:", err)
		return nil, ErrorInternalError
	}
	return &User{userData: userdata}, nil
}

// UserByEmail returns user by its email.
func (us *UserStorage) UserByEmail(email string) (model.User, error) {
	// TODO: implement dynamodb UserByEmail
	return nil, errors.New("Not implemented. ")
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
		return "", model.ErrUserNotFound
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
		IndexName:              aws.String(UserTableUsernameIndexName),
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
		log.Println("Error unmarshalling item:", err)
		return nil, ErrorInternalError
	}
	return userdata, nil
}

// userIdxByPhone returns user data projected on the phone index.
func (us *UserStorage) userIdxByPhone(phone string) (*userIndexByPhoneData, error) {
	result, err := us.db.C.Query(&dynamodb.QueryInput{
		TableName:              aws.String(UsersTableName),
		IndexName:              aws.String(UsersPhoneNumbersIndexName),
		KeyConditionExpression: aws.String("phone = :n"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":n": {S: aws.String(phone)},
		},
		Select: aws.String("ALL_PROJECTED_ATTRIBUTES"),
	})
	if err != nil {
		log.Println("Error querying for user by phone number:", err)
		return nil, ErrorInternalError
	}
	if len(result.Items) == 0 {
		return nil, model.ErrUserNotFound
	}

	item := result.Items[0]
	userdata := new(userIndexByPhoneData)
	if err = dynamodbattribute.UnmarshalMap(item, userdata); err != nil {
		log.Println("Error unmarshalling user:", err)
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
		return nil, model.ErrUserNotFound
	}

	user, err := us.UserByID(userIdx.ID)
	if err != nil {
		log.Println("Error querying user by id:", err)
		return nil, ErrorInternalError
	}
	user.Sanitize()
	return user, nil
}

// UserByPhone fetches user by the phone number.
func (us *UserStorage) UserByPhone(phone string) (model.User, error) {
	userIdx, err := us.userIdxByPhone(phone)
	if err != nil {
		log.Println("Error getting user by phone:", err)
		return nil, err
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
	preparedUser, err := us.prepareUserForSaving(usr)
	if err != nil {
		return nil, err
	}

	if len(password) > 0 {
		preparedUser.userData.Pswd = PasswordHash(password)
	}

	updatedUser, err := us.addNewUser(preparedUser)
	return updatedUser, err
}

func (us *UserStorage) prepareUserForSaving(usr model.User) (*User, error) {
	u, ok := usr.(*User)
	if !ok {
		return nil, model.ErrorWrongDataFormat
	}
	// generate new ID if it's not set
	if _, err := xid.FromString(u.ID()); err != nil {
		u.userData.ID = xid.New().String()
	}
	u.userData.Username = strings.ToLower(u.userData.Username)

	return u, nil
}

func (us *UserStorage) addNewUser(u *User) (*User, error) {
	u.userData.NumOfLogins = 0
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
	return u, err
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

// AddUserByNameAndPassword registers new user.
func (us *UserStorage) AddUserByNameAndPassword(name, password string, profile map[string]interface{}) (model.User, error) {
	name = strings.ToLower(name)
	_, err := us.userIdxByName(name)
	if err != nil && err != model.ErrUserNotFound {
		log.Println(err)
		return nil, err
	} else if err == nil {
		return nil, model.ErrorUserExists
	}
	u := userData{Active: true, Username: name, Profile: profile}
	return us.AddNewUser(&User{userData: u}, password)
}

// AddUserWithFederatedID adds new user with social ID.
func (us *UserStorage) AddUserWithFederatedID(provider model.FederatedIdentityProvider, federatedID string) (model.User, error) {
	_, err := us.userIDByFederatedID(provider, federatedID)
	if err != nil && err != model.ErrUserNotFound {
		log.Println("Error getting user by name:", err)
		return nil, err
	} else if err == nil {
		return nil, model.ErrorUserExists
	}

	fid := string(provider) + ":" + federatedID

	user, err := us.userIdxByName(fid)
	if err != nil && err != model.ErrUserNotFound {
		log.Println("Error getting user by name:", err)
		return nil, err
	} else if err == model.ErrUserNotFound {
		// no such user, let's create it
		uData := userData{Username: fid, Active: true}
		u, creationErr := us.AddNewUser(&User{userData: uData}, "")
		if creationErr != nil {
			log.Println("Error adding new user:", creationErr)
			return nil, creationErr
		}
		user = &userIndexByNameData{ID: u.ID(), Username: u.Username()}
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

	udata := userData{ID: user.ID, Username: user.Username, Active: true}
	return &User{userData: udata}, nil
}

// AddUserByPhone registers new user with phone number.
func (us *UserStorage) AddUserByPhone(phone string) (model.User, error) {
	_, err := us.userIdxByPhone(phone)
	if err != nil && err != model.ErrUserNotFound {
		log.Println(err)
		return nil, err
	} else if err == nil {
		return nil, model.ErrorUserExists
	}
	u := userData{Active: true, Phone: phone, ID: xid.New().String()}
	return us.AddNewUser(&User{userData: u}, "")
}

// UpdateUser updates user in DynamoDB storage.
func (us *UserStorage) UpdateUser(userID string, newUser model.User) (model.User, error) {
	if _, err := xid.FromString(userID); err != nil {
		log.Println("Incorrect userID: ", userID)
		return nil, model.ErrorWrongDataFormat
	}

	res, ok := newUser.(*User)
	if !ok || res == nil {
		return nil, model.ErrorWrongDataFormat
	}

	// use ID from the request if it's not set
	if len(res.ID()) == 0 {
		res.userData.ID = userID
	}

	if err := us.DeleteUser(userID); err != nil {
		log.Println("Error deleting old user:", err)
		return nil, err
	}

	preparedUser, err := us.prepareUserForSaving(res)
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

// ResetUsername sets user username.
func (us *UserStorage) ResetUsername(id, username string) error {
	idx, err := xid.FromString(id)
	if err != nil {
		log.Println("Incorrect user ID: ", id)
		return model.ErrorWrongDataFormat
	}

	_, err = us.db.C.UpdateItem(&dynamodb.UpdateItemInput{

		TableName: aws.String(UsersTableName),
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

	if !user.Active() {
		return "", ErrorInactiveUser
	}

	return user.ID(), nil
}

// FetchUsers fetches users which name satisfies provided filterString.
// Supports pagination. Search is case-senstive for now.
func (us *UserStorage) FetchUsers(filterString string, skip, limit int) ([]model.User, int, error) {
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(UsersTableName),
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
		return []model.User{}, 0, ErrorInternalError
	}

	users := make([]model.User, len(result.Items))
	for i := 0; i < len(result.Items); i++ {
		if i < skip {
			continue // TODO: use internal pagination mechanism
		}
		user := new(User)
		if err = dynamodbattribute.UnmarshalMap(result.Items[i], user); err != nil {
			log.Println("Error unmarshalling user:", err)
			return []model.User{}, 0, ErrorInternalError
		}
		users[i] = user
	}
	return users, len(result.Items), nil
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
		TableName: aws.String(UsersTableName),
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
	exists, err := us.db.IsTableExists(UsersTableName)
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
					IndexName: aws.String(UserTableUsernameIndexName),
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
					IndexName: aws.String(UsersPhoneNumbersIndexName),
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
			TableName:   aws.String(UsersTableName),
		}
		if _, err = us.db.C.CreateTable(input); err != nil {
			log.Println("Error creating table:", err)
			return err
		}
	}

	// create table to handle federated ID's
	exists, err = us.db.IsTableExists(UsersFederatedIDTableName)
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
			TableName:   aws.String(UsersFederatedIDTableName),
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
