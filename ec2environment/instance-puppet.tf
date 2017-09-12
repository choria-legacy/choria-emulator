resource "aws_instance" "puppetmaster" {
  count = 1
  ami = "${lookup(var.amis, var.region)}"
  instance_type = "t2.large"
  subnet_id = "${aws_subnet.choria_emulator.id}"
  vpc_security_group_ids = ["${aws_security_group.internal.id}", "${aws_security_group.management.id}"]
  source_dest_check = false
  user_data = "${file("cloud-init/puppet-master.txt")}"
  root_block_device {
    volume_type = "standard"
    volume_size = 8
    delete_on_termination = true
  }
  tags = {
    Project = "choria_emulator"
  }
}

output "puppetmaster" {
  value = "${aws_instance.puppetmaster.public_dns}"
}

