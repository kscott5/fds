# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_rest_api
# https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-integration-async.html
resource "aws_api_gateway_rest_api" "getusers" {
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
            authorizerUri                = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/arn:aws:lambda:${var.region}:${data.aws_caller_identity.current.account_id}:function:${aws_lambda_function.user_lambda_auth.function_name}/invocations"
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
            uri                 = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${aws_lambda_function.getusers_go.arn}/invocations"
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
            uri                 = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${aws_lambda_function.getuser_go.arn}/invocations"
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
            uri                 = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${aws_lambda_function.getusers.arn}/invocations"
          }
        }
      }
    }
  })
}

#https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_deployment
resource "aws_api_gateway_deployment" "getusers" {
  rest_api_id = aws_api_gateway_rest_api.getusers.id

  triggers = {
    redeployment = sha1(jsonencode(aws_api_gateway_rest_api.getusers.body))
  }
  lifecycle {
    create_before_destroy = true
  }
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_stage
resource "aws_api_gateway_stage" "getusers" {
  rest_api_id          = aws_api_gateway_rest_api.getusers.id
  stage_name           = "dev"
  deployment_id        = aws_api_gateway_deployment.getusers.id
  xray_tracing_enabled = true
}
resource "aws_lambda_permission" "api_getusers_go" {
  statement_id  = "${var.app_prefix}LambdaPermission"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.getusers_go.function_name
  principal     = "apigateway.${var.region}.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.getusers.execution_arn}/*/*/*"
}
resource "aws_lambda_permission" "api_getuser_go" {
  statement_id  = "${var.app_prefix}LambdaPermission"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.getuser_go.function_name
  principal     = "apigateway.${var.region}.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.getusers.execution_arn}/*/*/*"
}
resource "aws_lambda_permission" "api_putuser_go" {
  statement_id  = "${var.app_prefix}LambdaPermission"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.putuser_go.function_name
  principal     = "apigateway.${var.region}.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.getusers.execution_arn}/*/*/*"
}
output "aws_api_gateway_stage" {
  value = aws_api_gateway_stage.getusers.invoke_url
}