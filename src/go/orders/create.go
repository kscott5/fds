package orders

import (
	"context"
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

func Create(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logger,_ := zap.NewDevelopment()
	logger.Info("lambda function: dynamodb create new order")

	if tableName == "" {
		tableName = "FDSAppsOrders"
	}

	// cmapper := request.RequestContext.Authorizer["claims"]
	// claims := cmapper.(map[string]string)
	
	params, _ := client.ParseJSONRequestBody(request.Body)
	requires := map[string]string{"restaurantid": "string", "totalamount": "decimal", "items": "map"}
	if err := client.ParametersExists(*params, requires); err != nil {
		logger.Error(fmt.Sprint(err))
		
		return nil, fmt.Errorf("requires: %s", requires)
	}

	attrs := *params
	attrs["_id"] = uuid.New().String()
	//attrs["userid"] = claims["sub"]
	attrs["creation"] = time.Now().String()

	if input, err := attributevalue.MarshalMap(attrs); err != nil {
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
				Headers: client.HttpResponseHeaders,
				Body: fmt.Sprintf("{\"orderid\": \"%s\"}", attrs["_id"]),
			}

			return &response, nil
		}
	}

}