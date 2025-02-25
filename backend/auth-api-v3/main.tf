provider "aws" {
  region = "ap-southeast-2"
}

resource "aws_cognito_user_pool" "main" {
  name = "canvas_congito_user_pool"

  password_policy {
    minimum_length    = 8
    require_numbers   = true
    require_symbols   = true
    require_lowercase = true
    require_uppercase = true
  }

  mfa_configuration = "OFF"

  username_attributes = ["email"]

  auto_verified_attributes = ["email"]

  schema {
    name = "email"

    attribute_data_type = "String"
    required            = true
    mutable             = true

    string_attribute_constraints {
      min_length = 1
      max_length = 256
    }
  }
}


resource "aws_cognito_user_pool_client" "client" {

  name         = "canvas-client"
  user_pool_id = aws_cognito_user_pool.main.id

  generate_secret = true

  explicit_auth_flows = [

  ]

}

resource "aws_cognito_user_pool_domain" "main" {
  domain       = "canvas-client"
  user_pool_id = aws_cognito_user_pool.main.id
}

output "user_pool_id" {
  value = aws_cognito_user_pool.main.id
}


output "client_id" {
  value = aws_cognito_user_pool_client.client.id
}

output "client_secret" {
  value     = aws_cognito_user_pool_client.client.client_secret
  sensitive = true
}


