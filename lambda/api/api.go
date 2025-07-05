package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	//* when we call api functions, we want to interact with the database
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("RegisterUserHandler called")

	var registerUser types.RegisterUser

	err := json.Unmarshal([]byte(request.Body), &registerUser)
	if err != nil {
		fmt.Println("Error unmarshalling register request body:", err)
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	fmt.Println("Register user payload:", registerUser)

	if registerUser.Username == "" || registerUser.Password == "" {
		fmt.Println(("Validation failed - fields empty"))
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request - fields empty",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	//* does a user with this username already exist?
	userExists, err := api.dbStore.DoesUserExist(registerUser.Username)
	if err != nil {
		fmt.Println("Error checking if user exists:", err)
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if userExists {
		fmt.Println("User already exists:", registerUser.Username)
		return events.APIGatewayProxyResponse{
			Body:       "User already exists",
			StatusCode: http.StatusConflict,
		}, nil
	}

	user, err := types.NewUser(registerUser)
	if err != nil {
		fmt.Println("Error creating user struct:", err)
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error creating user: %w", err)
	}

	//* we know that a user does not exists
	//* insert the user into the database
	err = api.dbStore.InsertUser(user)
	if err != nil {
		fmt.Println("Error inserting user into database:", err)
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error inserting user: %w", err)
	}

	fmt.Println("User registered successfully:", user.Username)
	return events.APIGatewayProxyResponse{
		Body:       "User registered",
		StatusCode: http.StatusOK,
	}, nil
}

func (api ApiHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("LoginUserHandler called")
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginRequest)
	if err != nil {
		fmt.Println("Error unmarshalling login request body:", err)
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	fmt.Println("Login attempt for user:", loginRequest.Username)

	user, err := api.dbStore.GetUser(loginRequest.Username)
	if err != nil {
		fmt.Println("Error retrieving user from database:", err)
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if !types.ValidatePassword(user.PasswordHash, loginRequest.Password) {
		fmt.Println("Invalid password attempt for user:", loginRequest.Username)
		return events.APIGatewayProxyResponse{
			Body:       "Invalid credentials",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	accessToken := types.CreateToken(user)

	if accessToken == "" {
		fmt.Println("Error creating access token for user:", user.Username)
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	//* Prepare the response with the access token
	responseBody, err := json.Marshal(map[string]string{"access_token": accessToken})
	if err != nil {
		fmt.Println("Error marshalling response body:", err)
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error marshalling response: %w", err)
	}

	fmt.Println("User logged in successfully:", user.Username)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(responseBody),
	}, nil
}
