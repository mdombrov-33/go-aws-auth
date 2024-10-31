package database

import (
	"lambda-func/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	TABLE_NAME = "userTable"
)

type DynamoDBClient struct {
	dataBaseStore *dynamodb.DynamoDB
}

func NewDynamoDBClient() DynamoDBClient {
	dbSession := session.Must(session.NewSession()) // creates a new AWS session and panics if there is an error
	db := dynamodb.New(dbSession)                   // creates a new DynamoDB client using the session

	return DynamoDBClient{
		dataBaseStore: db,
	}
}

//* Does this user exists?
//* How do i insert a new record into DynamoDB?

func (u DynamoDBClient) DoesUserExist(username string) (bool, error) {
	result, err := u.dataBaseStore.GetItem(&dynamodb.GetItemInput{ // using & by AWS design. Pass by reference is more performant
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{ // basically the records we are looking for. We are looking for primary key - username
			"username": {
				S: aws.String(username),
			},
		},
	})
	if err != nil {
		return true, err
	}

	if result.Item == nil {
		return false, nil
	}

	return true, nil
}

func (u DynamoDBClient) InsertUser(user types.RegisterUser) error {
	// assemble the item
	item := &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
			"password": {
				S: aws.String(user.Password),
			},
		},
	}
	// insert the item
	_, err := u.dataBaseStore.PutItem(item)
	if err != nil {
		return err
	}

	return nil
}
