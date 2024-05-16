# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: MIT-0

data "archive_file" "placeorder_lambda_zip" {
  type        = "zip"
  output_path = "../dist/${var.app_prefix}.lambda.orders.zip"
  source_file = "../dist/orders/bootstrap"
}

resource "aws_lambda_function" "placeorder" {
  filename         = data.archive_file.placeorder_lambda_zip.output_path
  function_name    = "${var.app_prefix}PlaceOrder"
  role             = aws_iam_role.lambda_role.arn
  handler          = "bootstrap"
  source_code_hash = data.archive_file.placeorder_lambda_zip.output_base64sha256
  runtime          = var.lambda_runtime[1]
  architectures    = var.architectures
  timeout          = var.lambda_timeout
  tracing_config {
    mode = var.lambda_tracing_config
  }
  environment {
    variables = {
      FDS_APPS_ORDERS_TABLE = aws_dynamodb_table.orders_table.id
    }
  }
}

output "placeorde_lambda" {
  #value = aws_lambda_function.getusers.function_name
  value = "${var.arn_aws_lambda_base}:${var.region}:${var.account_id}:function:${aws_lambda_function.placeorder.function_name}"
}
