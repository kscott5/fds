# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: MIT-0

data "archive_file" "userfunctions_lambda_zip" {
  type        = "zip"
  output_path = "/tmp/userfunctions_lambda.zip"
  source_dir  = "src/users/dist"
}

resource "aws_lambda_function" "userfunctions_lambda" {
  filename         = data.archive_file.userfunctions_lambda_zip.output_path
  function_name    = "${var.workshop_stack_base_name}_userfunctions_lambda"
  description      = "Handler for all users related operations"
  role             = aws_iam_role.userfunctions_lambda_role.arn
  handler          = "lambda_function.lambda_handler"
  source_code_hash = data.archive_file.userfunctions_lambda_zip.output_base64sha256
  runtime          = var.lambda_runtime
  timeout          = var.lambda_timeout
  tracing_config {
    mode = var.lambda_tracing_config
  }
  environment {
    variables = {
      USERS_TABLE = aws_dynamodb_table.users_table.id
    }
  }
}

output "userfunctions_lambda" {
  value = aws_lambda_function.userfunctions_lambda.function_name
}
