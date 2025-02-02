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

variable "app_name" {
  description = "Application name prefix for resources"
  type        = string
  default     = "canvis-collab"
}
