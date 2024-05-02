# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: MIT-0

resource "aws_dynamodb_table" "users_table" {
  name         = "${var.workshop_stack_base_name}.users"
  billing_mode = "PROVISIONED"
  hash_key     = "_id"

  read_capacity = 10
  write_capacity = 10
  attribute {
    name = "_id"
    type = "S"
  }
}

output "users_table" {
  value = aws_dynamodb_table.users_table.id
}