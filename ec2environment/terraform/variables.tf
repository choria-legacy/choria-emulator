variable "access_key" {
  description = "AWS access key"
}

variable "secret_key" {
  description = "AWS secret access key"
}

variable "region" {
  description = "AWS region to host your networks"
  default     = "eu-central-1"
}

variable "broker_count" {
  description = "Amount of Broker instances"
  default     = 1
}

variable "emulator_count" {
  description = "Amount of Emulator Instances"
  default     = 1
}

variable "vpc_cidr" {
  description = "CIDR for VPC"
  default     = "10.128.0.0/16"
}

variable "emulator_subnet_cidr" {
  description = "CIDR for public subnet in availability zone A"
  default     = "10.128.1.0/24"
}

variable "management_networks" {
  description = "CIDRs from where you will be able to ssh into the network"
  default     = ["139.162.163.118/32", "84.255.40.82/32"]
}

variable "puppet_psk" {
  description = "A PSK to bake into puppet certs and to autosign with, should be unique to you and not shared"
}
