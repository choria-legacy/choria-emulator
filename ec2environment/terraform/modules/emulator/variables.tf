
variable "puppetmaster_ip" {
  description = "IP to use for the puppet master"
  type        = string
}

variable "emulator_count" {
  description = "How many emulators to start"
  default     = 1
}

variable "tags" {
  description = "Tags to apply to the resources"
}

variable "network_id" {
  description = "ID of the network to place the instances in"
}

variable "security_group_ids" {
  description = "IDs of the security groups to join the machines to"
  type        = list(string)
}

variable "puppet_psk" {
  description = "A PSK to bake into puppet certs and to autosign with, should be unique to you and not shared"
}
