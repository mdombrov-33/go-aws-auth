package app

import (
	"lambda-func/api"
	"lambda-func/database"
)

type App struct {
	APIHandler api.ApiHandler
}

func NewApp() App {
	// Here we actually initialize DB store
	// gets passed down into the api handler
	db := database.NewDynamoDBClient()
	apiHandler := api.NewApiHandler(db)

	return App{
		APIHandler: apiHandler,
	}
}
