package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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

const (
	OrderIdAttribute = "id"
	NonContextualUserId = "placeholder"
	DefaultOrderTable = "FDSAppsOrders"
)

type OrderStatus uint
const (
	Invalid OrderStatus = iota
	Placed
	Acknowledged
	Cancelled
	Paused
)
func (os OrderStatus) MarshalJSON() ([]byte, error) {
	switch os {
	case Placed:
		return []byte("placed"), nil
	case Acknowledged:
		return []byte("acknowledged"), nil
	case Cancelled:
		return []byte("cancelled"), nil
	default:
		return []byte("invalid"), fmt.Errorf("invalid order status marshal json not available")
	}
}
func (os *OrderStatus) Unmarshaler(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch strings.TrimSpace(strings.ToLower(v)) {
	case "placed":
		*os = Placed
	case "acknowledged":
		*os = Acknowledged
	case "cancelled":
		*os = Cancelled
	default:
		*os = Invalid
		return fmt.Errorf("order status unmarshaler not available")
	}
	return nil
}
func (os OrderStatus) String() string {
	switch os {
	case Placed:
		return "placed"
	case Acknowledged:
		return "acknowledged"
	case Cancelled:
		return "cancelled"
	default:
		return "invalid"
	}
}

type Items struct {
	ItemId		string `json:"itemid"`
	Description string  `json:"description"`
	Quanity     int     `json:"quanity"`
	Amount      float64 `json:"amount"`
}

type Order struct {
	RestaurantId string  		`json:"restaurantid"`
	TotalAmount  float64 		`json:"totalamount"`
	Items        []Items 		`json:"items"`
	OrderId      string	 		`json:"orderid"`
	UserId       string	 		`json:"userid"`
	Status       OrderStatus	`json:"status"`
	PlacedOn     time.Time		`json:"placedon"`
	ModifiedOn	 time.Time		`json:"modifiedon"`
}

func GetUserFromRequestContext(authorizer map[string]interface{}) (string, error) {
	userId := NonContextualUserId

	// Cognitio user pool authentication and authorization
	if cmapper := authorizer["claims"]; cmapper != nil {
		claims := cmapper.(map[string]string)
		if found := claims["sub"]; found != "" {
			userId = claims["sub"]
		}

		return userId, nil
	}

	return userId, fmt.Errorf("user id not found with request context")
}

func CreateOrder(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewDevelopment()
	logger.Info("lambda function: dynamodb create new order")
	logger.Debug(fmt.Sprintf("%v", request.Body))

	tableName := os.Getenv("FDS_APPS_ORDERS_TABLE")
	if tableName == "" {
		tableName = DefaultOrderTable
	}

	// extract and validate request body
	data := Order{}
	if err := json.Unmarshal([]byte(request.Body), &data); err != nil {
		return nil, err
	}

	requires := map[string]string{"restaurantid": "string", "totalamount": "decimal", "items": "map"}
	if data.RestaurantId == "" || len(data.Items) == 0 || data.TotalAmount <= 0 {
		return nil, fmt.Errorf("requires: %s", requires)
	}

	data.UserId, _ = GetUserFromRequestContext(request.RequestContext.Authorizer)
	data.OrderId = uuid.New().String()
	data.Status = Placed
	data.PlacedOn = time.Now()
	data.ModifiedOn = data.PlacedOn

	order := map[string]interface{}{
		"orderid": data.OrderId,
		"userid":  data.UserId,
		"data":    data,
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
