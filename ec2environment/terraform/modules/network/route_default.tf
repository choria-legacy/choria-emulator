resource "aws_route_table" "default" {
  vpc_id = aws_vpc.choria_emulator.id
  tags   = var.tags

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gateway.id
  }
}
