class roles::puppet_master {
  include profiles::common
  include profiles::broker
  include profiles::webserver
}
