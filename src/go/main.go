package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-lambda-go/lambda"

	// ##########################################################################################################
	//						ERROR message on AWS Lambda -> Function Name -> Test section
	// ##########################################################################################################
	//
	// INIT_REPORT Init Duration: 1.57 ms	Phase: init	Status: error	Error Type: Runtime.InvalidEntrypoint
	// INIT_REPORT Init Duration: 1.51 ms	Phase: invoke	Status: error	Error Type: Runtime.InvalidEntrypoint
	// START RequestId: 9b402d39-e58f-4719-ae89-711d4da3740a Version: $LATEST
	// RequestId: 9b402d39-e58f-4719-ae89-711d4da3740a Error: fork/exec /var/task/bootstrap: exec format error
	// Runtime.InvalidEntrypoint
	// END RequestId: 9b402d39-e58f-4719-ae89-711d4da3740a
	// REPORT RequestId: 9b402d39-e58f-4719-ae89-711d4da3740a	Duration: 15.03 ms	Billed Duration: 16 ms	Memory Size: 128 MB	Max Memory Used: 3 MB
	//
	//
	// https://docs.aws.amazon.com/lambda/latest/dg/troubleshooting-deployment.html
	//
	// NOTE: underscope, _, invokes packages init function only
	// ##########################################################################################################
	_ "github.com/aws/aws-lambda-go/lambdacontext"

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

var logger *zap.Logger

type LocalCredentials aws.AnonymousCredentials

func (local LocalCredentials) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Source:          os.Getenv("AWS_REGION"),
		CanExpire:       false,	
		// ##########################################################################################################
		//						ERROR message on AWS Lambda -> Function Name -> Test section
		// ##########################################################################################################
		// {
		// 		"errorMessage": "operation error DynamoDB: Scan, https response error StatusCode: 400, 
		//		RequestID: 5FKRVNKQPAQ4MHOS32U7NO8ILNVV4KQNSO5AEMVJF66Q9ASUAAJG, api error 
		//		UnrecognizedClientException: The security token included in the request is invalid.",
		// 		"errorType": "OperationError"
		// }
		//
		// https://repost.aws/questions/QUJ_tReBmXQzOFot6PMuY5AA/an-error-occurred-unrecognizedclientexception-when-calling-the-listclusters-operation-the-security-token-included-in-the-request-is-invalid
		//
		// Token Validity: Make sure the session token hasn't expired; the default duration is 1 hour, but it can be extended up to 12 hours.
		// ##########################################################################################################
		Expires:         time.Now().Add(time.Hour * 1),
		// https://github.com/aws/aws-sdk-go-v2/blob/main/config/env_config.go
		SessionToken: os.Getenv("AWS_SESSION_TOKEN"),
	}, nil // error
}

var handlers = make(map[string]func(context.Context, *Request) ([]byte, error), 5)

func getUser(ctx context.Context, in *Request) ([]byte, error) {
	return nil, errors.New("getusers not available")
}

func getUsers(ctx context.Context, in *Request) ([]byte, error) {
	var table_name string = os.Getenv("FDS_APPS_USERS_TABLE")
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

	// handlers["handler_key0"] = getUser
	// handlers["handler_key1"] = getUsers

	lambda.Start(lambda_handler)
}
