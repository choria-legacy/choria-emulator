# Creates a NATS Server on the broker node with leaf node
# NATS Servers from each of the emulator.
#
# On each emulator by default 500 Choria Servers are run connecting
# to the leafnode on localhost
#
#                             +---------------------------------+
#                             | Measure on Shell in eu central 1|
#                             +----------------+----------------+
#                                              |
#                                +-------------v-------------+
#                                |                           |
#           +-------------------->         Synadia NGS       <-------------------------+
#           |                    |                           |                         |
#           |                    +--------------^------------+                         |
#           |                                   |                                      |
#           |                                   |                                      |
#           |                                   |                                      |
#           |                                   |                                      |
# +---------+---------+               +---------+---------+                 +----------+-----------+
# |                   |               |                   |                 |                      |
# | Leaf in us-west-1 |               | Leaf in us-east-1 |                 | Leaf in eu-central-1 |
# |                   |               |                   |                 |                      |
# +---------+---------+               +---------+---------+                 +-----------+----------+
#           ^                                   ^                                       ^
#           |                                   |                                       |
#  +--------+-------+                  +--------+-------+                      +--------+-------+
#  | Choria Servers |                  | Choria Servers |                      | Choria Servers |
#  +----------------+                  +----------------+                      +----------------+
plan emulator::scenario::ngs_leafnodes (
  String $credentials,
  Integer $instances = 10
) {
  $_credentials = emulator::data("emulator_credentials", $credentials)

  choria::run_playbook("emulator::scenario::stop_all", {})

  info("Starting NATS servers on emulators with leaf connections to connect.ngs.global:7422")
  choria::run_playbook("emulator::nats::emulators_start", {
    "leafnode" => true,
    "servers" => "nats://connect.ngs.global:7422",
    "credentials" => $_credentials
  })

  choria::run_playbook("emulator::choria::regional_start", {
    "servers" => "localhost:4222",
    "instances" => $instances
  })
}