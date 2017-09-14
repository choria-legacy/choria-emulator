class profiles::emulator (
  Boolean $client = false,
  Boolean $server = true
) {
  class{"mcollective_agent_emulator":
    client => $client,
    server => $server
  }
}
