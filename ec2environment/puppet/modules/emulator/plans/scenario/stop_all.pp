plan emulator::scenario::stop_all() {
  info("Stopping all emulated choria instances")
  choria::run_playbook("emulator::choria::stop", {})

  info("Stopping all nats-server instances")
  choria::run_playbook("emulator::nats::stop", {})
}