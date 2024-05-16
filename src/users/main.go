package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	"github.com/kscott5/fds/internal/client"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/aws/aws-lambda-go/lambdacontext" // IMPORTANT: package level init() in use.

	"go.uber.org/zap"
)

var tableName string = os.Getenv("FDS_APPS_USERS_TABLE")

func putUser(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewDevelopment()
	logger.Info("lambda function: dynamodb put user")

	if tableName == "" {
		tableName = "FDSAppsUsers"
	}

	user := struct {
		UserName string `json:"username"`
		Fullname string `json:"fullname"`
	}{}

	json.Unmarshal([]byte(request.Body), &user)
	requires := map[string]string{"username": "string", "fullname": "string"}
	if user.UserName == "" || user.Fullname == "" {
		return nil, fmt.Errorf("requires: %s", requires)
	}

	_id := uuid.New().String()
	attrId, _ := attributevalue.Marshal(_id)
	if input, err := attributevalue.MarshalMap(user); err != nil {
		return nil, err
	} else {
		input["_id"] = attrId
		ddb := client.NewDynamodb(tableName)
		params := dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      input,
		}

		if _, err := ddb.PutItem(ctx, &params); err != nil {
			return nil, err
		} else {
			response := events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers:    client.HttpResponseHeaders,
				Body:       fmt.Sprintf("{\"_id\": \"%s\"}", _id),
			}

			return &response, nil
		}
	}
}

func getUser(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewDevelopment()

	logger.Info("lambda function: dynamodb get item user")
	logger.Debug(fmt.Sprint(request.PathParameters))

	if tableName == "" {
		tableName = "FDSAppsUsers"
	}

	id := ""
	json.Unmarshal([]byte(request.PathParameters["_id"]), &id)
	requires := map[string]string{"_id": "string"}
	if id == "" {
		return nil, fmt.Errorf("requires: %s", requires)
	}

	attr, _ := attributevalue.Marshal(id)
	key := map[string]types.AttributeValue{
		"_id": attr,
	}

	ddb := client.NewDynamodb(tableName)
	params := dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	}

	out := map[string]interface{}{}
	if output, err := ddb.GetItem(ctx, &params); err != nil {
		return nil, err
	} else if err := attributevalue.UnmarshalMap(output.Item, &out); err != nil {
		return nil, err
	} else if body, err := json.Marshal(out); err != nil {
		return nil, err
	} else {
		response := events.APIGatewayProxyResponse{
			StatusCode:      200,
			Headers:         client.HttpResponseHeaders,
			Body:            string(body),
			IsBase64Encoded: true,
		}

		return &response, nil
	}
}

func getUsers(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewDevelopment()
	logger.Info("lambda function: dynamodb scan get users")
	logger.Warn("filter expression or parameters not in use with this request")

	if tableName == "" {
		tableName = "FDSAppsUser"
	}

	ddb := client.NewDynamodb(tableName)
	params := dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	var out interface{}
	if output, err := ddb.Scan(ctx, &params); err != nil {
		return nil, err
	} else if err := attributevalue.UnmarshalListOfMaps(output.Items, &out); err != nil {
		return nil, err
	} else if body, err := json.Marshal(out); err != nil {
		return nil, err
	} else {
		response := events.APIGatewayProxyResponse{
			StatusCode:      200,
			Headers:         client.HttpResponseHeaders,
			Body:            string(body),
			IsBase64Encoded: true,
		}

		return &response, nil

	}
}

// curl -s -X POST http://localhost:2026/2015-03-31/functions/function/invocations -d '{"parameters": {"hello": "world", "event": "key", "list": [0,1,2,3,4]} }' | jq
func main() {
	// AWS SDK lambda function handler
	lambdaHandler := lambda.NewHandler(func(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		logger, _ := zap.NewDevelopment()
		logger.Info("FDS lambda.Start")

		switch key, _ := client.GetRequestKeyFrom(request.HTTPMethod, request.Resource); key {
		case "GET /users":
			return getUsers(ctx, request)
		case "GET /users/{_id}":
			return getUser(ctx, request)
		case "PUT /users", "PUT /user":
			return putUser(ctx, request)
		default:
			return nil, fmt.Errorf("(%s) not valid. valid request requires httpmethod and resource", key)
		}
	})

	lambda.Start(lambdaHandler)
}
