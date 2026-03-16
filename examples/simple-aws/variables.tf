variable "vpc_cidr" {
  type    = string
  default = "10.0.0.0/16"
}

variable "public_subnet_cidr" {
  type    = string
  default = "10.0.1.0/24"
}

variable "availability_zone" {
  type    = string
  default = "us-east-1a"
}

variable "db_name" {
  type    = string
  default = "mydb"
}

variable "db_instance_class" {
  type    = string
  default = "db.t3.micro"
}

variable "db_storage" {
  type    = string
  default = "20"
}
