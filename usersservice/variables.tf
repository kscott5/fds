# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: MIT-0

variable "region" {}
variable "workshop_stack_base_name" {
  default = "fds.apps"
}
variable "lambda_memory" {
  default = "128"
}
variable "lambda_runtime" {
  default = "python3.12"
}
variable "lambda_timeout" {
  default = "100"
}
variable "lambda_tracing_config" {
  default = "Active"
}
