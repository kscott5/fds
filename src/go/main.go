package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/aws/aws-lambda-go/lambda"

	// s"github.com/aws/aws-sdk-go-v2/aws"
	// s"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
)

// Generic key-value mapper
type Event map[string]interface{}

var logger *zap.Logger

var handlers = make(map[string]func(context.Context, *Event) ([]byte, error), 5)

func getUser(ctx context.Context, in *Event) ([]byte, error) {
	return nil, errors.New("getusers not available")
}

func getUsers(ctx context.Context, in *Event) ([]byte, error) {
	return nil, errors.New("getusers not available")
}

// curl -s -X POST http://localhost:2026/2015-03-31/functions/function/invocations -d '{"data": [{"hello": "world"}, {"event": "key"}, {"list": [0,1,2,3,4]} ]}' | jq
// curl -s -X POST http://localhost:2026/2015-03-31/functions/function/invocations -d '{"data": {"hello": "world", "event": "key", "list": [0,1,2,3,4]} }' | jq
// curl -s -X POST http://localhost:2026/2015-03-31/functions/function/invocations -d '{"hello": "world", "event": "key", "list": [0,1,2,3,4]} ' | jq
func lambda_handler(ctx context.Context, in *Event) (*Event, error) {
	req := *in

	required := map[string]string{"hello": "string:required", "event": "string:required"}
	//optional := map[string]string{"list": "[]int"}

	mapper := map[string]interface{}{"hello":"", "event":"", "list": make([]int,0)}

	fmt.Println(mapper)

	v := reflect.ValueOf(req)
	iter := v.MapRange()
	for iter.Next() {
		if _, found := required[iter.Key().String()]; !found {
			return nil, errors.New("request body not correct")
		}
	}
	
	return in, nil //errors.New("lambda_handler not available\n")
}

func main() {
	logger, _ = zap.NewDevelopment()

	handlers["handler_key0"] = getUser
	handlers["handler_key1"] = getUsers

	lambda.Start(lambda_handler)
}
