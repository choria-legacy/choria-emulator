variable "management_networks" {
  description = "CIDRs from where you will be able to ssh into the network"
}

variable "emulator_subnet_cidr" {
  description = "CIDR for public subnet in availability zone A"
}

variable "vpc_cidr" {
  description = "CIDR for VPC"
}

variable "tags" {
  description = "Resource tags"
  type        = map(string)
}
