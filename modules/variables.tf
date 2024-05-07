# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: MIT-0

variable "account_id" {
  default = "154337954194"
}
variable "region" {
  default = "us-east-1"
}
variable "arn_aws_lambda_base" {
  default = "arn:aws:lambda"
}
variable "workshop_stack_base_name" {
  default = "FDSApps"
}
variable "lambda_memory" {
  default = "128"
}
variable "lambda_runtime" {
  type = list(string)
  default = [
    "python3.12",
    "provided.al2023"
  ]
}
variable "architectures" {
  type    = list(string)
  default = ["arm64"]
}
variable "lambda_timeout" {
  default = "100"
}
variable "lambda_tracing_config" {
  default = "Active"
}
