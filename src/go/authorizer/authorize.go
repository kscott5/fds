package authorizer

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/aws/aws-lambda-go/lambdacontext" // IMPORTANT: package level init() in use.
)

const (
	// The policy version used for the evaluation. This should always be '2012-10-17'
	version = "2012-10-17"
	// The regular expression used to validate resource paths for the policy
	pattern = `^[/.a-zA-Z0-9-\*]+$`
)

var (
	userPoolId = os.Getenv("USER_POOL_ID")
	appClientId = os.Getenv("APPLICATION_CLIENT_ID")
	adminGroupName = os.Getenv("ADMIN_GROUP_NAME")

	httpVerb = map[string]string{
		"GET":"GET",
		"POST": "POST",
		"PUT": "PUT",
		"PATCH": "PATCH",
		"HEAD": "HEAD",
		"DELETE": "DELETE",
		"OPTIONS": "OPTIONS",
		"ALL": "*",
	}
)

type PolicyResponse struct {
	events.APIGatewayCustomAuthorizerResponse
	
	allowMethods []map[string]string
	denyMethods  []map[string]string

	accountId string
	region string
	route string
	stage string
	apiId string
}

func (pr *PolicyResponse) addMethod(effect, verb, resource string, conditions []string) error {
	/* Adds a method to the internal lists of allowed or denied methods. Each object in
    the internal list contains a resource ARN and a condition statement. The condition
    statement can be null. */
    if found := httpVerb[verb]; found != "" && verb != "*" {
        return fmt.Errorf("Invalid HTTP verb '%s'.", verb);
    }
    if found, _ := regexp.Match(pattern, resource.([]byte)); !found {
         return fmt.Errorf("Invalid resource path: '%s'. Path should match '%s'", resource, pattern);
    }
    if resource[0:1] == "/" {
        resource = resource[1:]
    }

	// https://pkg.go.dev/strings#Builder
    resourceArnBuilder := strings.Builder{}
	fmt.Fprintf(&resourceArnBuilder, "arn:aws:execute-api:%s:%s:%s/%s/%s/%s", region, accountId, apiId, stage, verb, resource)
	
    if (strings.ToLower(effect) == "allow") {
		pr.allowMethods = append(pr.allowMethods, map[string]string{
        	"resourceArn": resourceArnBuilder.String(),
        	"conditions": conditions,
      })
    } else if (strings.ToLower(effect) == "deny") {
        pr.denyMethods = append(pr.denyMethods, map[string]string{
        	"resourceArn": resourceArnBuilder.String(),
        	"conditions": conditions,
      })
    }
}

func (pr PolicyResponse) getEmptyStatement(effect string) events.IAMPolicyStatement {
	/* Returns an empty statement object prepopulated with the correct action and the
    desired effect. */
    var statement events.IAMPolicyStatement = events.IAMPolicyStatement{
        Action: []string{"execute-api:Invoke"},
        Effect: effect.substring(0,1).toUpperCase() + effect.substring(1).toLowerCase(),
        Resource: []string{},
    };
	
    return statement;
}
func validateToken(region, token string) (string, error) {
	// KEYS URL -- REPLACE WHEN CHANGING IDENTITY PROVIDER
	keysUrl = fmt.Sprintf("https://cognito-idp.${%s}.amazonaws.com/${userPoolId}/.well-known/jwks.json", region, token)

	if res, err := http.Get(keysUrl); err != nil {
		return "", err
	} else {
		return "", nil
	}
}


func main() {
	lambdaHandler := lambda.NewHandler(func(ctx context.Context, request *events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
		// Parse the input for the parameter values
		methodArn := strings.Split(request.MethodArn, ":")
	
		if validateToken(request.AuthorizationToken, methodArn[2])
		
		apiGatewayArn := strings.Split(methodArn[5], "/")
		policy := PolicyResponse {
		// Save the ARN parts
			accountId: methodArn[4],
			region: methodArn[3],
			route: methodArn[2],
			stage: apiGatewayArn[1],
			apiId: methodArn[0],
		}

		// Perform authorization to return the Allow policy for correct parameters
    	// and the 'Unauthorized' error, otherwise.
		if headers["HeaderAuth1"] == "headerValue1" && queryParams["QueryString1"] == "queryValue1" {
			return generateAllow("me", request.MethodArn)
		}
		return nil, fmt.Errorf("not available")
	})

	lambda.Start(lambdaHandler)
}