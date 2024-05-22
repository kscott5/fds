# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: MIT-0

resource "aws_cognito_user_pool" "user_pool" {
  name = "${var.app_prefix}UserPool"
  admin_create_user_config {
    allow_admin_create_user_only = false
  }
  auto_verified_attributes = ["email"]
  schema {
    name                = "name"
    attribute_data_type = "String"
    mutable             = true
    required            = true
  }
  schema {
    name                = "email"
    attribute_data_type = "String"
    mutable             = true
    required            = true
  }
  username_attributes = ["email"]
  tags = {
    Name = "User Pool"
  }
  lifecycle {
    ignore_changes = [
      schema
    ]
  }
}

resource "aws_cognito_user_pool_client" "user_pool_client" {
  name = "${var.app_prefix}UserPoolClient"
  explicit_auth_flows = [
    "ALLOW_ADMIN_USER_PASSWORD_AUTH",
    "ALLOW_USER_PASSWORD_AUTH",
    "ALLOW_USER_SRP_AUTH",
    "ALLOW_REFRESH_TOKEN_AUTH"
  ]
  generate_secret               = true
  prevent_user_existence_errors = "ENABLED"

  token_validity_units {
    refresh_token = "days"
    access_token  = "minutes"
    id_token      = "minutes"
  }

  refresh_token_validity               = 30
  id_token_validity                    = 60
  access_token_validity                = 60
  supported_identity_providers         = ["COGNITO"]
  user_pool_id                         = aws_cognito_user_pool.user_pool.id
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["code", "implicit"]
  allowed_oauth_scopes                 = ["email", "openid", "aws.cognito.signin.user.admin"]
  callback_urls                        = ["http://localhost:8080"]

}

resource "aws_cognito_user_pool_domain" "user_pool_domain" {
  domain       = aws_cognito_user_pool_client.user_pool_client.id
  user_pool_id = aws_cognito_user_pool.user_pool.id
}

resource "aws_cognito_user_group" "api_administrator_user_pool_group" {
  name         = var.user_pool_admin_group_name
  user_pool_id = aws_cognito_user_pool.user_pool.id
  description  = "FDS User group for API Administrators"
  precedence   = 0
}

output "user_pool" {
  value = aws_cognito_user_pool.user_pool.id
}
output "user_pool_client_id" {
  value = aws_cognito_user_pool_client.user_pool_client.id
}

output "user_pool_admin_group" {
  value = var.user_pool_admin_group_name
}

output "cognito_login_url" {
  value = "https://${aws_cognito_user_pool_client.user_pool_client.id}.auth.${var.region}.amazoncognito.com/oauth2/authorize?client_id=${aws_cognito_user_pool_client.user_pool_client.id}&response_type=token&state=request&scope=aws.cognito.signin.user.admin+openid+email&redirect_uri=http://localhost:8080"
}

output "cognito_login_auth_command" {
  value = "aws cognito-idp initiate-auth --auth-flow USER_PASSWORD_AUTH --client-id ${aws_cognito_user_pool_client.user_pool_client.id} --region ${var.region} --auth-parameters USERNAME=<user@example.com>,PASSWORD=<password>"
}
