# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: MIT-0

resource "aws_dynamodb_table" "orders_table" {
  name         = "${var.app_prefix}Orders"
  billing_mode = "PROVISIONED"
  hash_key     = "userid"
  range_key    = "orderid"

  read_capacity  = 5
  write_capacity = 5
  attribute {
    name = "userid"
    type = "S"
  }
  attribute {
    name = "orderid"
    type = "S"
  }

}

output "orders_table" {
  value = aws_dynamodb_table.orders_table.id
}