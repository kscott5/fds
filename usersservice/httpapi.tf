# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_rest_api
# https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-integration-async.html
resource "aws_api_gateway_rest_api" "getusers" {
  name = "${var.workshop_stack_base_name}.api.getusers"
  body = jsonencode({
    openapi = "3.0.2"
    info = {
      title   = "FDS RestAPI Get Users",
      version = "1.0"
    },
    paths = {
      "/users" = {
        get = {
          x-amazon-apigateway-integration = {
            httpMethod = "POST"
            type       = "aws_proxy"
            uri        = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${aws_lambda_function.getusers.arn}/invocations"
          }
        }
      }
    }
  })
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_resource
# resource "aws_api_gateway_resource" "getusers" {
#   rest_api_id = aws_api_gateway_rest_api.getusers.id
#   parent_id   = aws_api_gateway_rest_api.getusers.root_resource_id
#   path_part   = "users"
# }

# resource "aws_api_gateway_method" "getusers" {
#   rest_api_id   = aws_api_gateway_rest_api.getusers.id
#   resource_id   = aws_api_gateway_resource.getusers.id
#   http_method   = "GET"
#   authorization = "NONE"
#}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_integration
# resource "aws_api_gateway_integration" "getusers" {
#   rest_api_id             = aws_api_gateway_rest_api.getusers.id
#   resource_id             = aws_api_gateway_resource.getusers.id
#   http_method             = aws_api_gateway_method.getusers.http_method
#   uri                     = aws_lambda_function.getusers.invoke_arn
#   integration_http_method = "POST"
#   type                    = "AWS_PROXY"
# }

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
  rest_api_id   = aws_api_gateway_rest_api.getusers.id
  stage_name    = "dev"
  deployment_id = aws_api_gateway_deployment.getusers.id
}

resource "aws_lambda_permission" "api_getusers" {
  statement_id  = "AllowHttpGetUsers"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.getusers.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = aws_lambda_function.getusers.arn
}
output "aws_api_gateway_stage" {
  value = aws_api_gateway_stage.getusers.invoke_url
}