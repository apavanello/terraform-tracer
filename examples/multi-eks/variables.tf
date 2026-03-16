variable "project_name" {
  type    = string
  default = "acme"
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "vpc_cidr" {
  type    = string
  default = "10.0.0.0/16"
}

variable "public_subnet_a_cidr" {
  type    = string
  default = "10.0.1.0/24"
}

variable "public_subnet_b_cidr" {
  type    = string
  default = "10.0.2.0/24"
}

variable "private_subnet_a_cidr" {
  type    = string
  default = "10.0.10.0/24"
}

variable "private_subnet_b_cidr" {
  type    = string
  default = "10.0.11.0/24"
}

variable "eks_version" {
  type    = string
  default = "1.29"
}

# Platform cluster sizing
variable "platform_node_instance_type" {
  type    = string
  default = "t3.medium"
}
variable "platform_node_desired" {
  type    = string
  default = "2"
}
variable "platform_node_max" {
  type    = string
  default = "5"
}
variable "platform_node_min" {
  type    = string
  default = "1"
}

# Workloads cluster sizing
variable "workloads_node_instance_type" {
  type    = string
  default = "t3.large"
}
variable "workloads_node_desired" {
  type    = string
  default = "3"
}
variable "workloads_node_max" {
  type    = string
  default = "10"
}
variable "workloads_node_min" {
  type    = string
  default = "2"
}

# Database
variable "db_instance_class" {
  type    = string
  default = "db.t3.micro"
}
variable "db_storage" {
  type    = string
  default = "20"
}

# Redis
variable "redis_node_type" {
  type    = string
  default = "cache.t3.micro"
}
variable "redis_num_nodes" {
  type    = string
  default = "2"
}
