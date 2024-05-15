package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/aws/aws-lambda-go/lambdacontext" // IMPORTANT: package level init() in use.

	"go.uber.org/zap"
)

const (
	// The policy version used for the evaluation. This should always be '2012-10-17'
	version = "2012-10-17"
	// The regular expression used to validate resource paths for the policy
	pattern = `^[/.a-zA-Z0-9-\*]+$`
)

var (
	UserPoolId     = os.Getenv("FDS_USER_POOL_ID")
	AppClientId    = os.Getenv("FDS_APPLICATION_CLIENT_ID")
	AdminGroupName = os.Getenv("FDS_ADMIN_GROUP_NAME")
	
	logger    *zap.Logger
	HttpVerb = map[string]string{
		"GET":     "GET",
		"POST":    "POST",
		"PUT":     "PUT",
		"PATCH":   "PATCH",
		"HEAD":    "HEAD",
		"DELETE":  "DELETE",
		"OPTIONS": "OPTIONS",
		"ALL":     "*",
	}
)

type Method struct {
	ResourceArn string
	Conditions  []string
}
type LocalAuthorizerResponse struct {
	events.APIGatewayCustomAuthorizerResponse

	allowMethods []Method `json:"-"`
	denyMethods  []Method `json:"-"`

	AccountId string `json:"-"`
	Region    string `json:"-"`
	Route     string `json:"-"`
	Stage     string `json:"-"`
	ApiId     string `json:"-"`
}

func (pr *LocalAuthorizerResponse) addMethod(effect, verb, resource string, conditions []string) error {
	/* Adds a method to the internal lists of allowed or denied methods. Each object in
	   the internal list contains a resource ARN and a condition statement. The condition
	   statement can be null. */
	if found := HttpVerb[verb]; found != "" && verb != "*" {
		return fmt.Errorf("invalid HTTP verb '%s'", verb)
	}

	if found, _ := regexp.Match(pattern, []byte(resource)); !found {
		return fmt.Errorf("invalid resource path: '%s'. path should match '%s'", resource, pattern)
	}
	if resource[0:1] == "/" {
		resource = resource[1:]
	}

	// https://pkg.go.dev/strings#Builder
	resourceArnBuilder := strings.Builder{}
	fmt.Fprintf(&resourceArnBuilder, "arn:aws:execute-api:%s:%s:%s/%s/%s/%s", pr.Region, pr.AccountId, pr.ApiId, pr.Stage, verb, resource)

	method := Method{ResourceArn: resourceArnBuilder.String(), Conditions: conditions}
	if strings.ToLower(effect) == "allow" {
		pr.allowMethods = append(pr.allowMethods, method)
	} else if strings.ToLower(effect) == "deny" {
		pr.denyMethods = append(pr.denyMethods, method)
	}

	return nil
}

func (pr *LocalAuthorizerResponse) AllowAllMethods() {
	//Adds a '*' allow to the policy to authorize access to all methods of an API
	pr.addMethod("Allow", HttpVerb["ALL"], "*", []string{})
}

func (pr *LocalAuthorizerResponse) DenyAllMethods() {
	//Adds a '*' allow to the policy to deny access to all methods of an API
	pr.addMethod("Deny", HttpVerb["ALL"], "*", []string{})
}

func (pr *LocalAuthorizerResponse) AllowMethod(verb, resource string) {
	/*Adds an API Gateway method (Http verb + Resource path) to the list of allowed\
	  methods for the policy';*/
	pr.addMethod("Allow", verb, resource, []string{})
}

func (pr *LocalAuthorizerResponse) DenyMethod(verb, resource string) {
	/*Adds an API Gateway method (Http verb + Resource path) to the list of denied\n' +
	  methods for the policy*/
	pr.addMethod("Deny", verb, resource, []string{})
}

func (pr *LocalAuthorizerResponse) getEmptyStatement(effect string) events.IAMPolicyStatement {
	/* Returns an empty statement object prepopulated with the correct action and the
	   desired effect. */

	var statement events.IAMPolicyStatement = events.IAMPolicyStatement{
		Action:   []string{"execute-api:Invoke"},
		Effect:   strings.Join([]string{strings.ToUpper(effect[0:1]), strings.ToLower(effect[1:])}, ""),
		Resource: []string{},
	}

	return statement
}

func (pr *LocalAuthorizerResponse) getStatementForEffect(effect string, methods []Method) []events.IAMPolicyStatement {
	/* This function loops over an array of objects containing a resourceArn and
	   conditions statement and generates the array of statements for the policy. */
	var statements []events.IAMPolicyStatement

	for _, v := range methods {
		statement := pr.getEmptyStatement(effect)
		statement.Resource = append(statement.Resource, v.ResourceArn)
		statements = append(statements, statement)
	}
	return statements
}

