# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_rest_api
# https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-integration-async.html
resource "aws_api_gateway_rest_api" "rest_api" {
  name = "${var.app_prefix}GetUsers"
  endpoint_configuration {
    types = ["REGIONAL"]
  }
  body = jsonencode({
    openapi = "3.0.1"
    info = {
      title   = "FDS RestAPI Get Users (Python)",
      version = "1.0"
    },
    components = {
      securitySchemes = {
        lambdaTokenAuthorizer = {
          type                         = "apiKey"
          name                         = "Authorization"
          in                           = "header"
          x-amazon-apigateway-authtype = "custom"
          x-amazon-apigateway-authorizer = {
            authorizerUri                = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/arn:aws:lambda:${var.region}:${data.aws_caller_identity.current.account_id}:function:${aws_lambda_function.users_authorizer.function_name}/invocations"
            authorizerResultTtlInSeconds = 300
            type                         = "token"
          }
        }
      }
    },
    paths = {
      "/users" = {
        get = {
          security = [
            {
              "lambdaTokenAuthorizer" : []
            }
          ]
          x-amazon-apigateway-integration = {
            httpMethod          = "POST"
            type                = "aws_proxy"
            passthroughBehavior = "WHEN_NO_MATCH"
            uri                 = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${aws_lambda_function.getusers.arn}/invocations"
          }
        },
        put = {
          security = [
            {
              "lambdaTokenAuthorizer" : []
            }
          ]
          x-amazon-apigateway-integration = {
            httpMethod          = "POST"
            type                = "aws_proxy"
            passthroughBehavior = "WHEN_NO_MATCH"
            uri                 = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${aws_lambda_function.putuser.arn}/invocations"
          }
        }
      },

      "/users/{_id}" = {
        get = {
          security = [
            {
              "lambdaTokenAuthorizer" : []
            }
          ]
          x-amazon-apigateway-integration = {
            httpMethod          = "POST"
            type                = "aws_proxy"
            passthroughBehavior = "WHEN_NO_MATCH"
            uri                 = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${aws_lambda_function.getuser.arn}/invocations"
          }
        }
      },

      "/user" = {
        put = {
          security = [
            {
              "lambdaTokenAuthorizer" : []
            }
          ]
          x-amazon-apigateway-integration = {
            httpMethod          = "POST"
            type                = "aws_proxy"
            passthroughBehavior = "WHEN_NO_MATCH"
            uri                 = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${aws_lambda_function.putuser.arn}/invocations"
          }
        }
      }
    }
  })
}

#https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_deployment
resource "aws_api_gateway_deployment" "rest_api" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id

  triggers = {
    redeployment = sha1(jsonencode(aws_api_gateway_rest_api.rest_api.body))
  }
  lifecycle {
    create_before_destroy = true
  }
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_stage
resource "aws_api_gateway_stage" "rest_api" {
  rest_api_id          = aws_api_gateway_rest_api.rest_api.id
  stage_name           = "dev"
  deployment_id        = aws_api_gateway_deployment.rest_api.id
  xray_tracing_enabled = true
}

resource "aws_lambda_permission" "allow_api_on_getusers" {
  statement_id  = "${var.app_prefix}LambdaPermission"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.getusers.function_name
  principal     = "apigateway.${var.region}.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.rest_api.execution_arn}/*/*/*"
}
resource "aws_lambda_permission" "allow_api_on_getuser" {
  statement_id  = "${var.app_prefix}LambdaPermission"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.getuser.function_name
  principal     = "apigateway.${var.region}.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.rest_api.execution_arn}/*/*/*"
}
resource "aws_lambda_permission" "allow_api_on_putuser" {
  statement_id  = "${var.app_prefix}LambdaPermission"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.putuser.function_name
  principal     = "apigateway.${var.region}.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.rest_api.execution_arn}/*/*/*"
}

resource "aws_lambda_permission" "allow_api_on_authorizer" {
  statement_id  = "${var.app_prefix}LambdaPermission"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.users_authorizer.function_name
  principal     = "apigateway.${var.region}.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.rest_api.execution_arn}/*/*/*"
}
output "rest_api" {
  value = aws_api_gateway_stage.rest_api.invoke_url
}