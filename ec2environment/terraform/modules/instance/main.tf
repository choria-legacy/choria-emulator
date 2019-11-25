data "aws_region" "current" {}

resource "aws_instance" "self" {
  count                  = var.instance_count
  ami                    = lookup(var.amis, data.aws_region.current.name)
  instance_type          = var.type
  subnet_id              = var.subnet_id
  vpc_security_group_ids = var.security_group_ids
  source_dest_check      = false
  user_data              = var.user_data
  tags                   = var.tags
  root_block_device {
    volume_type           = "standard"
    volume_size           = 8
    delete_on_termination = true
  }
}

output "public_dns" {
  value = aws_instance.self.*.public_dns
}

output "public_ips" {
  value = aws_instance.self.*.public_ip
}

output "private_ips" {
  value = aws_instance.self.*.private_ip
}
