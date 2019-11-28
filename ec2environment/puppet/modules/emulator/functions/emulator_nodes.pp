function emulator::emulator_nodes() >> Array[String] {
  choria::discover(
    "test" => true,
    "agents" => ["emulator"],
    "facts" => ["role=emulator"]
  )
}
