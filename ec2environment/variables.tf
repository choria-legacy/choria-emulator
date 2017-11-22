variable "access_key" {
  description = "AWS access key"
}

variable "secret_key" {
  description = "AWS secret access key"
}

variable "shell_count" {
  description = "Amount of shell instances"
  default = 1
}

variable "nats_count" {
  description = "Amount of NATS instances"
  default = 1
}

variable "emulator_count" {
  description = "Amount of Emulator Instances"
  default = 16
}

variable "region" {
  description = "AWS region to host your networks"
  default     = "eu-central-1"
}

variable "avail_zone" {
  description = "AWS availability zone to host your network"
  default     = "eu-central-1a"
}

variable "vpc_cidr" {
  description = "CIDR for VPC"
  default     = "10.128.0.0/16"
}

variable "emulator_subnet_cidr" {
  description = "CIDR for public subnet in availability zone A"
  default     = "10.128.1.0/24"
}

/* centos 7 with updates in eu-central-1 */
variable "amis" {
  description = "Base AMI to launch the instances with"
  default = {
    eu-central-1 = "ami-fa2df395"
  }
}
