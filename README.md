# Lambda Auth API

A serverless user authentication service built with AWS Lambda and Go, providing secure user registration, JWT-based login, and protected API routes.

## 🚀 Features

- **User Registration** - Secure account creation with username and password
- **JWT Authentication** - Token-based login system with secure access tokens
- **Protected Routes** - Middleware-based route protection with JWT validation
- **Persistent Storage** - DynamoDB integration for reliable user data storage
- **Distributed Tracing** - AWS X-Ray monitoring for performance insights
- **Infrastructure as Code** - Complete AWS CDK deployment setup

## 🏗️ Architecture

![Architecture](https://img.shields.io/badge/Architecture-Serverless-orange)
![Go](https://img.shields.io/badge/Go-1.20+-blue)
![AWS](https://img.shields.io/badge/AWS-Lambda-yellow)
![DynamoDB](https://img.shields.io/badge/Database-DynamoDB-blue)
![JWT](https://img.shields.io/badge/Auth-JWT-green)
![API Gateway](https://img.shields.io/badge/API-Gateway-purple)
![X-Ray](https://img.shields.io/badge/Tracing-X--Ray-red)
![CDK](https://img.shields.io/badge/IaC-CDK-orange)
![Docker](https://img.shields.io/badge/Container-Docker-blue)
![Security](https://img.shields.io/badge/Security-Bcrypt-red)

### Tech Stack

- **Go** - Core language for Lambda functions and API logic
- **AWS Lambda** - Serverless compute platform
- **API Gateway** - HTTP request routing to Lambda handlers
- **DynamoDB** - NoSQL database for user data storage
- **AWS X-Ray** - Distributed tracing and monitoring
- **AWS CDK** - Infrastructure as code deployment

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

- [AWS CLI](https://aws.amazon.com/cli/) configured with proper credentials
- [AWS CDK](https://docs.aws.amazon.com/cdk/v2/guide/getting_started.html) (`npm install -g aws-cdk`)
- [Go 1.20+](https://golang.org/dl/)
- AWS account with permissions for Lambda, DynamoDB, and API Gateway

## 🚀 Quick Start

### 1. Bootstrap AWS Environment

```bash
cdk bootstrap aws://YOUR_ACCOUNT/YOUR_REGION
```

### 2. Deploy the Stack

```bash
cdk deploy
```

This creates:
- DynamoDB table (`userTable`)
- Lambda function with Go runtime
- API Gateway with endpoints: `/register`, `/login`, `/protected`
- X-Ray tracing configuration

## 📖 API Documentation

### Register User

**POST** `/register`

```bash
curl -X POST https://YOUR_API_GATEWAY_URL/prod/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your_username",
    "password": "your_password"
  }'
```

**Responses:**
- `200 OK` - User registered successfully
- `409 Conflict` - User already exists
- `400 Bad Request` - Validation error or malformed request

### Login

**POST** `/login`

```bash
curl -X POST https://YOUR_API_GATEWAY_URL/prod/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your_username",
    "password": "your_password"
  }'
```

**Responses:**
- `200 OK` - Returns JSON with JWT access token
- `400 Bad Request` - Invalid credentials or malformed request

### Protected Route

**GET** `/protected`

```bash
curl -X GET https://YOUR_API_GATEWAY_URL/prod/protected \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Responses:**
- `200 OK` - Access granted
- `401 Unauthorized` - Missing or invalid token

## 🐳 Docker Support

### Build and Run Locally

```bash
# Build the Docker image
docker build -t lambda-auth .

# Run the container locally
docker run --rm lambda-auth
```

### Deploy with ECR (Optional)

For container-based Lambda deployment:

```bash
# Create ECR repository
aws ecr create-repository --repository-name lambda-auth

# Authenticate Docker to ECR
aws ecr get-login-password --region YOUR_REGION | \
  docker login --username AWS --password-stdin \
  YOUR_AWS_ACCOUNT_ID.dkr.ecr.YOUR_REGION.amazonaws.com

# Tag and push image
docker tag lambda-auth:latest \
  YOUR_AWS_ACCOUNT_ID.dkr.ecr.YOUR_REGION.amazonaws.com/lambda-auth:latest

docker push \
  YOUR_AWS_ACCOUNT_ID.dkr.ecr.YOUR_REGION.amazonaws.com/lambda-auth:latest
```

## 📊 Monitoring & Observability

- **AWS X-Ray** - Distributed tracing for Lambda executions and DynamoDB calls
- **CloudWatch Logs** - Centralized logging for debugging and monitoring
- **Performance Metrics** - Latency and error analysis through X-Ray console

## 🔧 Local Development

You can test Lambda handlers locally using:

- **AWS SAM CLI** - For local Lambda simulation
- **Integration Tests** - Direct API Gateway endpoint testing

## 🛠️ Troubleshooting

### Common Issues

- **Deployment Fails**: Verify AWS credentials and permissions
- **DynamoDB Access**: Ensure table exists and Lambda has proper IAM roles
- **JWT Validation**: Check token generation and verification logic
- **Runtime Errors**: Review CloudWatch logs for detailed error messages

### Debug Steps

1. Check AWS CLI configuration: `aws sts get-caller-identity`
2. Verify CDK deployment: `cdk ls`
3. Review Lambda logs in CloudWatch
4. Test API endpoints with proper headers and payloads

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
