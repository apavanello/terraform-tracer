# Security group for the web servers
resource "aws_security_group" "web_sg" {
  name        = "${var.environment}-web-sg"
  description = "Security group for web module"
  vpc_id      = var.vpc_id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# The actual instances deployed across the provided subnets
resource "aws_instance" "web" {
  count = var.instance_count

  # Ami logic
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t3.micro"
  
  # Distributing instances across subnets using element()
  subnet_id = element(var.subnet_ids, count.index)
  vpc_security_group_ids = [aws_security_group.web_sg.id]

  tags = {
    Name        = "${var.environment}-web-${count.index}"
    Environment = var.environment
  }
}

output "security_group_id" {
  value = aws_security_group.web_sg.id
}
