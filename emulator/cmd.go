package emulator

import (
	"context"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/alecthomas/kingpin"
	"github.com/choria-io/go-choria/choria"
	"github.com/sirupsen/logrus"
)

func Run() error {
	app := kingpin.New("choria-emulator", "Emulator for Choria Networks")
	app.Flag("name", "Instance name prefix").Default("").StringVar(&name)
	app.Flag("instances", "Number of instances to start").Short('i').Required().IntVar(&instanceCount)
	app.Flag("agents", "Number of emulated agents to start").Short('a').Default("1").IntVar(&agentCount)
	app.Flag("collectives", "Number of emulated subcollectives to create").Default("1").IntVar(&collectiveCount)
	app.Flag("config", "Choria configuration file").Short('c').StringVar(&configFile)
	app.Flag("broker", "NATS Server pool, specify multiple times (eg one:4222)").StringsVar(&brokers)
	app.Flag("http-port", "Port to listen for /debug/vars").Short('p').Default("8080").IntVar(&statusPort)
	app.Flag("tls", "Enable TLS on the NATS connections").Default("false").BoolVar(&enableTLS)
	app.Flag("verify", "Enable TLS certificate verifications on the NATS connections").Default("false").BoolVar(&enableTLSVerify)
	app.Flag("secure", "Enable Choria protocol security").Default("false").BoolVar(&protocolSecure)

	kingpin.MustParse(app.Parse(os.Args[1:]))

	if name == "" {
		name, err = os.Hostname()
		if err != nil {
			panic(fmt.Sprintf("Name is not given and cannot determine hostname: %s", err.Error()))
		}
	}

	if configFile == "" {
		configFile = choria.UserConfig()
	}

	fw, err = choria.New(configFile)
	if err != nil {
		return err
	}

	logger := logrus.New()
	logger.Out = os.Stdout
	log = logrus.NewEntry(logger)

	wg = &sync.WaitGroup{}
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	go startHTTP()
	startInstances()

	wg.Wait()

	return nil
}

func startInstances() error {
	instances, err = NewEmulator()

	return err
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
}
