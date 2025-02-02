variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "ap-southeast-2"
}

variable "app_name" {
  description = "Application name"
  type        = string
  default     = "canvis-collab"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "development"
}

variable "region_name" {
  description = "Region name for tags"
  type        = string
  default     = "Sydney"
}

variable "callback_urls" {
  description = "Callback URLs for Cognito client"
  type        = list(string)
  default     = ["http://localhost:3000/callback"]
}

variable "logout_urls" {
  description = "Logout URLs for Cognito client"
  type        = list(string)
  default     = ["http://localhost:3000"]
}
