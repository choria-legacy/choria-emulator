#!/bin/bash

set -e

puppet module install choria-choria
puppet module install puppetlabs-apache
puppet module install puppetlabs-ntp
puppet module install puppetlabs-puppetdb
puppet module install puppetlabs-concat
puppet module install puppetlabs-hocon
puppet module install saz-limits
puppet module install puppetlabs-puppet_authorization --ignore-dependencies
puppet module install camptocamp-puppetserver --ignore-dependencies