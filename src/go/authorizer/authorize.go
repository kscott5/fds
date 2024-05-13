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

type Method struct { 
	ResourceArn string
	Conditions []string
}
type PolicyResponse struct {
	events.APIGatewayCustomAuthorizerResponse
	
	allowMethods []Method
	denyMethods  []Method

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

    if found, _ := regexp.Match(pattern,[]byte(resource)); !found {
         return fmt.Errorf("Invalid resource path: '%s'. Path should match '%s'", resource, pattern);
    }
    if resource[0:1] == "/" {
        resource = resource[1:]
    }

	// https://pkg.go.dev/strings#Builder
    resourceArnBuilder := strings.Builder{}
	fmt.Fprintf(&resourceArnBuilder, "arn:aws:execute-api:%s:%s:%s/%s/%s/%s", pr.region, pr.accountId, pr.apiId, pr.stage, verb, resource)
	
	method := Method{ResourceArn: resourceArnBuilder.String(), Conditions: conditions}
    if (strings.ToLower(effect) == "allow") {
		pr.allowMethods = append(pr.allowMethods, method)
    } else if (strings.ToLower(effect) == "deny") {
        pr.denyMethods = append(pr.denyMethods, method)
    }

	return nil
}


func (pr *PolicyResponse) allowAllMethods() {
    //Adds a '*' allow to the policy to authorize access to all methods of an API
    pr.addMethod("Allow", httpVerb["ALL"], "*", []string{})
}

func (pr *PolicyResponse) denyAllMethods() {
    //Adds a '*' allow to the policy to deny access to all methods of an API
    pr.addMethod("Deny", httpVerb["ALL"], "*", []string{})
 }

func (pr *PolicyResponse) allowMethod(verb, resource string) {
    /*Adds an API Gateway method (Http verb + Resource path) to the list of allowed\
    methods for the policy';*/
    pr.addMethod("Allow", verb, resource, []string{})
}

func (pr *PolicyResponse) denyMethod(verb, resource string) {
    /*Adds an API Gateway method (Http verb + Resource path) to the list of denied\n' +
    methods for the policy*/
    pr.addMethod("Deny", verb, resource, []string{})
}

func (pr *PolicyResponse) getEmptyStatement(effect string) events.IAMPolicyStatement {
	/* Returns an empty statement object prepopulated with the correct action and the
    desired effect. */
	
    var statement events.IAMPolicyStatement = events.IAMPolicyStatement{
        Action: []string{"execute-api:Invoke"},
        Effect: strings.Join([]string{ strings.ToUpper(effect[0:1]), strings.ToLower(effect[1:])}, ""),
        Resource: []string{},
    };
	
    return statement;
}

func (pr *PolicyResponse) getStatementForEffect(effect string, methods []Method) []events.IAMPolicyStatement {
    /* This function loops over an array of objects containing a resourceArn and
    conditions statement and generates the array of statements for the policy. */
    var statements []events.IAMPolicyStatement
    
    for _, v := range methods {
		statement := pr.getEmptyStatement(effect)
		statement.Resource = append( statement.Resource, v.ResourceArn)
		statements = append(statements, statement)
    }
    return statements;
}

func (pr * PolicyResponse)  build() {
    /*Generates the policy document based on the internal lists of allowed and denied
    conditions. This will generate a policy with two main statements for the effect:
    one statement for Allow and one statement for Deny.
    Methods that includes conditions will have their own statement in the policy.*/
    if ((pr.allowMethods == nil || pr.allowMethods.length == 0) &&
      (this.denyMethods === null || this.denyMethods.length == 0)) {
      throw Error('No statements defined for the policy');
    }
    var policy = {
        'principalId': this.principalId,
        'policyDocument': {
            'Version': this.version,
            'Statement': []
        }
    };

    var allowMethodsStatement = this.getStatementForEffect('Allow', this.allowMethods)
    var denyMethodsStatement = this.getStatementForEffect('Deny', this.denyMethods)
    var allMethodsStatement = allowMethodsStatement.concat(denyMethodsStatement);

    if (allMethodsStatement != null) {
      policy['policyDocument']['Statement'] = allMethodsStatement;
    }
    console.log(JSON.stringify(policy))
    return policy;
}

func validateToken(region, token string) (string, error) {
	// KEYS URL -- REPLACE WHEN CHANGING IDENTITY PROVIDER
	keysUrl := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, userPoolId)

	if res, err := http.Get(keysUrl); err != nil {
		return "", err
	} else {
		return "", nil
	}
}

func main() {
	lambdaHandler := lambda.NewHandler(func(ctx context.Context, request *events.APIGatewayCustomAuthorizerRequest) (*events.APIGatewayCustomAuthorizerResponse, error) {
		response := events.APIGatewayCustomAuthorizerResponse{}

		// Parse the input for the parameter values
		methodArn := strings.Split(request.MethodArn, ":")
	
		if len(methodArn) < 6 {
			return &response, fmt.Errorf("request method arn not available")
		}

		if principalId, err := validateToken(request.AuthorizationToken, methodArn[2]); err != nil {

		} else {	
			apiGatewayArn := strings.Split(methodArn[5], "/")
			policy := PolicyResponse {
			// Save the ARN parts
				accountId: methodArn[4],
				region: methodArn[3],
				route: methodArn[2],
				stage: apiGatewayArn[1],
				apiId: methodArn[0],
			}

			// *** Section 2 : authorization rules
			// Allow all public resources/methods explicitly

			var seperator = ""
			var singleResource = strings.Join([]string{"/users/", principalId}, seperator)
			var multiResource = strings.Join([]string{"/users/", principalId, "/*"}, seperator)

			policy.allowMethod(httpVerb["GET"], singleResource)
			policy.allowMethod(httpVerb["PUT"], singleResource)
			policy.allowMethod(httpVerb["DELETE"], singleResource)
			policy.allowMethod(httpVerb["GET"], multiResource)
			policy.allowMethod(httpVerb["PUT"], multiResource)
			policy.allowMethod(httpVerb["PUT"],  multiResource)
			policy.allowMethod(httpVerb["DELETE"], multiResource)

			return &response, fmt.Errorf("not available")
		}
	})

	lambda.Start(lambdaHandler)
}