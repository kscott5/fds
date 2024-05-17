package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kscott5/fds/internal/client"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/aws/aws-lambda-go/events"
	_ "github.com/aws/aws-lambda-go/lambdacontext" // IMPORTANT: package level init() in use.

	"go.uber.org/zap"
)

func GetOrder(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewDevelopment()
	logger.Info("lambda function: dynamodb get order")
	logger.Debug(fmt.Sprint(request.PathParameters))

	tableName := os.Getenv("FDS_APPS_ORDERS_TABLE")
	if tableName == "" {
		tableName = "FDSAppsUsers"
	}

	orderid := request.PathParameters["id"]
	requires := map[string]string{"id": "string"}
	if orderid == "" {
		return nil, fmt.Errorf("requires: %s", requires)
	}

	attr, _ := attributevalue.Marshal(orderid)
	key := map[string]types.AttributeValue{
		"orderid": attr,
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

func ListOrders(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewDevelopment()
	logger.Info("lambda function: dynamodb list orders")
	logger.Debug(fmt.Sprintf("%v", request.Body))

	tableName := os.Getenv("FDS_APPS_ORDERS_TABLE")
	if tableName == "" {
		tableName = "FDSAppsOrders"
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
	} else if body, err := json.Marshal(out); err != nil{
		return nil, err
	} else {
		response := events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    client.HttpResponseHeaders,
			Body:       string(body),
		}

		return &response, nil
	}
}
