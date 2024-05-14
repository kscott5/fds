# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: MIT-0

resource "aws_iam_role" "userfunctions_lambda_auth_role" {
  name               = "${var.app_prefix}LambdaAuthRole"
  description        = "Lambda function authorizer IAM role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_policy" "userfunctions_lambda_auth_role_policy" {
  name        = "${var.app_prefix}LambdaAuthRolePolicy"
  description = "Lambda function authorizer policy"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": [
        "xray:*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_policy_attachment" "userfunctions_lambda_auth_attach" {
  name       = "${var.app_prefix}LambdaAuthAttachment"
  roles      = ["${aws_iam_role.userfunctions_lambda_auth_role.name}"]
  policy_arn = aws_iam_policy.userfunctions_lambda_auth_role_policy.arn
}

