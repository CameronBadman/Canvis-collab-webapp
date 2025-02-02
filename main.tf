terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

module "dynamodb" {
  source = "./terraform/dynamodb"
  
  app_name    = var.app_name
  environment = var.environment
  region_name = var.region_name
}

module "cognito" {
  source = "./terraform/cognito"
  
  app_name     = var.app_name
  environment  = var.environment
  callback_urls = var.callback_urls
  logout_urls  = var.logout_urls
}
