output "cdk_access_key_id" {
  value = aws_iam_access_key.cdk.id
}

output "cdk_secret_access_key" {
  value     = aws_iam_access_key.cdk.secret
  sensitive = true
}

output "tfc_role_arn" {
  description = "ARN of the IAM role for Terraform Cloud OIDC authentication"
  value       = aws_iam_role.hcp.arn
}
