function emulator::regional_nodes(
  String $region
) >> Array[String] {
  choria::discover(
    "test" => true,
    "agents" => ["emulator"],
    "facts" => ["region=${region}"]
  )
}
