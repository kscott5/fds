package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var HttpResponseHeaders map[string]string = map[string]string{
	"content-type":               "application/json",
	"access-control-allow-orgin": "*",
}

func GetRequestKeyFrom(httpMethod, resource string) (string, error) {
	if httpMethod != "" && resource != "" {
		return fmt.Sprintf("%s %s", httpMethod, resource), nil
	}

	return "", fmt.Errorf("requires http method and resouce path")
}

// Local Credentials implements CredentialsProvider.Retrieve method
type LocalCredentials aws.AnonymousCredentials

func (local LocalCredentials) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Source:          os.Getenv("AWS_REGION"),
		CanExpire:       false,
		Expires:         time.Now().Add(time.Hour * 1),
		SessionToken:    os.Getenv("AWS_SESSION_TOKEN"),
	}, nil // error
}

func NewDynamodb(tableName string) *dynamodb.Client {
	cfg := aws.NewConfig()
	return dynamodb.NewFromConfig(*cfg, func(options *dynamodb.Options) {
		options.Region = os.Getenv("AWS_REGION")
		options.Credentials = aws.NewCredentialsCache(LocalCredentials{})
	})
}

func ParseJSONRequestBody(data string) (*map[string]string, error) {
	mapper := make(map[string]string)
	err := json.Unmarshal([]byte(data), &mapper)

	return &mapper, err
}

func ParametersExists(parameters map[string]string, requires map[string]string) error {
	for k := range requires {
		if found := parameters[k]; found == "" {
			return fmt.Errorf("requires request parameters: %s", requires)
		}
	}
	return nil
}
