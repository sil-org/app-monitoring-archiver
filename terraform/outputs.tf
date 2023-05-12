
output "serverless-access-key-id" {
  value = module.serverless-user.aws_access_key_id
}
output "serverless-secret-access-key" {
  value = nonsensitive(module.serverless-user.aws_secret_access_key)
}
