provider "aws" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = "eu-central-1"
  alias      = "eucentral1"
}

module "network" {
  source = "./modules/network"

  management_networks  = var.management_networks
  emulator_subnet_cidr = var.emulator_subnet_cidr
  vpc_cidr             = var.vpc_cidr

  providers = {
    aws = aws.eucentral1
  }

  tags = {
    Project = "choria_emulator"
  }
}

module "control" {
  source             = "./modules/control"
  network_id         = module.network.subnet_id
  security_group_ids = [module.network.internal_security_group_id, module.network.management_security_group_id]
  puppet_psk         = var.puppet_psk
  emulator_networks  = var.emulator_networks
  vpc_id             = module.network.vpc_id

  providers = {
    aws = aws.eucentral1
  }

  tags = {
    Project = "choria_emulator"
  }
}

module "emulators" {
  source = "./modules/emulator"

  puppetmaster_ip    = module.control.puppetmaster_ip
  network_id         = module.network.subnet_id
  security_group_ids = [module.network.internal_security_group_id, module.network.management_security_group_id]
  puppet_psk         = var.puppet_psk

  providers = {
    aws = aws.eucentral1
  }

  tags = {
    Project = "choria_emulator"
  }
}

output "puppetmaster" {
  value = module.control.puppetmaster_dns
}

output "shell" {
  value = module.control.shell
}

output "emulator" {
  value = module.emulators.emulators
}

output "brokers" {
  value = module.control.brokers
}
