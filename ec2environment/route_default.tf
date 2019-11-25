resource "aws_route_table" "default" {
  vpc_id = aws_vpc.choria_emulator.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gateway.id
  }
  tags {
    Project = "choria_emulator"
  }
}
