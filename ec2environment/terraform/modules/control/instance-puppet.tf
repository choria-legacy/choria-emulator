data "template_file" "puppetmaster_init" {
  template = file("cloud-init/puppet-master.txt")
  vars = {
    puppet_psk = var.puppet_psk
  }
}

module "puppetmaster" {
  source = "../instance"

  subnet_id          = var.network_id
  security_group_ids = var.security_group_ids
  user_data          = data.template_file.puppetmaster_init.rendered
  tags               = var.tags
}

output "puppetmaster_dns" {
  value = module.puppetmaster.public_dns[0]
}

output "puppetmaster_ip" {
  value = module.puppetmaster.public_ips[0]
}

