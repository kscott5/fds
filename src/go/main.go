package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
)

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

type Response struct {
	Data []map[string]string `json:"data"`
}

var logger *zap.Logger

type LocalCredentials aws.AnonymousCredentials

func (local LocalCredentials) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Source:          os.Getenv("AWS_REGION"),
		CanExpire:       false,
		Expires: time.Now().Add(time.Hour * 1),
		SessionToken: os.Getenv("AWS_SESSION_TOKEN"),
	}, nil // error
}

var handlers = make(map[string]func(context.Context, *Request) (*Response, error), 5)

func getUser(ctx context.Context, in *Request) (*Response, error) {
	return nil, errors.New("getusers not available")
}

func getUsers(ctx context.Context, in *Request) (*Response, error) {
	var table_name string = os.Getenv("FDS_APPS_USERS_TABLE")
	cfg := aws.NewConfig()
	client := dynamodb.NewFromConfig(*cfg, func(options *dynamodb.Options) {
		options.Region = os.Getenv("AWS_REGION")
		options.Credentials = aws.NewCredentialsCache(LocalCredentials{})
	})

	params := dynamodb.ScanInput{
		TableName: aws.String(table_name),
		//AttributesToGet: []string {"_id", "username", "fullname"},
		//FilterExpression: aws.String("exists(fullname)"),
	}

	if output, err := client.Scan(context.Background(), &params); err != nil {
		logger.Error(fmt.Sprintln(err))
		return nil, err
	} else {
		logger.Info("Scan complete")
		data := make([]map[string]string, 1)

		for _, item := range output.Items {
			value := make(map[string]string)
			for k, v := range item {
				value[fmt.Sprintf("%s",k)] = fmt.Sprintf("%s",v)
			}

			data = append(data, value)
		}

		response := Response{
			Data: data,
		}

		fmt.Println(response)
		return &response, nil
	}
}

// curl -s -X POST http://localhost:2026/2015-03-31/functions/function/invocations -d '{"mapper": {"hello": "world", "event": "key", "list": [0,1,2,3,4]} }' | jq
// Where "TIn" and "TOut" are types compatible with the "encoding/json" standard library.
// See https://golang.org/pkg/encoding/json/#Unmarshal for how deserialization behaves
func lambda_handler(ctx context.Context, in *Request) (*Response, error) {

	switch in.Mapper[""]
	return getUsers(ctx, in)
}

func main() {
	logger, _ = zap.NewDevelopment()

	// handlers["handler_key0"] = getUser
	// handlers["handler_key1"] = getUsers

	lambda.Start(lambda_handler)
}
