data "aws_region" "current" {}

resource "aws_cognito_user_pool" "main" {
  name = "${var.app_name}-${var.environment}"

  # Password policy for users
  password_policy {
    minimum_length    = 10
    require_lowercase = true
    require_numbers   = true
    require_symbols   = true
    require_uppercase = true
  }

  # Using email as primary way to sign in
  username_attributes = ["email"]
  
  # Auto verification
  auto_verified_attributes = ["email"]

  # Email verification settings
  verification_message_template {
    default_email_option = "CONFIRM_WITH_CODE"
    email_subject = "Welcome to Canvis Collab - Verify Your Account"
    email_message = "Welcome to Canvis! Your verification code is {####}"
  }

  # Required attributes
  schema {
    attribute_data_type = "String"
    name               = "email"
    required           = true
    string_attribute_constraints {
      min_length = 7
      max_length = 256
    }
  }

  schema {
    attribute_data_type = "String"
    name               = "name"
    required           = true
    string_attribute_constraints {
      min_length = 2
      max_length = 100
    }
  }

  # Account recovery setting
  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }

  tags = {
    Environment = var.environment
    Application = var.app_name
  }
}

resource "aws_cognito_user_pool_client" "main" {
  name = "${var.app_name}-client-${var.environment}"
  user_pool_id = aws_cognito_user_pool.main.id
  
  # Generate a client secret
  generate_secret = true
  
  # Auth flows
  explicit_auth_flows = [
    "ALLOW_USER_SRP_AUTH",
    "ALLOW_REFRESH_TOKEN_AUTH",
    "ALLOW_USER_PASSWORD_AUTH"
  ]

  # Token validity
  refresh_token_validity = 30
  access_token_validity  = 1
  id_token_validity     = 1

  token_validity_units {
    refresh_token = "days"
    access_token  = "hours"
    id_token     = "hours"
  }

  prevent_user_existence_errors = "ENABLED"

  # OAuth settings
  allowed_oauth_flows = ["implicit"]
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_scopes = [
    "email",
    "openid",
    "profile"
  ]
  
  callback_urls = var.callback_urls
  logout_urls   = var.logout_urls
}

resource "aws_cognito_user_pool_domain" "main" {
  domain       = "${var.app_name}-${var.environment}"
  user_pool_id = aws_cognito_user_pool.main.id
}