func (pr *LocalAuthorizerResponse) Build(principalId string) error {
	/*Generates the policy document based on the internal lists of allowed and denied
	  conditions. This will generate a policy with two main statements for the effect:
	  one statement for Allow and one statement for Deny.
	  Methods that includes conditions will have their own statement in the policy.*/
	if len(pr.allowMethods) == 0 && len(pr.denyMethods) == 0 {
		return fmt.Errorf("no statements defined for the policy")
	}

	pr.PrincipalID = principalId
	pr.PolicyDocument.Version = version
	pr.PolicyDocument.Statement = []events.IAMPolicyStatement{}

	var allowMethodsStatement = pr.getStatementForEffect("Allow", pr.allowMethods)
	var denyMethodsStatement = pr.getStatementForEffect("Deny", pr.denyMethods)
	var allMethodsStatement = append(allowMethodsStatement, denyMethodsStatement...)

	if len(allMethodsStatement) > 0 {
		pr.PolicyDocument.Statement = append(pr.PolicyDocument.Statement, allMethodsStatement...)
	}

	return nil
}

func GetWellKnownJwksKeys(region, userPoolId string)(map[string]interface{}, error) {
	logger.Info("GetWellKnowJwksKeys")

	// KEYS URL -- REPLACE WHEN CHANGING IDENTITY PROVIDER
	keysUrl := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, userPoolId)

	if res, err := http.Get(keysUrl); err != nil {
		return nil, err
	} else {
		size := res.ContentLength
		body := make([]byte, size)
		
		if size, err := res.Body.Read(body); err != nil {
			return nil, err
		} else if int64(size) == res.ContentLength {
			keys := make(map[string]interface{}, 1)
			json.Unmarshal(body, &keys)

			return keys, nil
		} else {
			return nil, err
		}
	}
}

func ValidateAuthToken(region, authToken string) (map[string]interface{}, error) {
	logger.Debug("validateAuthToken")
	logger.Debug(fmt.Sprintf("validateAuthToken: %s %s", region, authToken))

	if keys, err := GetWellKnownJwksKeys(region, UserPoolId); err != nil {
		return nil, err
	} else {

		rs256 := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}))
	
		_, err := rs256.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			header := token.Header
			for _, v := range keys {
				key := v.(map[string]string)
				if key["kid"] == header["kid"] {
					fmt.Println(header)
				}

			}
			return token, fmt.Errorf("%v", token)
		})

		if err != nil {
			fmt.Println(err)
			return nil, err
		}
			
		return nil, fmt.Errorf("not available")
	}
}

func main() {
	logger, _ := zap.NewDevelopment()
	logger.Info("FDS main authorizer")

	lambdaHandler := lambda.NewHandler(func(ctx context.Context, request *events.APIGatewayCustomAuthorizerRequest) (*LocalAuthorizerResponse, error) {
		logger.Info("FDS lambda.Start authorizer")

		// Parse the input for the parameter values
		// methodArn := []string{"arn", "aws", "execute-api", "{region}", "{accountid}" "{apiid}/{stage}/GET/request"}
		methodArn := strings.Split(request.MethodArn, ":")
		logger.Debug(fmt.Sprint(methodArn))
		
		if len(methodArn) < 6 {
			return &LocalAuthorizerResponse{}, fmt.Errorf("request method arn not available")
		}

		if token, err := ValidateAuthToken(methodArn[3], request.AuthorizationToken); err != nil {
			return &LocalAuthorizerResponse{}, err
		} else {
			apiGatewayArn := strings.Split(methodArn[5], "/")
			response := LocalAuthorizerResponse{
				// Save the ARN parts
				//PrincipalID: token["sub"],
				AccountId: methodArn[4],
				Region:    methodArn[3],
				Route:     methodArn[2],
				Stage:     apiGatewayArn[1],
				ApiId:     apiGatewayArn[0],
			}

			// *** Section 2 : authorization rules
			// Allow all public resources/methods explicitly

			var seperator = ""
			var singleResource = strings.Join([]string{"/users/", response.PrincipalID}, seperator)
			var multiResource = strings.Join([]string{"/users/", response.PrincipalID, "/*"}, seperator)

			response.AllowMethod(HttpVerb["GET"], singleResource)
			response.AllowMethod(HttpVerb["PUT"], singleResource)
			response.AllowMethod(HttpVerb["DELETE"], singleResource)
			response.AllowMethod(HttpVerb["GET"], multiResource)
			response.AllowMethod(HttpVerb["PUT"], multiResource)
			response.AllowMethod(HttpVerb["PUT"], multiResource)
			response.AllowMethod(HttpVerb["DELETE"], multiResource)

			
			// Look for admin group in Cognito groups
			// Assumption: admin group always has higher precedence
			// if found := token["cognito:groups"]; found && token["cognito:groups"][0] == adminGroupName {
			// 	// add administrative privileges
			// 	response.AllowMethod(HttpVerb["GET"], "users")
			// 	response.AllowMethod(HttpVerb["GET"], "users/*")
			

			// 	response.AllowMethod(HttpVerb["DELETE"], "users")
			// 	response.AllowMethod(HttpVerb["DELETE"], "users/*")
			// 	response.AllowMethod(HttpVerb["PUT"], "users")
			// 	response.AllowMethod(HttpVerb["PUT"], "users/*")
			// }

			//response.Build()
			return &response, fmt.Errorf("not available: token-> %s", token)
		}
	})

	lambda.Start(lambdaHandler)
}
