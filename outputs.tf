output "cognito_user_pool_id" {
  description = "The ID of the Cognito User Pool"
  value       = module.cognito.user_pool_id
}

output "cognito_client_id" {
  description = "The ID of the Cognito User Pool Client"
  value       = module.cognito.client_id
}

output "cognito_client_secret" {
  description = "The Client Secret of the Cognito User Pool Client"
  value       = module.cognito.client_secret
  sensitive   = true
}

output "cognito_domain" {
  description = "The Cognito Domain URL"
  value       = module.cognito.domain_url
}
