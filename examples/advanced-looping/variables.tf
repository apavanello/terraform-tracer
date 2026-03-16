variable "region" {
  description = "AWS Region"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
}

variable "public_subnets" {
  description = "Map of public subnets"
  type        = map(string)
}

variable "instance_count" {
  description = "Number of web instances"
  type        = number
  default     = 2
}
