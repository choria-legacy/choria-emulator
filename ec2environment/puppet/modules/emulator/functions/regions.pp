function emulator::regions() {
  $nodes = emulator::all_nodes()

  $regions = choria::task("mcollective",
    "nodes" => $nodes,
    "action" => "rpcutil.get_fact",
    "silent" => true,
    "properties" => {"fact" => "region"}
  ) 

  $regions.ok_set.map |$r| {
    $r["data"]["value"]
  }
}
