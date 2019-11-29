plan emulator::nats::stop () {
) {
  $nodes = emulator::all_nodes()

  $results = choria::task("mcollective",
    "nodes" => $nodes,
    "action" => "emulator.stop_nats",
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

}