resource "aws_security_group" "management" {
  vpc_id = aws_vpc.choria_emulator.id
  tags   = var.tags

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = var.management_networks
  }
}
