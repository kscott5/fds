package client

import (
	"context"
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