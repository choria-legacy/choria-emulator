plan emulator::choria::regional_start (
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

  $_servers = emulator::data("emulator_servers", $servers)
  
  if $credentials {
    $_creds = emulator::data("emulator_credentials", $credentials)
  } else {
    $_creds = undef
  }

  info(sprintf("Starting %d regions", $regions.length))

  $regions.each |$region| {
    info("Starting ${region}")

    choria::run_playbook("emulator::start",   
      "region" => $region,
      "servers" => $_servers,
      "tls" => $tls,
      "credentials" => $_creds,
      "monitor" => $monitor,
      "instances" => $instances,
      "agents" => $agents,
      "collectives" => $collectives,
    )
  }
}
