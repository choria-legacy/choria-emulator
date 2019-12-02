# A broker is created on the emulator broker node.
#
# On each emulator by default 500 Choria Servers are run connecting
# to the broker over the Internet.
#
#                             +---------------------------------+
#                             | Measure on Shell in eu central 1|
#                             +----------------+----------------+
#                                              |
#                                +-------------v-------------+
#                                |                           |
#           +-------------------->    NATS in eu-central-1   <-------------------------+
#           |                    |                           |                         |
#           |                    +--------------^------------+                         |
#           |                                   |                                      |
#           |                                   |                                       |
#  +--------+-------+                  +--------+-------+                      +--------+-------+
#  | Choria Servers |                  | Choria Servers |                      | Choria Servers |
#  |    us-west-1   |                  |    us-east-1   |                      |  eu-central-1  |
#  +----------------+                  +----------------+                      +----------------+
plan emulator::scenario::broker_basic (
  Optional[String] $broker = undef,
  Integer $instances = 500,
) {
  $_broker = emulator::data("emulator_broker", $broker)

  choria::run_playbook("emulator::scenario::stop_all", {})

  info("Starting NATS on broker instances")
  choria::run_playbook("emulator::nats::brokers_start", {})

  info("Starting ${instances} Choria Servers in each region")
  choria::run_playbook("emulator::choria::regional_start", {
    "servers" => $_broker,
    "instances" => $instances
  })
}