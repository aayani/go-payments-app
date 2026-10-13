provider "aws" {
  region                      = "us-east-1"
  access_key                  = "test"
  secret_key                  = "test"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    ec2         = var.aws_endpoint
    ecs         = var.aws_endpoint
    elasticache = var.aws_endpoint
    elbv2       = var.aws_endpoint
    iam         = var.aws_endpoint
    logs        = var.aws_endpoint
    rds         = var.aws_endpoint
  }
}
