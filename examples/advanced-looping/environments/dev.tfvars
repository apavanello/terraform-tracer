environment    = "dev"
region         = "sa-east-1"
vpc_cidr       = "10.10.0.0/16"
instance_count = 1

public_subnets = {
  "sa-east-1a" = "10.10.1.0/24"
  "sa-east-1b" = "10.10.2.0/24"
}
