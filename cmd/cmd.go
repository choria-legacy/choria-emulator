package cmd

import (
	"expvar"
	"fmt"
	"net/http"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/choria-io/choria-emulator/emulator"
	"github.com/choria-io/go-choria/mcollective"
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
	servers         []string
	err             error
	wg              sync.WaitGroup
	choria          *mcollective.Choria
	instances       []*emulator.ChoriaEmulationInstance
)

func parseCLI() {
	app := kingpin.New("choria-emulator", "Emulator for Choria Networks")
	app.Author("R.I.Pienaar <rip@devco.net>")
	app.Version("0.0.1")
	app.Flag("name", "Instance name prefix").Default("").StringVar(&name)
	app.Flag("instances", "Number of instances to start").Short('i').Required().IntVar(&instanceCount)
	app.Flag("agents", "Number of emulated agents to start").Short('a').Default("1").IntVar(&agentCount)
	app.Flag("collectives", "Number of emulated subcollectives to create").Default("1").IntVar(&collectiveCount)
	app.Flag("config", "Choria configuration file").Short('c').StringVar(&configFile)
	app.Flag("tls", "Enable TLS on the NATS connections").Default("false").BoolVar(&enableTLS)
	app.Flag("verify", "Enable TLS certificate verifications on the NATS connections").Default("false").BoolVar(&enableTLSVerify)
	app.Flag("server", "NATS Server pool, specify multiple times (eg one:4222)").StringsVar(&servers)
	app.Flag("http-port", "Port to listen for /debug/vars").Short('p').Default("8080").IntVar(&statusPort)

	if name == "" {
		name, err = os.Hostname()
		if err != nil {
			panic(fmt.Sprintf("Name is not given and cannot determine hostname: %s", err.Error()))
		}
	}

	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func Run() {
	parseCLI()

	initChoria()

	go startHTTP()
	startInstances()

	wg.Wait()
}

func startInstances() {
	for instance := 0; instance <= instanceCount-1; instance++ {
		ichoria, err := mcollective.New(choria.Config.ConfigFile)
		if err != nil {
			panic(fmt.Sprintf("Could not initialize Choria for instance %d: %s", instance, err.Error()))
		}

		ichoria.Config.Identity = fmt.Sprintf("%s-%d", name, instance)
		ichoria.Config.Collectives = []string{"mcollective"}

		for i := 1; i < collectiveCount; i++ {
			collective := fmt.Sprintf("collective%d", i)
			ichoria.Config.Collectives = append(ichoria.Config.Collectives, collective)
		}

		if !enableTLS {
			ichoria.Config.DisableTLS = true
		}

		if !(enableTLS && !enableTLSVerify) {
			ichoria.Config.OverrideCertname = ichoria.Config.Identity
		}

		if len(servers) > 0 {
			ichoria.Config.Choria.MiddlewareHosts = servers
		}

		emu, err := emulator.NewInstance(ichoria)
		if err != nil {
			panic(fmt.Sprintf("Could not start emulator: %s", err.Error()))
		}

		err = emu.Init(agentCount)
		if err != nil {
			panic(err)
		}

		wg.Add(1)

		go emu.ProcessRequests(&wg)
	}
}

func parseServers() ([]mcollective.Server, error) {
	s, err := mcollective.StringHostsToServers(servers, "nats")

	return s, err
}

func startHTTP() {
	exportConfig()

	port := fmt.Sprintf(":%d", statusPort)
	log.Infof("Starting to listen on HTTP %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func exportConfig() {
	expvar.NewString("name").Set(name)
	c := expvar.NewMap("config")
	c.Add("instances", int64(instanceCount))
	c.Add("agents", int64(agentCount))
	c.Add("collectives", int64(collectiveCount))
	c.Add("pid", int64(os.Getpid()))

	if enableTLS {
		c.Add("TLS", 1)
	} else {
		c.Add("TLS", 0)
	}

	if enableTLSVerify {
		c.Add("tlsVerify", 1)
	} else {
		c.Add("tlsVerify", 0)
	}
}

func initChoria() {
	if configFile == "" {
		configFile = mcollective.UserConfig()
	}

	if choria, err = mcollective.New(configFile); err != nil {
		panic(fmt.Sprintf("Could not initialize Choria: %s", err.Error()))
	}
}
