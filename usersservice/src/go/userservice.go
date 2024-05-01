package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
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
	aws.String("")
	return e, nil
}

func getUsers() {
	var ddb = aws.
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
