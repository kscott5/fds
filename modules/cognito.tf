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
    "ALLOW_USER_PASSWORD_AUTH",
    "ALLOW_USER_SRP_AUTH",
    "ALLOW_REFRESH_TOKEN_AUTH"
  ]
  generate_secret                      = false
  prevent_user_existence_errors        = "ENABLED"
  refresh_token_validity               = 30
  supported_identity_providers         = ["COGNITO"]
  user_pool_id                         = aws_cognito_user_pool.user_pool.id
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["code"]
  allowed_oauth_scopes                 = ["email", "openid"]
  callback_urls                        = ["http://localhost"]
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

output "user_pool_client" {
  value = aws_cognito_user_pool_client.user_pool_client.id
}

output "user_pool_admin_group" {
  value = var.user_pool_admin_group_name
}

output "cognito_login_url" {
  value = "https://${aws_cognito_user_pool_client.user_pool_client.id}.auth.${var.region}.amazoncognito.com/login?client_id=${aws_cognito_user_pool_client.user_pool_client.id}&response_type=code&redirect_uri=http://repost.aws"
}

output "cognito_login_auth_command" {
  value = "aws cognito-idp initiate-auth --auth-flow USER_PASSWORD_AUTH --client-id ${aws_cognito_user_pool_client.user_pool_client.id} --region ${var.region} --auth-parameters USERNAME=<user@example.com>,PASSWORD=<password>"
}
