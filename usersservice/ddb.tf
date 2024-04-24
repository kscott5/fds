# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: MIT-0

resource "aws_dynamodb_table" "users_table" {
  name         = "${var.workshop_stack_base_name}.users"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "_id"

  attribute {
    name = "_id"
    type = "S"
  }
  attribute {
    name = "userid"
    type = "S"
  }


  attribute {
    name = "username"
    type = "S"
  }

  attribute {
    name = "surname"
    type = "S"
  }
  attribute {
    name = "givennames"
    type = "S"
  }
}
output "users_table" {
  value = aws_dynamodb_table.users_table.id
}
