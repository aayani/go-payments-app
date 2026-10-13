// Default VPC lookup - supplies the network for payments resources - reuses the existing AWS network
data "aws_vpc" "default" {
  cidr_block = "172.31.0.0/16"
}

// Default subnet lookup - places the ALB and ECS tasks in the VPC - lets the ALB span availability zones
data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}
