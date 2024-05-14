module github.com/kscott5/fds/authorizer

go 1.22.1

replace github.com/kscott5/fds/authorizer => ./authorizer

require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/golang-jwt/jwt/v5 v5.2.1
)
