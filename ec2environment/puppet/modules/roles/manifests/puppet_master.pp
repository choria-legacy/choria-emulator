class roles::puppet_master {
  include profiles::common
  include profiles::nats
  include profiles::webserver
}
