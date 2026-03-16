provider "aws" {
  region = var.region
}

# 1. Defining a core VPC
resource "aws_vpc" "core" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name        = "${var.environment}-core-vpc"
    Environment = var.environment
  }
}

# 2. Creating multiple subnets using for_each and a map variable
resource "aws_subnet" "public" {
  for_each = var.public_subnets

  vpc_id                  = aws_vpc.core.id
  cidr_block              = each.value
  availability_zone       = each.key
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.environment}-public-${each.key}"
    Type = "Public"
  }
}

# 3. Generating a list of subnet IDs purely for dependency testing
locals {
  public_subnet_ids = [for s in aws_subnet.public : s.id]
}

# 4. Using a module to deploy an application in those subnets
module "web_app" {
  source = "./modules/web_server"

  environment    = var.environment
  vpc_id         = aws_vpc.core.id
  # References the subnets created by the for_each loop above
  subnet_ids     = local.public_subnet_ids
  instance_count = var.instance_count
}

# 5. Iterating with count to create security group rules based on var.instance_count
resource "aws_security_group_rule" "app_ingress" {
  count = var.instance_count

  type              = "ingress"
  from_port         = 8080 + count.index
  to_port           = 8080 + count.index
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  # The module outputs its security group ID
  security_group_id = module.web_app.security_group_id
}
