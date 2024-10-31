package api

import (
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
)

type ApiHandler struct {
	dbStore database.DynamoDBClient // when we call api functions, we want to interact with the database
}

func NewApiHandler(dbStore database.DynamoDBClient) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(event types.RegisterUser) error {
	if event.Username == "" || event.Password == "" {
		return fmt.Errorf("request has empty parameters")
	}

	// does a user with this username already exist?
	userExists, err := api.dbStore.DoesUserExist(event.Username)
	if err != nil {
		return fmt.Errorf("there was an error checking if user exists %w", err)
	}

	if userExists {
		return fmt.Errorf("a user with this username already exists")
	}

	// we know that a user does not exists
	// insert the user into the database
	err = api.dbStore.InsertUser(event)
	if err != nil {
		return fmt.Errorf("there was an error registering the user %w", err)
	}

	return nil

}
