package dynamodb

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/rs/xid"
)

const (
	invitesTableName     = "Invites"
	inviteEmailIndexName = "invite-email"
	maxInvitesLimit      = 20
)

type inviteIndexByEmailData struct {
	ID    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
}

// InviteStorage is a DynamoDB invite storage.
type InviteStorage struct {
	db *DB
}

// NewInviteStorage creates new DynamoDB invite storage.
func NewInviteStorage(settings model.DynamoDatabaseSettings) (*InviteStorage, error) {
	if len(settings.Endpoint) == 0 || len(settings.Region) == 0 {
		return nil, l.ErrorAPIDataError
	}

	// create database
	db, err := NewDB(settings.Endpoint, settings.Region)
	if err != nil {
		return nil, err
	}

	is := &InviteStorage{db: db}
	err = is.ensureTable()
	return is, err
}

// ensureTable ensures that invite storage exists in the database.
func (is *InviteStorage) ensureTable() error {
	exists, err := is.db.IsTableExists(invitesTableName)
	if err != nil {
		log.Println("Error checking Invites table existence:", err)
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
			{
				AttributeName: aws.String("email"),
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
				IndexName: aws.String(inviteEmailIndexName),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("email"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("KEYS_ONLY"),
				},
			},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
		TableName:   aws.String(invitesTableName),
	}

	_, err = is.db.C.CreateTable(input)
	return err
}

// Save creates and saves new invite to a database.
func (is *InviteStorage) Save(email, inviteToken, role, appID, createdBy string, expiresAt time.Time) error {
	invite := model.Invite{
		ID:        xid.New().String(),
		AppID:     appID,
		Token:     inviteToken,
		Archived:  false,
		Email:     email,
		Role:      role,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	iv, err := dynamodbattribute.MarshalMap(invite)
	if err != nil {
		log.Println("error marshalling invite: ", err)
		return l.ErrorAPIDataError
	}

	input := &dynamodb.PutItemInput{
		Item:      iv,
		TableName: aws.String(invitesTableName),
	}

	if _, err = is.db.C.PutItem(input); err != nil {
		log.Println("error putting invite to storage: ", err)
		return l.ErrorAPIDataError
	}
	return nil
}

// inviteIdxByEmail returns invite data projected on the email index.
func (is *InviteStorage) inviteIdxByEmail(email string) (*inviteIndexByEmailData, error) {
	result, err := is.db.C.Query(&dynamodb.QueryInput{
		TableName:              aws.String(invitesTableName),
		IndexName:              aws.String(inviteEmailIndexName),
		KeyConditionExpression: aws.String("email = :n"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":n": {S: aws.String(email)},
		},
		Select: aws.String("ALL_PROJECTED_ATTRIBUTES"),
	})
	if err != nil {
		log.Println("error querying for invite by email: ", err)
		return nil, l.ErrorAPIDataError
	}
	if len(result.Items) == 0 {
		return nil, l.NewError(l.ErrorNotFound, "invite")
	}

	item := result.Items[0]
	inviteData := new(inviteIndexByEmailData)
	if err = dynamodbattribute.UnmarshalMap(item, inviteData); err != nil {
		log.Println("error while unmarshal invite: ", err)
		return nil, l.NewError(l.ErrorNotFound, "invite")
	}
	return inviteData, nil
}

// GetByEmail returns not archived and not expired invite by email.
func (is *InviteStorage) GetByEmail(email string) (model.Invite, error) {
	inviteIdx, err := is.inviteIdxByEmail(email)
	if err != nil {
		log.Println("error getting invite by email: ", err)
		return model.Invite{}, err
	}

	invite, err := is.GetByID(inviteIdx.ID)
	if err != nil {
		log.Println("error querying invite by id: ", err)
		return model.Invite{}, l.NewError(l.ErrorNotFound, "invite")
	}

	return invite, nil
}

// GetByID returns invite by its ID.
func (is *InviteStorage) GetByID(id string) (model.Invite, error) {
	if len(id) == 0 {
		return model.Invite{}, l.NewError(l.ErrorNotFound, "invite")
	}

	result, err := is.db.C.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(invitesTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		log.Println("Error getting invite:", err)
		return model.Invite{}, l.ErrorAPIDataError
	}

	if result.Item == nil {
		return model.Invite{}, l.NewError(l.ErrorNotFound, "invite")
	}

	invite := model.Invite{}
	if err = dynamodbattribute.UnmarshalMap(result.Item, &invite); err != nil {
		log.Println("Error unmarshalling invite:", err)
		return model.Invite{}, l.ErrorAPIDataError
	}
	return invite, nil
}

// GetAll returns all active invites by default.
// To get an archived invites need to set withArchived argument to true.
func (is *InviteStorage) GetAll(withArchived bool, skip, limit int) ([]model.Invite, int, error) {
	if limit == 0 || limit > maxInvitesLimit {
		limit = maxInvitesLimit
	}

	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(invitesTableName),
		Limit:     aws.Int64(int64(limit)),
	}

	if !withArchived {
		scanInput.FilterExpression = aws.String("#archived = :archived")
		scanInput.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{
			":archived": {BOOL: aws.Bool(false)},
		}
		scanInput.ExpressionAttributeNames = map[string]*string{
			"#archived": aws.String("archived"),
		}
	}

	result, err := is.db.C.Scan(scanInput)
	if err != nil {
		log.Println("Error querying for invites:", err)
		return []model.Invite{}, 0, l.ErrorAPIDataError
	}

	invites := make([]model.Invite, len(result.Items))
	for i := 0; i < len(result.Items); i++ {
		if i < skip {
			continue
		}
		invite := model.Invite{}
		if err = dynamodbattribute.UnmarshalMap(result.Items[i], &invite); err != nil {
			log.Println("error while unmarshal invite: ", err)
			return []model.Invite{}, 0, l.ErrorAPIDataError
		}
		invites[i] = invite
	}
	return invites, len(result.Items), nil
}

// ArchiveAllByEmail archived all invites by email.
func (is *InviteStorage) ArchiveAllByEmail(email string) error {
	scanInput := &dynamodb.ScanInput{
		TableName:        aws.String(invitesTableName),
		FilterExpression: aws.String("#email = :email"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":email": {S: aws.String(email)},
		},
		ExpressionAttributeNames: map[string]*string{
			"#email": aws.String("email"),
		},
	}

	result, err := is.db.C.Scan(scanInput)
	if err != nil {
		log.Println("Error querying for invites:", err)
		return l.ErrorAPIDataError
	}

	for i := 0; i < len(result.Items); i++ {
		invite := model.Invite{}
		if err = dynamodbattribute.UnmarshalMap(result.Items[i], &invite); err != nil {
			log.Println("error while unmarshal invite: ", err)
		}
		if err := is.ArchiveByID(invite.ID); err != nil {
			log.Printf("error while ArchiveByID: %v", err)
		}
	}
	return nil
}

// ArchiveByID archived specific invite by its ID.
func (is *InviteStorage) ArchiveByID(id string) error {
	if _, err := xid.FromString(id); err != nil {
		return l.ErrorAPIDataError
	}
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v": {
				BOOL: aws.Bool(true),
			},
		},
		TableName: aws.String(invitesTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set archived = :v"),
	}

	if _, err := is.db.C.UpdateItem(input); err != nil {
		return l.ErrorAPIDataError
	}
	return nil
}

// Close does nothing here.
func (is *InviteStorage) Close() {}
