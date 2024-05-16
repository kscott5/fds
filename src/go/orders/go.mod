module github.com/kscott5/fds/orders

go 1.22.1

require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/aws/aws-sdk-go-v2 v1.26.2
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.13.16
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.32.2
	github.com/google/uuid v1.6.0
	github.com/kscott5/fds/internal/client v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.27.0
)

require (
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.20.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.7 // indirect
	github.com/aws/smithy-go v1.20.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
)

replace github.com/kscott5/fds/internal/client => ../internal/
