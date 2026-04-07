# Role for Continuous Deployment using CDK

resource "aws_iam_role" "cd" {
  description = "for GitHub Actions to deploy ${var.github_repository}"
  name        = "${var.app_name}-cd"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Sid    = "GitHub"
      Effect = "Allow"
      Action = "sts:AssumeRoleWithWebIdentity"
      Principal = {
        Federated = var.github_oidc_provider_arn
      }
      Condition = {
        StringEquals = {
          "token.actions.githubusercontent.com:aud" : "sts.amazonaws.com"
        },
        StringLike = {
          "token.actions.githubusercontent.com:sub" : "repo:${var.github_repository}:*"
        }
      }
    }]
  })
}

resource "aws_iam_role_policy" "cd" {
  name = "${var.app_name}-cd"
  role = aws_iam_role.cd.name

  policy = jsonencode({
    "Version" : "2012-10-17"
    "Statement" : [{
      Effect   = "Allow"
      Action   = "sts:AssumeRole"
      Resource = "arn:aws:iam::*:role/cdk-*"
    }]
  })
}
