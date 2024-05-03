package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
	
)

type Event interface{}

var logger *zap.Logger

func lambda_handler(c context.Context, e *Event) ([]byte, error) {
	logger.Info(fmt.Sprintln("Context:", c, "Event: ", &e))
	logger.Info("Start: Sample data lambda hander")

	cfg := aws.NewConfig()
	cfg.Region = "us-east-1"

	
	client := dynamodb.NewFromConfig(*cfg, func(options *dynamodb.Options){
		
	})
	
	params := dynamodb.ScanInput{
		TableName: aws.String("fds.apps.users"),
	}

	if o, err := client.Scan(context.Background(), &params); err != nil {
		logger.Error(fmt.Sprintln(err))
		return nil, err
	} else if b, err := json.Marshal(o.Items); err != nil {
		logger.Error(fmt.Sprintln(err))
		return nil, err
	} else {
		return b, nil
	}
}

func main() {
	logger, _ = zap.NewDevelopment()

	lambda.Start(lambda_handler)
}
