variable "aws_endpoint" {
  type = string
}

variable "app_image" {
  type        = string
  default     = "go-payments-app:local"
  description = "Container image for the payments ECS service."
}