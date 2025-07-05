build:
	@echo "Building lambda function"
	@GOOS=linux GOARCH=amd64 go build -o bootstrap

	# Create a zip file with the bootstrap binary, not used because we are using a Docker image with ECR
	# @zip function.zip bootstrap