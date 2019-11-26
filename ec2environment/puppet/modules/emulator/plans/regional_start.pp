plan emulator::regional_start (
  String $servers,
  Boolean $tls=false,
  Optional[String] $credentials=undef,
  Integer $monitor=8080,
  Integer $instances=1,
  Integer $agents=1,
  Integer $collectives=1,
) {
  $regions = emulator::regions()

  if $regions.length == 0 {
    fail("no region data found")
  }

  info(sprintf("Starting %d regions", $regions.length))

  $regions.each |$region| {
    info("Starting ${region}")

    choria::run_playbook("emulator::start",   
      "region" => $region,
      "servers" => $servers,
      "tls" => $tls,
      "credentials" => $credentials,
      "monitor" => $monitor,
      "instances" => $instances,
      "agents" => $agents,
      "collectives" => $collectives,
    )
  }
}
