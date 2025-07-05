#* Stage 1: Build the Go application
#* Use the official Go image to build our Go application.
#* This stage creates the executable binary.
FROM public.ecr.aws/docker/library/golang:1.23-alpine AS build

# Set the current working directory inside the container
WORKDIR /app

#* Copy go.mod and go.sum files and download dependencies
#* This helps with caching layers: dependencies are downloaded only if go.mod/go.sum change
COPY go.mod go.sum ./
RUN go mod download

#* Copy the rest of the application source code
COPY . .

#* Build the Go application for Linux (because Lambda runs on Linux)
#* The output binary will be named 'main'
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

#* Stage 2: Create the final, smaller image for Lambda
#* Use a minimal base image (alpine) for a smaller deployment size.
#* We only need the compiled binary from the previous stage.
FROM public.ecr.aws/lambda/go:1
#* public.ecr.aws/lambda/go:1 is specifically designed for Go Lambdas and provides a runtime environment.

#* Copy the executable binary from the 'build' stage into the final image
#* We copy 'main' (the compiled Go binary) to /var/task/
#* Lambda expects the executable to be at /var/task/
COPY --from=build /app/main /var/task/

#* Define the command that Lambda will execute when the container starts
#* This should be the name of our executable
CMD [ "main" ]