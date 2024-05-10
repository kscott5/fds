package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/aws/aws-lambda-go/lambdacontext" // IMPORTANT: package level init() in use.

	"go.uber.org/zap"
)

type Request struct {
	HttpMethod string                 		`json:"httpMethod"`
	Resource   string                 		`json:"resource"`
	Parameters map[string]interface{} 		`json:"parameters,omitempty"`
	PathParameters map[string]interface{} 	`json:"pathParameters,omitempty"`
}

type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Data       interface{}       `json:"data"`
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
		"header": "application/jsom",
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

func parametersExists(parameters map[string]interface{}, requires map[string]string) (error) {
	logger.Debug(fmt.Sprint("parametersExists", parameters, " requires:", requires))

	for k := range requires {
		if found := parameters[k]; found == nil {
			return fmt.Errorf("requires request parameters: %s", requires)
		}
	}
	return nil
}

func putUser(ctx context.Context, request *Request) (*Response, error) {
	logger.Info("lambda function: dynamodb scan users")

	requires := map[string]string{"username": "string", "fullname": "string"}
	if err := parametersExists(request.Parameters, requires); err != nil {
		logger.Error(fmt.Sprint(err))
		
		return nil, fmt.Errorf("requires: %s", requires)
	}

	attrs := request.Parameters
	attrs["_id"] = uuid.New().String()

	if input, err := attributevalue.MarshalMap(attrs); err != nil {
		logger.Error(fmt.Sprint(err))
		return nil, fmt.Errorf("json format and mappers requires: %s", requires)
	} else {
		client := newDynamodbClient()
		params := dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      input,
		}

		if _, err := client.PutItem(ctx, &params); err != nil {
			logger.Error(fmt.Sprint(err))
			return nil, err
		} else {
			response := Response{
				StatusCode: 200,
				Headers: headers,
				Data: map[string]string{"_id": attrs["_id"].(string)},
			}
	
			return &response, nil
		}
	}
}

func getUser(ctx context.Context, request *Request) (*Response, error) {
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

	if output, err := client.GetItem(ctx, &params); err != nil {
		logger.Error(fmt.Sprint(err))
		return nil, fmt.Errorf("could not access this user data")
	} else {
		
		var data map[string]struct{
			_id string
			FullName string
			UserName string
		}

		// NOTE: Expectation different 
		attributevalue.UnmarshalMap(output.Item, data)

		response := Response{
			StatusCode: 200,
			Headers: headers,
			Data: data,
		}

		return &response, nil
	}
}

func getUsers(ctx context.Context, request *Request) (*Response, error) {
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
	} else {
		response := Response{
			StatusCode: 200,
			Headers: headers,
			Data: out,
		}

		return &response, nil

	}
}

// curl -s -X POST http://localhost:2026/2015-03-31/functions/function/invocations -d '{"parameters": {"hello": "world", "event": "key", "list": [0,1,2,3,4]} }' | jq
func main() {
	logger, _ = zap.NewDevelopment()
	logger.Info("FDS main")

	// AWS SDK lambda function handler
	lambda.Start(func(ctx context.Context, request *Request) (*Response, error) {
		logger.Info("FDS lambda.Start")
		logger.Debug(fmt.Sprintf("request %s", request))

		requestKey := fmt.Sprintf("%s %s", request.HttpMethod, request.Resource)

		switch requestKey {
		case "GET /users":
			return getUsers(ctx, request)
		case "GET /users/{_id}":
			return getUser(ctx, request)
		case "PUT /user":
			return putUser(ctx, request)
		default:
			return nil, fmt.Errorf("invalid request: %s", requestKey)
		}
	})
}
