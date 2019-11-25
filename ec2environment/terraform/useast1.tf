provider "aws" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = "us-east-1"
  alias      = "useast1"
}

module "useast1_network" {
  source = "./modules/network"

  management_networks  = var.management_networks
  emulator_subnet_cidr = var.emulator_subnet_cidr
  vpc_cidr             = var.vpc_cidr

  providers = {
    aws = aws.useast1
  }

  tags = {
    Project = "choria_emulator"
  }
}

module "useast1_emulators" {
  source = "./modules/emulator"

  puppetmaster_ip    = module.control.puppetmaster_ip
  network_id         = module.useast1_network.subnet_id
  security_group_ids = [module.useast1_network.internal_security_group_id, module.useast1_network.management_security_group_id]
  puppet_psk         = var.puppet_psk

  providers = {
    aws = aws.useast1
  }

  tags = {
    Project = "choria_emulator"
  }
}

output "useast1" {
  value = module.useast1_emulators.emulators
}
