package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/aws/aws-lambda-go/lambdacontext" // IMPORTANT: package level init() in use.

	"go.uber.org/zap"
)

// AWS API Gateway Passthrough template requestParameters.  Reference OpenAPI
// This struct accepts any JSON parameters{} map object, key-value pair.
//
// Example:
//
// curl -s -X POST http://localhost:2026/2015-03-31/functions/function/invocations -d \
//	'{
//			"parameters": { \
//		 		"hello": "world", \
//		 		"event": "key", \
//		 		"list": [0,1,2,3,4] \
//			} \
//	}' | jq
type Request struct {
	HttpMethod string `json:"httpMethod"`
	Resource string `json:"resrouce"`
	Parameters map[string]interface{} `json:"parameters"`
}

// AWS API Gateway response template
// This struct returns a JSON data array or map object.
//
// Example:
//
// {
// 		"data":  
// 		[ 
// 			{
// 				"FullName": "Paulo Santos1",
// 				"Userid": "pasantos1",
// 				"_id": "589944140a20444fb3c85aa386acd9c4"
// 			},
// 			{
// 				"_id": "f6b3fb73-4fbb-40c0-9b4b-fa4c03c953ab",
// 				"age": 23,
// 				"disabilityTypes": [
// 					"independent living",
// 					"hearing",
// 					"vision",
// 					"mobility",
// 					"self-care"
// 				],
// 				"educationLevel": "Some College",
// 				"employmentStatus": "1099",
// 				"gender": "transwoman",
// 				"hasDisabilities": false,
// 				"healthTypes": [
// 					"Binge drinker",
// 					"Sleeplessness",
// 					"Smoker",
// 					"Obesity",
// 					"Sicklecell"
// 				],
// 				"martialStatus": "HomemakerMarried",
// 				"source": {
// 					"description": "Disability and Health Data System",
// 					"type": "Internal Marketing Research",
// 					"version": "1.0"
// 				},
// 				"userid": "f6b3fb73-4fbb-40c0-9b4b-fa4c03c953ab",
// 				"username": "f6b3fb73-4fbb-40c0-9b4b-fa4c03c953ab"
// 			}
// 		]
//  }
type Response struct {
	StatusCode int `json:"statusCode"`
	Headers map[string]string `json:"headers"`
	Data interface{} `json:"data"`
}

// Local Credentials implements CredentialsProvider.Retrieve method
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

var logger *zap.Logger

// curl -s -X POST http://localhost:2026/2015-03-31/functions/function/invocations -d '{"parameters": {"hello": "world", "event": "key", "list": [0,1,2,3,4]} }' | jq
func main() {
	logger, _ = zap.NewDevelopment()

	// AWS SDK lambda function handler
	lambda.Start(func (ctx context.Context, request *Request)(*Response, error) {
		logger.Info("lambda function: dynamodb scan users")
		logger.Debug(fmt.Sprintf("\t%s", request.Parameters))

		var table_name string = os.Getenv("FDS_APPS_USERS_TABLE")

		cfg := aws.NewConfig()
		client := dynamodb.NewFromConfig(*cfg, func(options *dynamodb.Options) {
			options.Region = os.Getenv("AWS_REGION")
			options.Credentials = aws.NewCredentialsCache(LocalCredentials{})
		})
	
		params := dynamodb.ScanInput{
			TableName: aws.String(table_name),
		}
	
		headers := make(map[string]string, 1)
		headers["content-type"] = "application/json"

		var out interface{}
		if output, err := client.Scan(context.Background(), &params); err != nil {
			return nil, err
		} else if err := attributevalue.UnmarshalListOfMaps(output.Items, &out); err != nil {
			return nil, err
		} else {
			response := Response{
				StatusCode: 200,
				Headers: headers,
				Data: out,
			}
			return &response, nil
		}
	})
}
