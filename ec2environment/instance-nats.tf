resource "template_file" "nats_init" {
  template = "${file("cloud-init/common.txt")}"
  vars {
    puppet_master_ip = "${aws_instance.puppetmaster.private_ip}"
    role = "nats"
  }
}

resource "aws_instance" "nats" {
  count = "${var.nats_count}"
  ami = "${lookup(var.amis, var.region)}"
  instance_type = "t2.medium"
  subnet_id = "${aws_subnet.choria_emulator.id}"
  vpc_security_group_ids = ["${aws_security_group.internal.id}", "${aws_security_group.management.id}"]
  source_dest_check = false
  user_data = "${template_file.nats_init.rendered}"
  root_block_device {
    volume_type = "standard"
    volume_size = 8
    delete_on_termination = true
  }
  tags = {
    Project = "choria_emulator"
  }
}

output "nats" {
  value = "${aws_instance.nats.*.public_dns}"
}
