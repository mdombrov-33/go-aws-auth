package database

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
