provider "aws" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = "us-west-1"
  alias      = "uswest1"
}

module "uswest1_network" {
  source = "./modules/network"

  management_networks  = var.management_networks
  emulator_subnet_cidr = var.emulator_subnet_cidr
  vpc_cidr             = var.vpc_cidr

  providers = {
    aws = aws.uswest1
  }

  tags = {
    Project = "choria_emulator"
  }
}

module "uswest1_emulators" {
  source = "./modules/emulator"

  puppetmaster_ip    = module.control.puppetmaster_ip
  network_id         = module.uswest1_network.subnet_id
  security_group_ids = [module.uswest1_network.internal_security_group_id, module.uswest1_network.management_security_group_id]
  puppet_psk         = var.puppet_psk
  emulator_count     = 2

  providers = {
    aws = aws.uswest1
  }

  tags = {
    Project = "choria_emulator"
  }
}

output "uswest1" {
  value = module.uswest1_emulators.emulators
}
