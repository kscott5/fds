package authorizer

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/aws/aws-lambda-go/lambdacontext" // IMPORTANT: package level init() in use.
)

func generateAllow(principalId, resource string) {
	generatePolicy(principalId, resource, "ALLOW")
}

func generateDeny(principalId, resource string) {
	generatePolicy(principalId, resource, "DENY")
}

func generatePolicy(principalId, resource, effect string) map[string]string {
	
	return nil
}

func main() {
	lambdaHandler := lambda.NewHandler(func(ctx context.Context, request *events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {

		headers := request.Headers
		queryParams := request.QueryStringParameters
		stageVars := request.StageVariables

		// Parse the input for the parameter values
		methodArn := strings.Split(request.MethodArn, ":")
		apiGatewayArn := strings.Split(methodArn[5], "/")

		// Save the ARN parts
		awsAccountId := methodArn[4]
		region := methodArn[3]
		route := methodArn[2]
		stage := methodArn[1]
		apiId := methodArn[0]


		// Perform authorization to return the Allow policy for correct parameters
    	// and the 'Unauthorized' error, otherwise.
		if headers["HeaderAuth1"] == "headerValue1" && queryParams["QueryString1"] == "queryValue1" {
			return generateAllow("me", request.MethodArn)
		}
		return nil, fmt.Errorf("not available")
	})

	lambda.Start(lambdaHandler)
}