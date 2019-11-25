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
	"github.com/pkg/errors"
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
	err             error
	ctx             context.Context
	cancel          func()
	wg              *sync.WaitGroup
	fw              *choria.Framework
	instances       []*server.Instance
	log             *logrus.Entry
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

	for i := 1; i <= instanceCount; i++ {
		name := fmt.Sprintf("%s-%d", name, i)
		log.Infof("Creating instance %s", name)
		srv, err := newInstance(name)
		if err != nil {
			return nil, errors.Wrapf(err, "could not start instance %d", i)
		}

		emulated = append(emulated, srv)
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

	srv, err := server.NewInstance(ichoria)
	if err != nil {
		return nil, errors.Wrapf(err, "could not start instance %d", instance)
	}

	wg.Add(1)
	err = srv.Run(ctx, wg)

	for i := 1; i <= agentCount; i++ {
		agent := NewEmulatedAgent(ichoria, i)
		err := srv.RegisterAgent(ctx, agent.Metadata().Name, agent)
		if err != nil {
			return nil, errors.Wrapf(err, "could not register agent %s", agent.Metadata().Name)
		}
	}

	return srv, err
}
