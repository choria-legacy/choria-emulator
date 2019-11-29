resource "aws_security_group" "nats" {
  vpc_id = var.vpc_id
  tags   = var.tags

  ingress {
    from_port   = 4222
    to_port     = 4222
    protocol    = "tcp"
    cidr_blocks = var.emulator_networks
  }

  ingress {
    from_port   = 7422
    to_port     = 7422
    protocol    = "tcp"
    cidr_blocks = var.emulator_networks
  }
}
