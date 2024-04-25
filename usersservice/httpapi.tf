# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_integration
resource "aws_api_gateway_rest_api" "getusers" {
    name = "${var.workshop_stack_base_name}.rest.api.users.tablescan"
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_resource
resource "aws_api_gateway_resource" "getusers" {
    rest_api_id = aws_api_gateway_rest_api.getusers.id
    parent_id = aws_api_gateway_rest_api.getusers.root_resource_id
    path_part = "somepart"
}

resource  "aws_api_gateway_method" "getusers" {
    rest_api_id = aws_api_gateway_rest_api.getusers.id
    resource_id = aws_api_gateway_resource.getusers.id
    http_method = "GET"
    authorization = "NONE"
}
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_integration
resource "aws_api_gateway_integration" "getusers" {
  rest_api_id = aws_api_gateway_rest_api.getusers.id
  resource_id = aws_api_gateway_resource.getusers.id
  http_method = aws_api_gateway_method.getusers.http_method
  uri = aws_lambda_function.getusers.invoke_arn
  integration_http_method = "POST"
  type = "AWS_PROXY"
  
}

#https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_deployment
resource "aws_api_gateway_deployment" "getusers" {
    rest_api_id = aws_api_gateway_rest_api.getusers.id
    triggers = {
        redeployment = sha1(jsondecode(
            aws_api_gateway_resource.getusers.id,
            aws_api_gateway_method.getusers.id,
            aws_api_gateway_integration.getusers.id))
    }

    lifecycle {
      create_before_destroy = true
    }
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_stage
resource "aws_api_gateway_stage" "getusers" {

  
}
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/apigatewayv2_api
resource "aws_apigatewayv2_api" "getusers" {
  name          = "${var.workshop_stack_base_name}.api.users.tablescan"
  protocol_type = "HTTP"
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/apigatewayv2_route
resource "aws_apigatewayv2_route" "getusers" {
  api_id    = aws_apigatewayv2_api.getusers.id
  route_key = "GET /users/tablescan"
  target = "integrations/${aws_apigatewayv2_integration.getusers.id}"
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/apigatewayv2_integration
resource "aws_apigatewayv2_integration" "getusers" {
  api_id             = aws_apigatewayv2_api.getusers.id
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.getusers.invoke_arn
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/apigatewayv2_deployment
resource "aws_apigatewayv2_deployment" "getusers" {
  api_id = aws_apigatewayv2_api.getusers.id
  

  triggers = {
    redeployment = sha1(join(",", tolist([
      jsonencode(aws_apigatewayv2_api.getusers),
      jsonencode(aws_apigatewayv2_route.getusers),
      jsonencode(aws_apigatewayv2_integration.getusers),
    ])))
  }

  lifecycle {
    create_before_destroy = true
  }
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/apigatewayv2_stage
resource "aws_apigatewayv2_stage" "getusers" {
  api_id        = aws_apigatewayv2_api.getusers.id
  #deployment_id = aws_apigatewayv2_deployment.getusers.id
  name          = "dev"
  auto_deploy = true  
}

resource "aws_lambda_permission" "api_getusers" {
    statement_id = "AllowHttpGetUsers"
    action =  "lambda:InvokeFunction"
    function_name = aws_lambda_function.getusers.function_name
    principal = "apigateway.amazonaws.com"
    source_arn = aws_lambda_function.getusers.arn
}

output "aws_apigatewayv2_stage" {
  value = aws_apigatewayv2_stage.getusers.invoke_url
  description = "Available staging url located on the HTTP v2 Gateway"
}
output "aws_apigatewayv2_api" {
  value = aws_apigatewayv2_api.getusers.api_endpoint
}