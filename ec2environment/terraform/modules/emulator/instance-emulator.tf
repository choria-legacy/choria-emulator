data "template_file" "emulator_init" {
  template = file("cloud-init/common.txt")
  vars = {
    puppet_master_ip = var.puppetmaster_ip
    role             = "emulator"
    puppet_psk       = var.puppet_psk
  }
}

module "emulator" {
  source = "../instance"

  instance_count     = var.emulator_count
  type               = "t2.medium"
  subnet_id          = var.network_id
  security_group_ids = var.security_group_ids
  user_data          = data.template_file.emulator_init.rendered
  tags               = var.tags
}

output "emulators" {
  value = module.emulator.*.public_dns
}
