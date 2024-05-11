package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/aws/aws-lambda-go/lambdacontext" // IMPORTANT: package level init() in use.

	"go.uber.org/zap"
)


func requestKey(httpMethod, resource string) (string, error) {
	if httpMethod != "" && resource != "" {
		return fmt.Sprintf("%s %s", httpMethod, resource), nil
	}

	return "", fmt.Errorf("requires http method and resouce path")
}

// Local Credentials implements CredentialsProvider.Retrieve method
type LocalCredentials aws.AnonymousCredentials

func (local LocalCredentials) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Source:          os.Getenv("AWS_REGION"),
		CanExpire:       false,
		Expires:         time.Now().Add(time.Hour * 1),
		SessionToken:    os.Getenv("AWS_SESSION_TOKEN"),
	}, nil // error
}

var (
	logger    *zap.Logger
	tableName string
	headers map[string]string = map[string]string{
		"content-type": "application/json",
		"access-control-allow-orgin": "*",
	}
)

func newDynamodbClient() *dynamodb.Client {
	var found bool
	if tableName, found = os.LookupEnv("FDS_APPS_USERS_TABLE"); !found {
		tableName = "FDSAppsUsers"
	}

	cfg := aws.NewConfig()
	return dynamodb.NewFromConfig(*cfg, func(options *dynamodb.Options) {
		options.Region = os.Getenv("AWS_REGION")
		options.Credentials = aws.NewCredentialsCache(LocalCredentials{})
	})
}

func parametersExists(parameters map[string]string, requires map[string]string) (error) {
	logger.Debug(fmt.Sprint("parametersExists", parameters, " requires:", requires))

	for k := range requires {
		if found := parameters[k]; found == "" {
			return fmt.Errorf("requires request parameters: %s", requires)
		}
	}
	return nil
}

func putUser(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logger.Info("lambda function: dynamodb scan users")

	requires := map[string]string{"username": "string", "fullname": "string"}
	if err := parametersExists(request.PathParameters, requires); err != nil {
		logger.Error(fmt.Sprint(err))
		
		return nil, fmt.Errorf("requires: %s", requires)
	}

	attrs := request.PathParameters
	attrs["_id"] = uuid.New().String()

	if input, err := attributevalue.MarshalMap(attrs); err != nil {
		return nil, err
	} else {
		client := newDynamodbClient()
		params := dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      input,
		}

		if _, err := client.PutItem(ctx, &params); err != nil {
			return nil, err
		} else {
			response := events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: headers,
				Body: fmt.Sprintf("{\"_id\": \"%s\"}", attrs["_id"]),
			}

			return &response, nil
		}
	}

}

func getUser(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logger.Info("lambda function: dynamodb get item user")
	logger.Debug(fmt.Sprint(request.PathParameters))

	requires := map[string]string{"_id": "string"}
	if err := parametersExists(request.PathParameters, requires); err != nil {
		logger.Error(fmt.Sprint(err))
		
		return nil, fmt.Errorf("requires: %s", requires)
	}

	attrs := request.PathParameters
	key, _ := attributevalue.MarshalMap(attrs)

	client := newDynamodbClient()
	params := dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: key,
	}

	out := map[string]interface{}{}
	if output, err := client.GetItem(ctx, &params); err != nil {
		return nil, err
	} else if err := attributevalue.UnmarshalMap(output.Item, &out); err != nil {
		return nil, err
	} else if body, err :=json.Marshal(out); err != nil {		
		return nil, err
	} else {
		response := events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: headers,
			Body: string(body),
			IsBase64Encoded: true,
		}

		return &response, nil
	}
}

func getUsers(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logger.Info("lambda function: dynamodb scan get users")
	logger.Warn(fmt.Sprintf("filter expression or parameters not in use with this request: %s", request))
	
	client := newDynamodbClient()
	params := dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	var out interface{}
	if output, err := client.Scan(ctx, &params); err != nil {
		return nil, err
	} else if err := attributevalue.UnmarshalListOfMaps(output.Items, &out); err != nil {
		return nil, err
	} else if body, err := json.Marshal(out); err != nil {
		return nil, err
	} else {
		response := events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: headers,
			Body:  string(body),
			IsBase64Encoded: true,
		}

		return &response, nil

	}
}

// curl -s -X POST http://localhost:2026/2015-03-31/functions/function/invocations -d '{"parameters": {"hello": "world", "event": "key", "list": [0,1,2,3,4]} }' | jq
func main() {
	logger, _ = zap.NewDevelopment()
	logger.Info("FDS main")

	// AWS SDK lambda function handler
	lambdaHandler := lambda.NewHandler(func(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		logger.Info("FDS lambda.Start")
		
		switch id, _ := requestKey(request.HTTPMethod, request.Resource); id {
		case "GET /users":
			return getUsers(ctx, request)
		case "GET /users/{_id}":
			return getUser(ctx, request)
		case "PUT /user":
			return putUser(ctx, request)
		default:
			return nil, fmt.Errorf("(%s) not valid. valid request requires httpmethod and resource", id)
		}
	})

	lambda.Start(lambdaHandler)
}
