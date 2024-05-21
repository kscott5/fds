resource aws_dynamodb_table "favoritetable" {
    name         = "${var.app_prefix}Favorite"
  billing_mode = "PROVISIONED"
  hash_key     = "id"
  
  read_capacity  = 5
  write_capacity = 5
  attribute {
    name = "userid"
    type = "S"
  }
}