plan emulator::nats::brokers_start (
  Boolean $leafnode=false,
  Optional[String] $servers=undef,
  Optional[String] $credentials=undef,
  Integer $monitor=8222,
  Integer $clients=4222,
) {
 $_nodes = emulator::broker_nodes()

 choria::run_playbook("emulator::nats::start",
   "nodes" => $_nodes,
   "leafnode" => $leafnode,
   "servers" => $servers,
   "credentials" => $credentials,
   "monitor" => $monitor,
   "clients" => $clients
 )
}