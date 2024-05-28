resource aws_dynamodb_table "addresstable" {
    name         = "${var.app_prefix}Address"
  billing_mode = "PROVISIONED"
  hash_key     = "userid"
  
  read_capacity  = 5
  write_capacity = 5
  
  attribute {
    name = "id"
    type = "S"
  }
}