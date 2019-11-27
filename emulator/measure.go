package emulator

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/choria-io/go-client/client"
	"github.com/choria-io/go-client/discovery/broadcast"
	"github.com/choria-io/go-protocol/protocol"
	"github.com/choria-io/mcorpc-agent-provider/mcorpc"
	mco "github.com/choria-io/mcorpc-agent-provider/mcorpc/client"
)

var (
	forceDirect bool
	testCount   int
	payloadSize int
	outDir      string
	description string

	dt time.Duration
)

type measureStats struct {
	times []time.Duration
	stats *mco.Stats
	sync.Mutex
}

func startMeasure() error {
	requestMetrics, err := os.Create(filepath.Join(outDir, "requests.csv"))
	if err != nil {
		return fmt.Errorf("could not open stats: %w", err)
	}
	defer requestMetrics.Close()

	requestTimes, err := os.Create(filepath.Join(outDir, "times.csv"))
	if err != nil {
		return fmt.Errorf("could not open stats: %w", err)
	}
	defer requestTimes.Close()

	log.Infof("Performing discovery of nodes with the emulator0 agent")
	nodes, err := discover()
	if err != nil {
		return fmt.Errorf("could not perform discovery: %w", err)
	}

	if len(nodes) == 0 {
		return fmt.Errorf("did not discover any nodes")
	}

	log.Infof("Performing %d tests against %d nodes with a payload of %d bytes, reports in %s", testCount, len(nodes), payloadSize, outDir)

	rm := csv.NewWriter(requestMetrics)
	rt := csv.NewWriter(requestTimes)

	for i := 1; i <= testCount; i++ {
		runTest(i, nodes, rm, rt)
	}

	rm.Flush()
	if err = rm.Error(); err != nil {
		log.Errorf("could not flush request metrics: %s", err)
	}

	rt.Flush()
	if err = rt.Error(); err != nil {
		log.Errorf("could not flush request times: %s", err)
	}

	return nil
}

func runTest(c int, nodes []string, rm *csv.Writer, rt *csv.Writer) error {
	log.Debugf("Starting test %d", c)

	type genRequest struct {
		Size int `json:"size"`
	}

	rpc, err := mco.New(fw, "emulated0")
	if err != nil {
		return err
	}

	times := []time.Duration{}
	mu := &sync.Mutex{}
	var startTime time.Time

	f := func(pr protocol.Reply, rpcr *mco.RPCReply) {
		mu.Lock()
		times = append(times, time.Now().Sub(startTime))
		mu.Unlock()

		if rpcr.Statuscode != mcorpc.OK {
			log.Errorf("......response from %s: %s", pr.SenderID(), rpcr.Statusmsg)
		}
	}

	opts := []mco.RequestOption{
		mco.ReplyHandler(f),
	}

	if forceDirect {
		log.Debugf("...using %d discovered nodes in directed mode", len(nodes))
		opts = append(opts, mco.Targets(nodes), mco.DirectRequest())
	} else {
		log.Debugf("...forcing discovery of the emulated0 nodes for %v", dt)
		f, err := client.NewFilter(client.AgentFilter("emulated0"))
		if err != nil {
			return err
		}

		opts = append(opts, mco.Filter(f))
		opts = append(opts, mco.DiscoveryTimeout(dt))
	}

	log.Debugf("...starting request")

	startTime = time.Now()
	result, err := rpc.Do(ctx, "generate", &genRequest{payloadSize}, opts...)
	if err != nil {
		return err
	}

	stats := result.Stats()

	stimes := make([]string, len(times))
	for i, t := range times {
		stimes[i] = strconv.Itoa(int(t))
	}

	rt.Write(stimes)
	err = rt.Error()
	if err != nil {
		log.Errorf("could not write times: %s", err)
	}

	dd := time.Duration(0)
	if !forceDirect {
		dd, err = stats.DiscoveryDuration()
		if err != nil {
			log.Errorf("could not determine discovery time: %s", err)
		}
	}

	pd, err := stats.PublishDuration()
	if err != nil {
		log.Errorf("could not determine publish time: %s", err)
	}

	rd, err := stats.RequestDuration()
	if err != nil {
		log.Errorf("could not determine total duration: %s", err)
	}

	if c == 1 {
		rm.Write([]string{"test", "description", "payload_bytes", "expected", "discovered", "discovery_duration", "publish_duration", "request_duration", "failed", "ok", "noresponses", "unexpected", "responded"})
	}

	rstats := []string{
		strconv.Itoa(c),
		description,
		strconv.Itoa(payloadSize),
		strconv.Itoa(len(nodes)),
		strconv.Itoa(stats.DiscoveredCount()),
		strconv.Itoa(int(dd)),
		strconv.Itoa(int(pd)),
		strconv.Itoa(int(rd)),
		strconv.Itoa(stats.FailCount()),
		strconv.Itoa(stats.OKCount()),
		strconv.Itoa(len(stats.NoResponseFrom())),
		strconv.Itoa(len(stats.UnexpectedResponseFrom())),
		strconv.Itoa(stats.ResponsesCount()),
	}

	rm.Write(rstats)
	err = rm.Error()
	if err != nil {
		log.Errorf("could not write measurement: %s", err)
	}

	log.Debugf("...finished request")

	logline := fmt.Sprintf("REQUEST #: %-4d OK: %d / %d ELAPSED TIME: %v", c, stats.OKCount(), len(nodes), dd+pd+rd)
	if stats.OKCount() == len(nodes) {
		log.Infof(logline)
	} else {
		log.Errorf(logline)
	}

	return nil
}

func discover() ([]string, error) {
	b := broadcast.New(fw)

	f, err := client.NewFilter(client.AgentFilter("emulated0"))
	if err != nil {
		return []string{}, err
	}

	return b.Discover(ctx, broadcast.Filter(f), broadcast.Timeout(dt))
}
