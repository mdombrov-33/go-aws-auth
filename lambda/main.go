package main

import (
	"fmt"
	"lambda-func/app"
	"net/http"

	"github.com/aws/aws-lambda-go/events" // allows us to extract paths, requests etc.
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"`
}

// * Take in a payload and do something with it
func HandleRequest(event MyEvent) (string, error) {
	if event.Username == "" {
		return "", fmt.Errorf("username is empty")
	}

	return fmt.Sprintf("Successfully called by - %s", event.Username), nil
}

func main() {
	lambdaApp := app.NewApp()
	//* Hook lambda function to the gateway
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/register":
			return lambdaApp.APIHandler.RegisterUserHandler(request)
		case "/login":
			return lambdaApp.APIHandler.LoginUser(request)
		default:
			return events.APIGatewayProxyResponse{
				Body:       "Not Found",
				StatusCode: http.StatusNotFound,
			}, nil
		}
	})

}
