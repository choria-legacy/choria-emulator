# On each emulator by default 10 Choria Servers are run connecting
# to to the NGS global name
#
#                             +---------------------------------+
#                             | Measure on Shell in eu central 1|
#                             +----------------+----------------+
#                                              |
#                                +-------------v-------------+
#                                |                           |
#           +-------------------->        Synadia NGS        <-------------------------+
#           |                    |                           |                         |
#           |                    +--------------^------------+                         |
#           |                                   |                                      |
#           |                                   |                                       |
#  +--------+-------+                  +--------+-------+                      +--------+-------+
#  | Choria Servers |                  | Choria Servers |                      | Choria Servers |
#  +----------------+                  +----------------+                      +----------------+
plan emulator::scenario::ngs_basic (
  Integer $instances = 10,
) {
  choria::run_playbook("emulator::scenario::stop_all", {})

  info("Starting ${instances} Choria Servers in each region")
  choria::run_playbook("emulator::choria::regional_start", {
    "servers" => "connect.ngs.global:4222",
    "instances" => $instances
  })
}