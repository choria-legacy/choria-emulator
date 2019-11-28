function emulator::broker_nodes() >> Array[String] {
  choria::discover(
    "test" => true,
    "agents" => ["emulator"],
    "facts" => ["role=broker"],
  )
}
