plan emulator::setup::download (
  Optional[String] $emulator_url = undef,
) {
  if $emulator_url {
    choria::data("emulator_url", $emulator_url, emulator::ds())
    $_emulator = $emulator_url
  } else {
    $_emulator = choria::data("emulator_url", emulator::ds())
  }

  $emulators = emulator::all_nodes()

  choria::task("mcollective",
    "nodes" => $emulators,
    "action" => "emulator.download",
    "batch_size" => 10,
    "silent" => true,
    "post" => ["summarize"],
    "properties" => {
      "http" => $_emulator,
      "target" => "choria-emulator"
    }
  ) 

  sprintf("Downloaded choria-emulator from %s to %d nodes", $_emulator, $emulators.length)
}
