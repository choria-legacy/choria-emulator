function emulator::all_nodes() >> Array[String] {
  choria::discover(
    "test" => true,
    "agents" => ["emulator"],
  )
}
