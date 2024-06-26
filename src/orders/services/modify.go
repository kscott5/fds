package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kscott5/fds/internal/client"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	_ "github.com/aws/aws-lambda-go/lambdacontext" // IMPORTANT: package level init() in use.

	"go.uber.org/zap"
)
const (
	MaxElapseTimeMilliSecs = 600000
)

func ModifyOrder(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewDevelopment()
	logger.Info("lambda function: dynamodb modify order")
	logger.Debug(fmt.Sprintf("%v", request.Body))

	userid, _ := GetUserFromRequestContext(request.RequestContext.Authorizer)

	// extract previous order and validate request.PathParameters with PUT /orders/{orderid} or PUT /order/{orderid}
	po := Order{}
	if response, err := GetOrder(ctx, request); err != nil {
		return nil, err
	} else if err := json.Unmarshal([]byte(response.Body), &po); err != nil {
		return nil, err
	} else if po.UserId != userid || po.Status == Acknowledged || (time.Now().UnixMilli() - int64(po.PlacedOn)) > MaxElapseTimeMilliSecs {
		return nil, fmt.Errorf("order updates not acceptable. previous order was acknowledged")
	}

	logger.Info("lambda function: processing order updates")
	logger.Debug(fmt.Sprintf("lambda function: order status %s", po.Status))

	tableName := os.Getenv("FDS_APPS_ORDERS_TABLE")
	if tableName == "" {
		tableName = DefaultOrderTable
	}

	// extract and validate request body
	data := Order{}
	if err :=json.Unmarshal([]byte(request.Body), &data); err != nil {
		return nil, err
	}

	requires := map[string]string{"restaurantid": "string", "totalamount": "decimal", "items": "map"}
	if data.RestaurantId == "" || len(data.Items) == 0  || data.TotalAmount <= 0  {
		return nil, fmt.Errorf("requires: %s", requires)
	}

	data.Status = Placed
	data.ModifiedOn = UnixMilliTime(time.Now().UnixMilli())

	// current order
	co := map[string]interface{}{
		"orderid": data.OrderId,
		"userid": data.UserId,
		"data": data.Items,
	}

	if input, err := attributevalue.MarshalMap(co); err != nil {
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
				Body:       fmt.Sprintf("{\"orderid\": \"%s\", \"description\": \"updates are complete\"}", co["orderId"]),
			}

			return &response, nil
		}
	}
}
