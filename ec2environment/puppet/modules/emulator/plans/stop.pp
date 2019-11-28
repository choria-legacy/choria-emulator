plan emulator::stop (
  Optional[String] $region=undef,
) {
  if $region {
    $nodes = emulator::regional_emulator_nodes($region)
  } else {
    $nodes = emulator::emulator_nodes()
  }

  $results = choria::task("mcollective",
    "nodes" => $nodes,
    "action" => "emulator.stop",
    "silent" => true,
    "properties" => {}
  )

  $results.each |$result| {
    if $result.ok {
      info(sprintf("%s: %s: stopped: %s", $result["sender"], $result["statusmsg"], !$result["data"]["status"]))
    } else {
      error(sprintf("%s: %s", $result["sender"], $result["statusmsg"]))
    }
  }

  undef
}
