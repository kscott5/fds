package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
)

type Event interface{}

var logger *zap.Logger

func lambda_handler(c context.Context, e *Event) (*Event, error) {
	logger.Info(fmt.Sprintln("Context:", c, "Event: ", &e))
	
	loadData()

	return e, nil
}

func loadData() {
	logger.Info("Start: Sample data lambda hander")
	// NOTE: Local development with docker pull amazon/dynamodb-local
	logger.Info("Load test data")	
	data := make(map[int][3]string, 10)
	data[0] = [3]string{uuid.New().String(), "marivera0", "Martha Rivera"}
	data[1] = [3]string{uuid.New().String(), "nikkwolf0", "Nikki Wolf"}
	data[2] = [3]string{uuid.New().String(), "pasantos0", "Paulo Santos"}

	data[3] = [3]string{uuid.New().String(), "marivera1", "Martha Rivera"}
	data[4] = [3]string{uuid.New().String(), "nikkwolf1", "Nikki Wolf"}
	data[5] = [3]string{uuid.New().String(), "pasantos1", "Paulo Santos"}

	data[6] = [3]string{uuid.New().String(), "marivera2", "Martha Rivera"}
	data[7] = [3]string{uuid.New().String(), "nikkwolf2", "Nikki Wolf"}
	data[8] = [3]string{uuid.New().String(), "pasantos2", "Paulo Santos"}

	for _, v := range data {
		params := dynamodb.PutItemInput {
			TableName: aws.String("fds.apps.users"),
			Item: map[string]types.AttributeValue{
				"_id": &types.AttributeValueMemberS{ Value: v[0]},
				"userid": &types.AttributeValueMemberS{ Value: v[1] },
				"fullname": &types.AttributeValueMemberS{ Value: v[2]},
			},
 		} // end params

		client := dynamodb.NewFromConfig(aws.Config{})
		client.PutItem(context.Background(), &params)
	}

	logger.Info("End: Sample data lambda hander")
}

func main() {
	logger, _ = zap.NewDevelopment()
	
	lambda.Start(lambda_handler)
}