package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	"github.com/kscott5/fds/internal/client"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	_ "github.com/aws/aws-lambda-go/lambdacontext" // IMPORTANT: package level init() in use.

	"go.uber.org/zap"
)

var tableName string = os.Getenv("FDS_APPS_ORDERS_TABLE")

func CreateOrder(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewDevelopment()
	logger.Info("lambda function: dynamodb create new order")
	logger.Debug(fmt.Sprintf("%v", request.Body))

	if tableName == "" {
		tableName = "FDSAppsOrders"
	}

	// anonymous structure
	order := struct { 
		RestaurantId string `json:"restaurantid"`
		TotalAmount float64 `json:"totalamount"`
		Items []string `json:"items"`
		OrderId string
		UserId string
		Status string
		PlacedOn string 
	}{}

	// extract and validate request body
	if err :=json.Unmarshal([]byte(request.Body), &order); err != nil {
		return nil, err
	}
	requires := map[string]string{"restaurantid": "string", "totalamount": "decimal", "items": "map"}
	if order.RestaurantId == "" || order.TotalAmount == 0 || len(order.Items) == 0 {
		return nil, fmt.Errorf("requires: %s", requires)
	}

	// Cognitio user pool authentication and authorization
	// cmapper := request.RequestContext.Authorizer["claims"]
	// claims := cmapper.(map[string]string)
	// order.UserId = claims["sub"]

	order.OrderId = uuid.New().String()
	order.Status = "PLACED"
	order.PlacedOn = time.Now().String()

	if input, err := attributevalue.MarshalMap(order); err != nil {
		return nil, err
	} else {
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
				Body:       fmt.Sprintf("{\"orderid\": \"%s\"}", order.OrderId),
			}

			return &response, nil
		}
	}
}
