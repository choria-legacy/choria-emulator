package emulator

import (
	"fmt"
	"sync"

	"github.com/choria-io/go-choria/protocol"

	"github.com/choria-io/go-choria/mcollective"
	log "github.com/sirupsen/logrus"
)

type ChoriaEmulationInstance struct {
	name      string
	choria    *mcollective.Choria
	connector mcollective.Connector
	agents    map[string]Agent
}

type Agent interface {
	Init() error
	Name() string
	HandleAgentMsg(msg string) (*[]byte, error)
}

func NewInstance(choria *mcollective.Choria) (i *ChoriaEmulationInstance, err error) {
	i = &ChoriaEmulationInstance{
		name:   choria.Config.Identity,
		choria: choria,
	}

	logger := log.WithFields(log.Fields{"emulator": i.name})

	servers := func() ([]mcollective.Server, error) {
		return i.choria.MiddlewareServers()
	}

	_, err = i.choria.MiddlewareServers()
	if err != nil {
		return nil, fmt.Errorf("Could not find initial NATS servers: %s", err.Error())
	}

	i.connector, err = choria.NewConnector(servers, choria.Certname(), logger)
	if err != nil {
		return nil, fmt.Errorf("Could not create connector for instance %s: %s", i.name, err.Error())
	}

	i.agents = make(map[string]Agent)

	return
}

func (self *ChoriaEmulationInstance) Init(agentCount int) error {
	log.Infof("Starting emulator instance %s with %d emulated agents in %d collectives", self.name, agentCount, len(self.choria.Config.Collectives))

	discovery := &DiscoveryAgent{}
	err := discovery.Init()
	if err != nil {
		return fmt.Errorf("Could not initialize discovery agent: %s", err.Error())
	}

	self.agents["discovery"] = discovery

	for i := 0; i <= agentCount-1; i++ {
		name := fmt.Sprintf("emulated%d", i)
		emulated := &EmulatedAgent{name: name}
		err = emulated.Init()
		if err != nil {
			return fmt.Errorf("Could not initialize emulated agent %d: %s", i, err.Error())
		}

		self.agents[name] = emulated
	}

	self.subscribeNode()
	self.subscribeAgents()

	return nil
}

func (self *ChoriaEmulationInstance) ProcessRequests(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		rawmsg := self.connector.Receive()

		// copy from choria.NewMessageFromTransportJSON with security turned off
		transport, err := self.choria.NewTransportFromJSON(string(rawmsg.Data))
		if err != nil {
			log.Warn(err.Error())
			continue
		}

		srequest, err := self.choria.NewSecureRequestFromTransport(transport, true)
		if err != nil {
			log.Warn(err.Error())
			continue
		}

		request, err := self.choria.NewRequestFromSecureRequest(srequest)
		if err != nil {
			log.Warn(err.Error())
			continue
		}

		if filter, ok := request.Filter(); ok {
			process := true

			if len(filter.AgentFilters()) > 0 {
				for _, agent := range filter.AgentFilters() {
					// mco find wants this but it will fail for actual requests
					if agent == "rpcutil" {
						continue
					}

					if _, ok := self.agents[agent]; !ok {
						log.Warnf("Ignoring message due to agent filter: %s", agent)
						process = false
						break
					}
				}
			}

			if !process {
				continue
			}
		}

		protocol.CopyFederationData(transport, request)

		msg, err := mcollective.NewMessageFromRequest(request, transport.ReplyTo(), self.choria)
		if err != nil {
			log.Warn(err.Error())
			continue
		}

		go self.dispatch(msg, request)
	}
}

func (self *ChoriaEmulationInstance) subscribeNode() {
	for _, collective := range self.choria.Config.Collectives {
		target := fmt.Sprintf("%s.node.%s", collective, self.name)
		self.connector.Subscribe("node", target, "")
		// log.Debugf("Subscribed to %s", target)
	}
}

func (self *ChoriaEmulationInstance) subscribeAgents() {
	for _, collective := range self.choria.Config.Collectives {
		for _, agent := range self.agents {
			target := fmt.Sprintf("%s.broadcast.agent.%s", collective, agent.Name())
			self.connector.Subscribe(fmt.Sprintf("agent.%s", agent), target, "")
			log.Debugf("Subscribed to %s", target)
		}
	}
}

func (self *ChoriaEmulationInstance) dispatch(msg *mcollective.Message, request protocol.Request) {
	agent, ok := self.agents[msg.Agent]
	if !ok {
		log.Warnf("Received a message for an unknown agent %s", msg.Agent)
		return
	}

	rawreply, err := agent.HandleAgentMsg(msg.Payload)
	if err != nil {
		log.Warnf("Handling %s failed: %s", msg.String(), err.Error())
		return
	}

	reply, err := mcollective.NewMessage(string(*rawreply), msg.Agent, msg.Collective(), "reply", msg, self.choria)
	if err != nil {
		log.Warnf("Could not create Reply Message: %s", err.Error())
		return
	}

	transport, err := self.choria.NewTransportFromMessage(reply, request)
	if err != nil {
		log.Warnf("Could not create Reply Transport: %s", err.Error())
		return
	}

	protocol.CopyFederationData(request, transport)

	j, err := transport.JSON()
	if err != nil {
		log.Warnf("Could not extract JSON data from transport: %s", err.Error())
		return
	}

	err = self.connector.PublishRaw(msg.ReplyTo(), []byte(j))
	if err != nil {
		log.Warnf("Sending reply from %s failed: %s", msg.String(), err.Error())
		return
	}
}
