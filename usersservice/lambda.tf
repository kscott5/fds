# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: MIT-0

data "archive_file" "userfunctions_lambda_zip" {
  type        = "zip"
  output_path = "./dist/${var.workshop_stack_base_name}.lambda.getusers.zip"
  source_dir  = "./src/users"
}

resource "aws_lambda_function" "getusers" {
  filename         = data.archive_file.userfunctions_lambda_zip.output_path
  function_name    = "tablescan"
  description      = "${var.workshop_stack_base_name}.users.table scan"
  role             = aws_iam_role.userfunctions_lambda_role.arn
  handler          = "tablescan.lambda_handler"
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
  value = aws_lambda_function.getusers.function_name
}
