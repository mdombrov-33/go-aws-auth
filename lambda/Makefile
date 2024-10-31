build:
	@echo "Building lambda function"
	@GOOS=linux GOARCH=amd64 go build -o bootstrap
	@zip function.zip bootstrap