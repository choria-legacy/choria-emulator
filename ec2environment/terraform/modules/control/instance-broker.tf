data "template_file" "broker_init" {
  template = file("cloud-init/common.txt")
  vars = {
    puppet_master_ip = module.puppetmaster.private_ips[0]
    role             = "broker"
    puppet_psk       = var.puppet_psk
    region           = data.aws_region.current.name
  }
}

module "broker" {
  source = "../instance"

  type               = "t2.medium"
  subnet_id          = var.network_id
  security_group_ids = concat(var.security_group_ids, [aws_security_group.nats.id])
  user_data          = data.template_file.broker_init.rendered
  tags               = var.tags
}

output "brokers" {
  value = module.broker.*.public_dns
}
