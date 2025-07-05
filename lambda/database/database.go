package database

import (
	"fmt"
	"lambda-func/types"

	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	TABLE_NAME = "userTable"
)

type UserStore interface {
	DoesUserExist(ctx context.Context, username string) (bool, error)
	GetUser(ctx context.Context, username string) (types.User, error)
	InsertUser(ctx context.Context, user types.User) error
}

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

func (u DynamoDBClient) DoesUserExist(ctx context.Context, username string) (bool, error) {
	result, err := u.dataBaseStore.GetItemWithContext(ctx, &dynamodb.GetItemInput{ //* using & by AWS design, pass by reference for performance
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{ //* basically what we're looking for, key-value pair, primary key
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

func (u DynamoDBClient) InsertUser(ctx context.Context, user types.User) error {
	item := &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
			"password": {
				S: aws.String(user.PasswordHash),
			},
		},
	}

	_, err := u.dataBaseStore.PutItemWithContext(ctx, item)
	if err != nil {
		return err
	}
	return nil
}

func (u DynamoDBClient) GetUser(ctx context.Context, username string) (types.User, error) {
	var user types.User

	result, err := u.dataBaseStore.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return user, err
	}

	//* user does not exist
	if result.Item == nil {
		return user, fmt.Errorf("user not found")
	}

	//* unmarshal the item
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}
