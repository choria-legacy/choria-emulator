# Environment Setup

## Requirements

You'll need the following:

  * A `Puppet Server`
  * 1 `Shell` node to run the tests from and gather stats
  * A number of `Emulator Nodes` to run the emulator on
  * 1, 3 or 5 nodes for `NATS` depending on your scenarios, these should not also run the emulator
  * Copies of `nats-server` and `choria-emulator` hosted on a reachable webserver
  * All the nodes set up with latest Choria and communicating with the Puppet Master over a dedicated NATS running on the Puppet Master (example, just not on the test nodes
  * The emulator nodes must have the `mcollective_agent_emulator` module deployed on them
  * The shell node need the `mcollective_agent_emulator` but only it's client

![overview](overview.png)

## Prepare Emulators

The emulators will need copies of `nats-server` and `choria-emulator` deployed on them, you can do this as follows:

```
% touch ~/.choria-emulator.yaml
% mco playbook run setup-prereqs.yaml \
   --emulator_url=https://shell.internal/choria-emulator-0.0.1 \
   --gnatsd_url=https://shell.internal/gnatsd \
   --choria_url=https://shell.internal/go-choria
```

The answers you supply here will be stored in `~/.choria-emulator.yaml` so next time you can leave off the arguments.

## CLI tools for non TLS operation

If the scenario is one run without TLS then you'll have to do some tweaks to your `mco` binary if you wanted to interact with the network using `mco ping` etc.

```
$ cp /opt/puppetlabs/puppet/bin/mco scripts/mco
```

Edit this and insert below the `#!` the following `$choria_unsafe_disable_nats_tls = true`. Use this `mco` with the config file the scenario gives you to interact with the network.

## Run the scenario

Individual scenario documentation will show you how to deploy and configure the emulator and test client. Use the basic flat scenario to confirm it all works using a 10 count test.

The individual scenarios will show example agents, collectives and instance counts, adjust it to your needs.  Generally I find a 4GB instance can manage 1 000 choria instances well enough.

After the run is completed you can stop the various parts using these commands:

```
$ mco playbook run stop-emulator.yaml
$ mco playbook run stop-nats.yaml
$ mco playbook run stop-federation.yaml
```

## Interpret the results

The test results are stored in directories like `reports/20170911-100024/` you can graph the stats using (TODO URL FOR A SHEET)