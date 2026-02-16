/*
 * App Monitoring Archiver Hashicorp Terraform Role
 */
resource "aws_iam_role" "hcp" {
  description = "for Terraform workspace app-monitoring-archiver"
  name        = "app-monitoring-archiver-hcp-terraform"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Federated = "arn:aws:iam::369020531563:oidc-provider/app.terraform.io"
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringEquals = {
            "app.terraform.io:aud" = "aws.workload.identity"
          }
          StringLike = {
            "app.terraform.io:sub" = "organization:gtis:project:*:workspace:app-monitoring-archiver:run_phase:*"
          }
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "hcp" {
  name = "Deployment"
  role = aws_iam_role.hcp.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "iam:CreateUser",
          "iam:DeleteUser",
          "iam:GetUser",
          "iam:TagUser",
          "iam:UntagUser",
          "iam:ListUserTags",
          "iam:CreateAccessKey",
          "iam:DeleteAccessKey",
          "iam:GetAccessKeyLastUsed",
          "iam:ListAccessKeys",
          "iam:AttachUserPolicy",
          "iam:DetachUserPolicy",
          "iam:ListAttachedUserPolicies",
          "iam:ListUserPolicies",
        ]
        Effect   = "Allow"
        Resource = "arn:aws:iam::369020531563:user/app-monitoring-archiver-*"
        Sid      = "IAMUserManagement"
      },
      {
        Action = [
          "iam:CreatePolicy",
          "iam:DeletePolicy",
          "iam:GetPolicy",
          "iam:GetPolicyVersion",
          "iam:ListPolicyVersions",
          "iam:CreatePolicyVersion",
          "iam:DeletePolicyVersion",
        ]
        Effect   = "Allow"
        Resource = "arn:aws:iam::369020531563:policy/app-monitoring-archiver-*"
        Sid      = "IAMPolicyManagement"
      },
      {
        Action = [
          "iam:GetRole",
          "iam:GetRolePolicy",
          "iam:ListRolePolicies",
          "iam:ListAttachedRolePolicies",
        ]
        Effect   = "Allow"
        Resource = "arn:aws:iam::369020531563:role/app-monitoring-archiver-hcp-terraform"
        Sid      = "RoleSelfRead"
      },
    ]
  })
}
