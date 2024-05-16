package main

import (
	"context"
	"fmt"

	"github.com/kscott5/fds/internal/client"
	"github.com/kscott5/fds/orders/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/aws/aws-lambda-go/lambdacontext" // IMPORTANT: package level init() in use.

	"go.uber.org/zap"
)

func main() {
	// AWS SDK lambda function handler
	lambdaHandler := lambda.NewHandler(func(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		logger, _ := zap.NewDevelopment()
		logger.Info("FDS lambda.Start orders")

		switch key, _ := client.GetRequestKeyFrom(request.HTTPMethod, request.Resource); key {
		case "PUT /orders", "PUT /order":
			return services.CreateOrder(ctx, request)
		default:
			return nil, fmt.Errorf("(%s) not valid. valid request requires httpmethod and resource", key)
		}
	})

	lambda.Start(lambdaHandler)
}
