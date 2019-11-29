class profiles::broker {
  class{"choria::broker":
    network_broker => true
  }
}