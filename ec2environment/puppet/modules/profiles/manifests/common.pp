class profiles::common {
  file{
    "/etc/systemd/system/choria-server.service.d":
      ensure => directory,
      owner => root;

    "/etc/systemd/system/choria-server.service.d/override.conf":
      notify => [Exec["systemctl daemon-reload"], Class["choria::service"]],
      content => "[Service]\nLimitNOFILE=infinity";
  }

  exec{"systemctl daemon-reload":
    refreshonly => true,
    before => Class["choria::service"]
  }

  include mcollective
  include limits

  package{"gcc":
    ensure => "present"
  }
}
