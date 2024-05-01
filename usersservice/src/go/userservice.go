package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
)

var logger *zap.Logger

type Event interface{}

// curl -s -X POST http://localhost:8080/2015-03-31/functions/function/invocations -d '{"data": [{"hello": "world"}, {"event": "key"}, {"list": [0,1,2,3,4]} ]}' | jq
//
//	{
//	  "data": [
//	    {
//	      "hello": "world"
//	    },
//	    {
//	      "event": "key"
//	    },
//	    {
//	      "list": [
//	        0,
//	        1,
//	        2,
//	        3,
//	        4
//	      ]
//	    }
//	  ]
//	}
type Event2 struct {
	Data interface{} `json:"data"`
}

func lambda_handler(c context.Context, e *Event) (*Event, error) {
	logger.Info(fmt.Sprintln("Context:", c, "Event: ", &e))
	
	
	return e, nil
}

func getUsers() []byte {
	// NOTE: Local development with docker pull amazon/dynamodb-local

	var scanInput dynamodb.ScanInput = dynamodb.ScanInput{} 
	scanInput.TableName = aws.String("fds.apps.users")

	client := dynamodb.NewFromConfig(aws.Config{})
	if output, err := client.Scan(context.Background(), &scanInput);  err != nil {
		logger.Error(fmt.Sprintln(err))
	} else if data, err:= json.Marshal(output.Items); err != nil {
		logger.Error(fmt.Sprintln(err))
	} else {
		return data
	}
	return nil
}

func main() {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		fmt.Printf("error creating zap logger, error:%v", err)
		return
	}
	logger.Info("starting lambda_handler")
	lambda.Start(lambda_handler)
}
