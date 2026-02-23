provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      managed_by        = "terraform"
      workspace         = terraform.workspace
      itse_app_customer = var.app_customer
      itse_app_env      = var.app_environment
      itse_app_name     = var.app_name
    }
  }
}
