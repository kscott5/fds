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

func CreateOrder(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewDevelopment()
	logger.Info("lambda function: dynamodb create new order")
	logger.Debug(fmt.Sprintf("%v", request.Body))


	tableName := os.Getenv("FDS_APPS_ORDERS_TABLE")
	if tableName == "" {
		tableName = "FDSAppsOrders"
	}

	// anonymous structure
	data := struct { 
		RestaurantId string `json:"restaurantid"`
		TotalAmount float64 `json:"totalamount"`
		Items []string `json:"items"`
		OrderId string
		UserId string
		Status string
		PlacedOn string 
	}{ OrderId: uuid.New().String(), Status: "PLACED"}

	// extract and validate request body
	if err :=json.Unmarshal([]byte(request.Body), &data); err != nil {
		return nil, err
	}

	requires := map[string]string{"restaurantid": "string", "totalamount": "decimal", "items": "map"}
	if data.RestaurantId == "" || len(data.Items) == 0  || data.TotalAmount <= 0  {
		return nil, fmt.Errorf("requires: %s", requires)
	}

	data.UserId = "PLACEHOLDER"
	// Cognitio user pool authentication and authorization
	if cmapper := request.RequestContext.Authorizer["claims"]; cmapper != nil {
		claims := cmapper.(map[string]string)
		if found := claims["sub"]; found != "" {
			data.UserId = claims["sub"]
		}
	}
	
	data.OrderId = uuid.New().String()
	data.Status = "PLACED"
	data.PlacedOn = time.Now().UTC().String()

	order := map[string]interface{}{
		"orderid": data.OrderId,
		"userid": data.UserId,
		"data": data,
	}

	if input, err := attributevalue.MarshalMap(order); err != nil {
		return nil, err
	} else {
		logger.Debug("new dynamodb client session")
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
				Body:       fmt.Sprintf("{\"orderid\": \"%s\"}", order["orderId"]),
			}

			return &response, nil
		}
	}
}
