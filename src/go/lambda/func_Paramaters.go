package lambda

import(
	"context"
	"os"
	"time"
	"github.com/aws/aws-sdk-go-v2/aws"
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
