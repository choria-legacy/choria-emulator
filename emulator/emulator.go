package emulator

import (
	"context"
	"fmt"
	"sync"

	"github.com/choria-io/go-protocol/protocol"
	"github.com/sirupsen/logrus"

	"github.com/choria-io/go-config"

	"github.com/choria-io/go-choria/choria"
	"github.com/choria-io/go-choria/server"
	gorpc "github.com/choria-io/mcorpc-agent-provider/mcorpc/golang"
)

var (
	instanceCount   int
	agentCount      int
	collectiveCount int
	statusPort      int
	name            string
	configFile      string
	enableTLS       bool
	enableTLSVerify bool
	protocolSecure  bool
	brokers         []string
	credentials     string

	err       error
	ctx       context.Context
	cancel    func()
	wg        *sync.WaitGroup
	fw        *choria.Framework
	instances []*server.Instance
	log       *logrus.Entry
	mu        *sync.Mutex
)

func NewEmulator() (emulated []*server.Instance, err error) {
	server.RegisterAdditionalAgentProvider(&gorpc.Provider{})

	if protocolSecure {
		log.Infof("Enabling Choria protocol security")
		protocol.Secure = "true"
	} else {
		log.Infof("Disabling Choria protocol security")
		protocol.Secure = "false"
	}

	log.Infof("Starting %d Choria Server instances each belonging to %d collective(s) with %d emulated agent(s)", instanceCount, collectiveCount, agentCount)

	mu = &sync.Mutex{}
	instanceID, _ := choria.NewRequestID()

	for i := 1; i <= instanceCount; i++ {
		name := fmt.Sprintf("%s-%s-%d", name, instanceID[1:6], i)
		log.Infof("Creating instance %s", name)

		wg.Add(1)
		startf := func() {
			defer wg.Done()
			srv, err := newInstance(name)
			if err != nil {
				log.Errorf("Could not start instance %d: %s", i, err)
				return
			}

			mu.Lock()
			emulated = append(emulated, srv)
			mu.Unlock()
		}

		startf()
	}

	return emulated, nil
}

func newInstance(name string) (instance *server.Instance, err error) {
	cfg, err := config.NewConfig(configFile)
	if err != nil {
		return nil, err
	}

	cfg.Identity = name
	cfg.OverrideCertname = cfg.Identity
	cfg.Collectives = []string{"mcollective"}
	cfg.InitiatedByServer = true
	cfg.DisableSecurityProviderVerify = true
	cfg.Choria.StatusFilePath = ""
	cfg.LogLevel = "warn"

	cfg.DisableTLS = !enableTLS
	cfg.DisableTLSVerify = !enableTLSVerify
	cfg.Choria.UseSRVRecords = false
	cfg.Choria.SecurityProvider = "file"

	if cfg.DisableTLS {
		cfg.Choria.SSLDir = "/nonexisting"
	}

	ichoria, err := choria.NewWithConfig(cfg)
	if err != nil {
		return nil, err
	}

	for i := 1; i < collectiveCount; i++ {
		ichoria.Config.Collectives = append(ichoria.Config.Collectives, fmt.Sprintf("collective%d", i))
	}

	if len(brokers) > 0 {
		ichoria.Config.Choria.MiddlewareHosts = brokers
	}

	if credentials != "" {
		ichoria.Configuration().Choria.NatsCredentials = credentials
	}

	srv, err := server.NewInstance(ichoria)
	if err != nil {
		log.Errorf("Could not create instance %s: %s", name, err)
		return
	}

	wg.Add(1)
	err = srv.Run(ctx, wg)
	if err != nil {
		log.Errorf("Could not run instance %s: %s", name, err)
		return
	}

	for i := 0; i < agentCount; i++ {
		agent := NewEmulatedAgent(ichoria, i)
		err := srv.RegisterAgent(ctx, agent.Metadata().Name, agent)
		if err != nil {
			log.Errorf("Could not register agent %s: %s", agent.Metadata().Name, err)
		}
	}

	return srv, err
}
