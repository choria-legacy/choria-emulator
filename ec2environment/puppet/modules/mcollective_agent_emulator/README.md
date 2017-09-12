# mcollective_agent_emulator version 0.0.1

#### Table of Contents

1. [Overview](#overview)
1. [Usage](#usage)
1. [Configuration](#configuration)

## Overview

Choria Agent emulated by choria-emulator

The mcollective_agent_emulator module is generated automatically, based on the source from http://choria.io.

Available Actions:

  * **generate** - Generates random data of a given size

## Usage

You can include this module into your infrastructure as any other module, but as it's designed to work with the [choria mcollective](http://forge.puppet.com/choria/mcollective) module you can configure it via Hiera:

```yaml
mcollective::plugin_classes:
  - mcollective_agent_emulator
```

## Configuration

Server and Client configuration can be added via Hiera and managed through tiers in your site Hiera, they will be merged with any included in this module

```yaml
mcollective_agent_emulator::config:
   example: value
```

This will be added to both the `client.cfg` and `server.cfg`, you can likewise configure server and client specific settings using `mcollective_agent_emulator::client_config` and `mcollective_agent_emulator::server_config`.

These settings will be added to the `/etc/puppetlabs/mcollective/plugin.d/` directory in individual files.

For a full list of possible configuration settings see the module [source repository documentation](http://choria.io).

## Data Reference

  * `mcollective_agent_emulator::gem_dependencies` - Deep Merged Hash of gem name and version this module depends on
  * `mcollective_agent_emulator::manage_gem_dependencies` - disable managing of gem dependencies
  * `mcollective_agent_emulator::package_dependencies` - Deep Merged Hash of package name and version this module depends on
  * `mcollective_agent_emulator::manage_package_dependencies` - disable managing of packages dependencies
  * `mcollective_agent_emulator::class_dependencies` - Array of classes to include when installing this module
  * `mcollective_agent_emulator::package_dependencies` - disable managing of class dependencies
  * `mcollective_agent_emulator::config` - Deep Merged Hash of common config items for this module
  * `mcollective_agent_emulator::server_config` - Deep Merged Hash of config items specific to managed nodes
  * `mcollective_agent_emulator::client_config` - Deep Merged Hash of config items specific to client nodes
  * `mcollective_agent_emulator::policy_default` - `allow` or `deny`
  * `mcollective_agent_emulator::policies` - List of `actionpolicy` policies to deploy with an agent
  * `mcollective_agent_emulator::client` - installs client files when true - defaults to `$mcollective::client`
  * `mcollective_agent_emulator::server` - installs server files when true - defaults to `$mcollective::server`
  * `mcollective_agent_emulator::ensure` - `present` or `absent`

## Development:

To contribute to this MCollective plugin please visit http://choria.io.

This module was generated using the Choria Plugin Packager based on templates found at the [GitHub Project](https://github.com/choria-io/).
