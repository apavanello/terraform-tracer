environment    = "prod"
region         = "us-east-1"
vpc_cidr       = "10.0.0.0/16"
instance_count = 3

public_subnets = {
  "us-east-1a" = "10.0.1.0/24"
  "us-east-1b" = "10.0.2.0/24"
  "us-east-1c" = "10.0.3.0/24"
}
