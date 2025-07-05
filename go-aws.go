package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GoAwsStackProps struct {
	awscdk.StackProps
}

func NewGoAwsStack(scope constructs.Construct, id string, props *GoAwsStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	//* Create DB table
	table := awsdynamodb.NewTable(stack, jsii.String("myUserTable"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("username"), // should be the same as primary key
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName:     jsii.String("userTable"),     // should be the same as TABLE_NAME
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY, // remove the DB when the stack is destroyed(cdk destroy)
	})

	myFunction := awslambda.NewFunction(stack, jsii.String("myLambdaFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_FROM_IMAGE(),                                   //* runtime for the lambda function (using docker image)
		Code:    awslambda.EcrImageCode_FromAsset(jsii.String(("./lambda")), nil), //* point to Dockerfile directory
		Handler: jsii.String("main"),                                              //* for images, this is usually the entrypoint in the Dockerfile
		//* Architecture for the lambda function (using ARM64 for cost efficiency), we specify this because in Dockerfile we have GOOS=linux go build -o main line
		//* By adding this, we ensure that the lambda function is built for Linux ARM64 architecture
		Architecture: awslambda.Architecture_ARM_64(),
	})

	//* Grant the lambda function read/write access to the table
	//* Doing this to connect the lambda function to the table
	table.GrantReadWriteData(myFunction)

	//* Create an API Gateway
	api := awsapigateway.NewRestApi(stack, jsii.String("myAPIGateway"), &awsapigateway.RestApiProps{
		//* Enable CORS
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowHeaders: jsii.Strings("Content-Type", "Authorization"),           // headers
			AllowMethods: jsii.Strings("GET", "POST", "DELETE", "PUT", "OPTIONS"), // methods
			AllowOrigins: jsii.Strings("*"),                                       // origins
		},
		DeployOptions: &awsapigateway.StageOptions{
			LoggingLevel: awsapigateway.MethodLoggingLevel_INFO,
		},
		CloudWatchRole: jsii.Bool(true), // fix cloudwatch error on deployment
	})

	integration := awsapigateway.NewLambdaIntegration(myFunction, nil)

	//* Define the routes
	//* Register route
	registerResource := api.Root().AddResource(jsii.String("register"), nil)
	registerResource.AddMethod(jsii.String("POST"), integration, nil)

	//* Login route
	loginResource := api.Root().AddResource(jsii.String("login"), nil)
	loginResource.AddMethod(jsii.String("POST"), integration, nil)

	//* Protected route
	protectedResource := api.Root().AddResource(jsii.String("protected"), nil)
	protectedResource.AddMethod(jsii.String("GET"), integration, nil)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewGoAwsStack(app, "AuthStack", &GoAwsStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil

}
