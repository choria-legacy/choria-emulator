class profiles::common {
  include mcollective
  include limits

  package{"gcc":
    ensure => "present"
  }
}
