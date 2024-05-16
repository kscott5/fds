# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: MIT-0

data "archive_file" "getusers_lambda_zip" {
  type        = "zip"
  output_path = "../dist/${var.app_prefix}.lambda.go.zip"
  source_file = "../dist/bootstrap"
}

resource "aws_lambda_function" "getusers" {
  filename         = data.archive_file.getusers_lambda_zip.output_path
  function_name    = "${var.app_prefix}GetUsers"
  role             = aws_iam_role.lambda_role.arn
  handler          = "bootstrap"
  source_code_hash = data.archive_file.getusers_lambda_zip.output_base64sha256
  runtime          = var.lambda_runtime[1]
  architectures    = var.architectures
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

resource "aws_lambda_function" "getuser" {
  filename         = data.archive_file.getusers_lambda_zip.output_path
  function_name    = "${var.app_prefix}GetUser"
  role             = aws_iam_role.lambda_role.arn
  handler          = "bootstrap"
  source_code_hash = data.archive_file.getusers_lambda_zip.output_base64sha256
  runtime          = var.lambda_runtime[1]
  architectures    = var.architectures
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

resource "aws_lambda_function" "putuser" {
  filename         = data.archive_file.getusers_lambda_zip.output_path
  function_name    = "${var.app_prefix}PutUser"
  role             = aws_iam_role.lambda_role.arn
  handler          = "bootstrap"
  source_code_hash = data.archive_file.getusers_lambda_zip.output_base64sha256
  runtime          = var.lambda_runtime[1]
  architectures    = var.architectures
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

output "getuser_lambda" {
  #value = aws_lambda_function.getusers.function_name
  value = "${var.arn_aws_lambda_base}:${var.region}:${var.account_id}:function:${aws_lambda_function.getusers.function_name}"
}

output "putuser_lambda" {
  #value = aws_lambda_function.getusers.function_name
  value = "${var.arn_aws_lambda_base}:${var.region}:${var.account_id}:function:${aws_lambda_function.getusers.function_name}"
}
