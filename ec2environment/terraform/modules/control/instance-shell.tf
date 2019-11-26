data "template_file" "shell_init" {
  template = file("cloud-init/common.txt")
  vars = {
    puppet_master_ip = module.puppetmaster.private_ips[0]
    role             = "shell"
    puppet_psk       = var.puppet_psk
    region           = data.aws_region.current.name
  }
}

module "shell" {
  source = "../instance"

  type               = "t2.medium"
  subnet_id          = var.network_id
  security_group_ids = var.security_group_ids
  user_data          = data.template_file.shell_init.rendered
  tags               = var.tags
}

output "shell" {
  value = module.shell.*.public_dns
}
