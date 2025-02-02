
variable "app_name" {
  description = "Application name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
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
