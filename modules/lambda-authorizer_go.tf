# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: MIT-0

data "archive_file" "user_lambda_auth_zip" {
  type        = "zip"
  output_path = "../dist/${var.app_prefix}.lambda.auth.go.zip"
  source_file = "../dist/auth/bootstrap"
}

resource "aws_lambda_function" "user_lambda_auth" {
  filename         = data.archive_file.user_lambda_auth_zip.output_path
  function_name    = "${var.app_prefix}Authozier"
  role             = aws_iam_role.user_lambda_auth_role.arn
  handler          = "boostrap"
  source_code_hash = data.archive_file.user_lambda_auth_zip.output_base64sha256
  runtime          = var.lambda_runtime[1]
  architectures = var.architectures
  timeout          = var.lambda_timeout
  tracing_config {
    mode = var.lambda_tracing_config
  }
  environment {
    variables = {
      FDS_USER_POOL_ID          = aws_cognito_user_pool.user_pool.id
      FDS_APPLICATION_CLIENT_ID = aws_cognito_user_pool_client.user_pool_client.id
      FDS_ADMIN_GROUP_NAME      = var.user_pool_admin_group_name
    }
  }
}

output "user_lambda_auth" {
  value = aws_lambda_function.user_lambda_auth.function_name
}
