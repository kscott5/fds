# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: MIT-0

data "archive_file" "userfunctions_lambda_auth_zip" {
  type        = "zip"
  output_path = "../dist/${var.app_prefix}.lambda.auth.go.zip"
  source_file = "../dist/auth/bootstrap"
}

resource "aws_lambda_function" "userfunctions_lambda_auth" {
  filename      = "${data.archive_file.userfunctions_lambda_auth_zip.output_path}"
  function_name = "${var.app_prefix}Authozier"
  description   = "Handler for Lambda authorizer"
  role          = "${aws_iam_role.userfunctions_lambda_auth_role.arn}"
  handler       = "boostrap"
  source_code_hash = "${data.archive_file.userfunctions_lambda_auth_zip.output_base64sha256}"
  runtime = var.lambda_runtime[1]
  timeout = var.lambda_timeout
  environment {
    variables = {
      USER_POOL_ID = aws_cognito_user_pool.user_pool.id
      APPLICATION_CLIENT_ID = aws_cognito_user_pool_client.user_pool_client.id
      ADMIN_GROUP_NAME = var.user_pool_admin_group_name
    }
  }
}

output "userfunctions_lambda_auth" {
  value = aws_lambda_function.userfunctions_lambda_auth.function_name
}
