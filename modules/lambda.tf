# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: MIT-0

data "archive_file" "getusers_lambda_zip" {
  type        = "zip"
  output_path = "../dist/${var.app_prefix}.lambda.pi.zip"
  source_dir  = "../src/pi"
}

resource "aws_lambda_function" "getusers" {
  filename         = data.archive_file.getusers_lambda_zip.output_path
  function_name    = "${var.app_prefix}GetUsersPi"
  role             = aws_iam_role.lambda_role.arn
  handler          = "userservice.lambda_handler"
  source_code_hash = data.archive_file.getusers_lambda_zip.output_base64sha256
  runtime          = var.lambda_runtime[0]
  timeout          = var.lambda_timeout
  tracing_config {
    mode = var.lambda_tracing_config
  }
  environment {
    variables = {
      FDS_APPS_USERS_TABLE = aws_dynamodb_table.users_table.id
    }
  }
}

output "getusers_lambda" {
  #value = aws_lambda_function.getusers.function_name
  value = "${var.arn_aws_lambda_base}:${var.region}:${var.account_id}:function:${aws_lambda_function.getusers.function_name}"
}
