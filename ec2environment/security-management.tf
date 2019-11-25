resource "aws_security_group" "management" {
  vpc_id = aws_vpc.choria_emulator.id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["139.162.163.118/32", "84.255.40.82/32"]
  }

  tags = {
    Project = "choria_test"
  }
}
