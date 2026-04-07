variable "aws_region" {
  description = "AWS region for managed resources"
  type        = string
  default     = "us-east-1"
}

/*
 * AWS tag values
 */

variable "app_customer" {
  description = "customer name to use for the itse_app_customer tag"
  type        = string
  default     = "gtis"
}

variable "app_environment" {
  description = "environment name to use for the itse_app_environment tag, e.g. staging, production"
  type        = string
  default     = "production"
}

variable "app_name" {
  description = "app name to use for the itse_app_name tag"
  type        = string
  default     = "app-monitoring-archiver"
}

/*
 * GitHub OIDC provider parameters
 */

variable "github_oidc_provider_arn" {
  description = <<-EOT
    ARN of the OIDC provider for GitHub in AWS IAM, used for GitHub Actions to authenticate to AWS. The provider
    can be created in Terraform using the `aws_iam_openid_connect_provider` resource. Specify the URL as
    "https://token.actions.githubusercontent.com" and the client_id_list as ["sts.amazonaws.com"].
  EOT
  type        = string
}

variable "github_repository" {
  description = <<-EOT
    GitHub repository that should be granted access to the OIDC provider for GitHub. Format should be 'owner/repo'.
  EOT
  type        = string
}
