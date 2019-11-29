# Creates a NATS Server on the broker node with leaf node
# NATS Servers from each of the emulator.
#
# On each emulator by default 500 Choria Servers are run connecting
# to the leafnode on localhost
#
#                +------------------+
#                | Measure on Shell |
#                +---------+--------+
#                          |
#                          |
#                          v
#                +---------+-------+
#     +--------> |  NATS on Broker | <------+
#     |          +-----------------+        |
#     |                                     |
#     |                                     |
#     |                                     |
#     |                                     |
#     |                                     |
# +---+---------------+      +--------------+----+
# | Leaf on emulator1 |      | Leaf on emulator2 |
# +---------+---------+      +----------+--------+
#           ^                           ^
#           |                           |
#     +-----+------+             +------+-----+
#     | 500 Choria |             | 500 Choria |
#     +------------+             +------------+
plan emulator::scenario::broker_leafnodes (
  Optional[String] $broker_public_name,
  Integer $instances = 500
) {
  $_pub_name = emulator::data("emulator_broker_name", $broker_public_name)

  choria::run_playbook("emulator::scenario::stop_all", {})

  info("Starting NATS on broker instances")
  choria::run_playbook("emulator::nats::brokers_start", {})

  info("Starting NATS servers on emulators with leaf connections to ${_pub_name}")
  choria::run_playbook("emulator::nats::emulators_start", {
    "leafnode" => true,
    "servers" => "nats://${_pub_name}:7422"
  })

  choria::run_playbook("emulator::choria::regional_start", {
    "servers" => "localhost:4222",
    "instances" => $instances
  })
}