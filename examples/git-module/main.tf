provider "aws" {
  region = "us-east-1"
}

# This module uses a git source — the tracer should clone and parse it
module "vpc" {
  source = "git::https://github.com/terraform-aws-modules/terraform-aws-vpc.git?ref=v5.1.0"

  name = var.vpc_name
  cidr = var.vpc_cidr

  azs             = var.azs
  public_subnets  = var.public_subnets
  private_subnets = var.private_subnets

  enable_nat_gateway = true
  single_nat_gateway = true
}

# A local resource that depends on the remote module
resource "aws_security_group" "app_sg" {
  name        = "app-sg"
  description = "Application security group"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "app" {
  ami                    = "ami-0c55b159cbfafe1f0"
  instance_type          = "t3.micro"
  subnet_id              = module.vpc.public_subnets[0]
  vpc_security_group_ids = [aws_security_group.app_sg.id]

  tags = {
    Name = "app-instance"
  }
}
