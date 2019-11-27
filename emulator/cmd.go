package emulator

import (
	"context"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/choria-io/go-choria/choria"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

func Run() error {
	app := kingpin.New("choria-emulator", "Emulator for Choria Networks")

	emulate := app.Command("emulate", "Starts the emulator")
	emulate.Flag("name", "Instance name prefix").Default("").StringVar(&name)
	emulate.Flag("instances", "Number of instances to start").Short('i').Required().IntVar(&instanceCount)
	emulate.Flag("credentials", "NATS credentials to use when connecting").ExistingFileVar(&credentials)
	emulate.Flag("agents", "Number of emulated agents to start").Short('a').Default("1").IntVar(&agentCount)
	emulate.Flag("collectives", "Number of emulated subcollectives to create").Default("1").IntVar(&collectiveCount)
	emulate.Flag("config", "Choria configuration file").Short('c').StringVar(&configFile)
	emulate.Flag("server", "NATS Server pool, specify multiple times (eg one:4222)").StringsVar(&brokers)
	emulate.Flag("http-port", "Port to listen for /debug/vars").Short('p').Default("8080").IntVar(&statusPort)
	emulate.Flag("tls", "Enable TLS on the NATS connections").Default("false").BoolVar(&enableTLS)
	emulate.Flag("verify", "Enable TLS certificate verifications on the NATS connections").Default("false").BoolVar(&enableTLSVerify)
	emulate.Flag("secure", "Enable Choria protocol security").Default("false").BoolVar(&protocolSecure)

	measure := app.Command("measure", "Perform requests and records various metrics")
	measure.Arg("count", "Number of tests to run").Required().IntVar(&testCount)
	measure.Arg("description", "Test scenario description").Required().StringVar(&description)
	measure.Arg("outdir", "Directory to write reports to").Required().ExistingDirVar(&outDir)
	measure.Flag("force-direct", "Force direct mode communication").BoolVar(&forceDirect)
	measure.Flag("discovery-timeout", "Discovery Timeout").Default("2s").DurationVar(&dt)
	measure.Flag("size", "Payload size").Default("100").IntVar(&payloadSize)
	measure.Flag("config", "Choria configuration file").Short('c').StringVar(&configFile)
	measure.Flag("tls", "Enable TLS on the NATS connections").Default("false").BoolVar(&enableTLS)
	measure.Flag("verify", "Enable TLS certificate verifications on the NATS connections").Default("false").BoolVar(&enableTLSVerify)
	measure.Flag("secure", "Enable Choria protocol security").Default("false").BoolVar(&protocolSecure)

	command := kingpin.MustParse(app.Parse(os.Args[1:]))

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

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	switch command {
	case "measure":
		err = startMeasure()
	case "emulate":
		err = startEmulator()
	default:
		err = fmt.Errorf("unknown command %s", command)
	}

	return err
}

func startEmulator() error {
	if name == "" {
		name, err = os.Hostname()
		if err != nil {
			panic(fmt.Sprintf("Name is not given and cannot determine hostname: %s", err.Error()))
		}
	}

	wg = &sync.WaitGroup{}

	go startHTTP()

	time.Sleep(time.Second)

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
