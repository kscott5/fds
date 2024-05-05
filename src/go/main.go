package main

import (
	"context"
	"errors"
	"fmt"
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
)

var table_name string = os.Getenv("FDS_APPS_USERS_TABLE")

// AWS API Gateway Passthrough template format
// curl -s -X POST http://localhost:2026/2015-03-31/functions/function/invocations -d \
//
//	'{
//			"mapper": { \
//		 		"hello": "world", \
//		 		"event": "key", \
//		 		"list": [0,1,2,3,4] \
//			} \
//	}' | jq
type Request struct {
	Mapper map[string]interface{} `json:"mapper"`
}

var logger *zap.Logger

type LocalCredentials aws.AnonymousCredentials

func (local LocalCredentials) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Source:          os.Getenv("AWS_REGION"),
		CanExpire:       true,
		Expires:         time.Now().Add(time.Second * 240),
	}, nil // error
}

var handlers = make(map[string]func(context.Context, *Request) ([]byte, error), 5)

func getUser(ctx context.Context, in *Request) ([]byte, error) {
	return nil, errors.New("getusers not available")
}

func getUsers(ctx context.Context, in *Request) ([]byte, error) {
	cfg := aws.NewConfig()
	client := dynamodb.NewFromConfig(*cfg, func(options *dynamodb.Options) {
		options.Region = os.Getenv("AWS_REGION")
		options.Credentials = aws.NewCredentialsCache(LocalCredentials{})
	})

	params := dynamodb.ScanInput{
		TableName: aws.String(table_name),
	}

	if output, err := client.Scan(context.Background(), &params); err != nil {
		logger.Error(fmt.Sprintln(err))
		return nil, err
	} else if data, err := json.Marshal(output.Items); err != nil {
		logger.Error(fmt.Sprintln(err))
		return nil, err
	} else {
		logger.Info("Scan complete")
		return data, nil
	}
}

// curl -s -X POST http://localhost:2026/2015-03-31/functions/function/invocations -d '{"mapper": {"hello": "world", "event": "key", "list": [0,1,2,3,4]} }' | jq
// Where "TIn" and "TOut" are types compatible with the "encoding/json" standard library.
// See https://golang.org/pkg/encoding/json/#Unmarshal for how deserialization behaves
func lambda_handler(ctx context.Context, in *Request) ([]byte, error) {
	// Generic example
	// var s string = in.Mapper["event"].(string)
	// l := in.Mapper["list"].([]interface{})
	// sum := l[3].(float64) + l[2].(float64)
	// fmt.Println(in.Mapper["hello"], s, l[3], "+", l[2], "=", sum)

	return getUsers(ctx, in)
}

func main() {
	logger, _ = zap.NewDevelopment()

	handlers["handler_key0"] = getUser
	handlers["handler_key1"] = getUsers

	lambda.Start(lambda_handler)
}
