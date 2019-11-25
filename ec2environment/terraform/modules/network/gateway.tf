resource "aws_internet_gateway" "gateway" {
  vpc_id = aws_vpc.choria_emulator.id
  tags   = var.tags
}
