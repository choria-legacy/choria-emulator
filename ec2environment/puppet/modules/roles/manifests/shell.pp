class roles::shell {
  include profiles::common

  class{"profiles::emulator":
    client => true,
    server => false
  }
}
