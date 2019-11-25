output "vpc_id" {
  value = aws_vpc.choria_emulator.id
}

output "subnet_id" {
  value = aws_subnet.choria_emulator.id
}

output "internal_security_group_id" {
  value = aws_security_group.internal.id
}

output "management_security_group_id" {
  value = aws_security_group.management.id
}
