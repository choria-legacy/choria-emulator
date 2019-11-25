resource "aws_vpc" "choria_emulator" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true
  tags                 = var.tags
}
