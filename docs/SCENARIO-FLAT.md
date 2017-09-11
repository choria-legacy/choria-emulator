# Flat Network with Single Broker

This is the most basic flat network, a single NATS broker with 1000s of Choria instances connected to it.

Use this to determine the top scale a certain spec NATS broker will handle unclustered, this is a good setup for a single Federation Member network, though you might want to test HA too.

![flat network](scenario_flat.png)

## Setup

### NATS

You need `gnatsd` running on one of the NATS servers where all the nodes can connect to, a simple TLS free `gnatsd` can be started with:

```
# ulimit -u unlimited
# gnatsd -T -m 8222 -p 4222
```

Lets say for example this is running on `192.168.1.1:4222`

### Choria Instances

We'll run 100 instances per node, 10 agents and 5 sub collectives:

```
$ mco playbook run start-emulator.yaml --agents 10 --collectives 5 --instances 100 --servers 192.168.1.1:4222
```

After your tests you can stop it using:

```
$ mco playbook run stop-emulator.yaml
```

### Client Config

You'll need a special client cfg to connect to this network, save it in `emulator-client.cfg`:

```
plugin.choria.middleware_hosts = 192.168.1.1:4222

# for very large networks you'll need to do some tweaks here
discovery_timeout = 5

collectives = mcollective
connection_timeout = 3
connector = nats
identity = emulator_client
libdir = /opt/puppetlabs/mcollective/plugins
logger_type = console
loglevel = warn
main_collective = mcollective
securityprovider = choria
plugin.choria.security.serializer = json
plugin.choria.use_srv_records = false
```

Run the scenario using `measure-collective.rb --config emulator-client.cfg` with additional options for amount of tests etc