# Background

When considering to build a large scale Choria network, one with more than 10 000 nodes, you have a number of things to consider when it comes to planning the optimal size of networks and network structure wrt Federation.

It is likely that a single flat network will not perform to your needs and this is highly dependant on your workloads.  You have to answer these questions:

  * Total acceptable response time to an estate wide `mco rpc rpcutil ping` call
  * Discovery time that, given optimal conditions, reliably discovers all your nodes
  * Number of agents you will be running
  * Number of sub collectives you will be running
  * Max acceptible downtime you can experience when rebooting the middleware layer

You have to supply the targets for these that fit your workloads and style.

With these targets, how do you figure out what are the right size networks? This repository provided tools to help you answer and validate these items.

The end result is a method to repeatedly and reliably study the performance of Choria under various configurations and topologies.  Armed with this you can determine optimal sizes for your Federated Networks etc. You'll be able to determine if your topology supports each of the target numbers you set above.

This is primarily a research tool the Choria developers use to validate Choria and the various architectures we support.  If we want to see how might 200 Federated Networks behave, it's easy to standup a network with exactly that given a some EC2 instances.

Expect the eventual outcome to be statements about the scalability of Choria given specific test environments.  But we'd like this to be usable and open so others can validate our methods and reproduce for their own networks.

# Choria Emualtor
As it's unrealistic that you'll have a 100 000 nodes just lying around your lab an emulator is included that creates multiple Choria instances in memory in a single process.

Each instance can have a number of emulated agents, belong to many sub collectives and generally you'll be able to interact with them from the normal Choria `mco` CLI in representitive ways and it will generate traffic matching a real Choria.

On my MBP I can run 2 000 instances of Choria along with a Choria Broker and the Choria client all on the same laptop and response times are around 1.5 seconds for all nodes. Of course what is realistic per VM varies on spec but I have had good success with 750-1000 emulated instances on a 2GB VM.

The idea is that you will use many VMs, say a few 100, deploy standard Choria to them along with an agent, that will be provided here, to manage a big network of emulated Choria instances.

Each of these 100s of VMs can run, lets say, a thousand Choria instances at a time and you can point them at different topologies of NATS, Federation etc and do tests with different concurrencies and payload sizes.

Using this is should be feasable to make realistic emulations of Choria networks with 100 000 to 300 000 nodes.

## Instrumentation Tool

A tool is included to make 100s or 1000s of requests against the emulated network, during this it records vast
amounts of statistics and write these out to CSV files for offline analysis.

Using this you can verify if your topology meets the requirements set out above thus helping you determine what is the optimal sizes for Choria Networks.

Scripts or a Google Sheet will be provided to analyze this data.

## Understanding how Choria use NATS.

Generally the cost and performance of a network broker comes down to:

  * Number of TCP Connections
  * TLS or Plain
  * Number of message targets and their types
  * Number of subscribers
  * Cluster overhead

1 single Choria node will:

  * Maintain a single TCP session to the NATS broker
  * Use TLS unconditionally
  * Subscribe to one queue unique to itself in every subcollective
  * For every agent like `puppet` a broadcast topic exist in every sub collective
  * Subscribe to `subcollectives * agents` broadcast queues. The least amount of agents is 2

So 10 agents in 5 sub collectives will use:

   * 50 broadcast target for agents
   * 5 targets for the node for directed traffic
   * 1 TCP Connection

100 nodes will have 5 500 subscription, 550 NATS targets and 100 NATS TCP connections.

So if you intend to run many sub collectives or many agents you need to consider that in your testing as it will impact the performance, bootstrap time and memory consumption of the NATS broker.

## Scenarios

Some setup is required to get this going, see the [Environment Setup](docs/PREPARE.md) guide and complete it before doing any of the scenarios.

A number of possible architectures can be built using this emulator, please see specific docs below:

  * [Flat network](docs/SCENARIO-FLAT.md) with Choria instances and Client sharing a single NATS broker
  * Flat netowrk with Choria instances and Client sharing a cluster of 3 or 5 NATS brokers
  * Federated network with NATS + Federation on every node connecting to central Federation

For each scenario you can monitor your NATS infrastructure for [connection and subscription rates](docs/NATS.md).

For each scenario you can adjust the amount of agents and sub collectives to your needs.
