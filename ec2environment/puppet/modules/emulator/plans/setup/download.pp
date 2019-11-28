plan emulator::setup::download (
  Optional[String] $emulator_url = undef,
  Optional[String] $nats_url = undef,
) {
  $_emu_url = emulator::data("emulator_url", $emulator_url)
  $_nats_url = emulator::data("nats_url", $nats_url)
  $emulators = emulator::all_nodes()

  info("Downloading choria-emulator from ${_emu_url}")
  choria::task("mcollective",
    "nodes" => $emulators,
    "action" => "emulator.download",
    "batch_size" => 10,
    "silent" => true,
    "properties" => {
      "http" => $_emu_url,
      "target" => "choria-emulator"
    }
  ) 

  info("Downloading nats-server from ${_nats_url}")
  choria::task("mcollective",
    "nodes" => $emulators,
    "action" => "emulator.download",
    "batch_size" => 10,
    "silent" => true,
    "properties" => {
      "http" => $_nats_url,
      "target" => "nats-server"
    }
  ) 

  undef
}
