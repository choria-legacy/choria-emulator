package emulator

import (
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/choria-io/go-choria/choria"
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

type Measure struct {
	TestDescription string `json:"description"`
	OutputDir       string `json:"output_dir"`
	PayloadSize     int    `json:"payload_size"`
	Count           int    `json:"count"`

	sync.Mutex

	fw         *choria.Framework
	requestCSV *CSV
	timesCSV   *CSV
	bucketsCSV *CSV
}

func NewMeasure() (*Measure, error) {
	return &Measure{
		TestDescription: description,
		OutputDir:       outDir,
		PayloadSize:     payloadSize,
		Count:           testCount,
		fw:              fw,
	}, nil
}

func (m *Measure) Close() {
	err := m.requestCSV.Close()
	if err != nil {
		log.Errorf("could not close request csv: %s", err)
	}

	err = m.timesCSV.Close()
	if err != nil {
		log.Errorf("could not close times csv: %s", err)
	}

	err = m.bucketsCSV.Close()
	if err != nil {
		log.Errorf("could not close time buckets csv: %s", err)
	}
}

func (m *Measure) MustOpen() {
	m.requestCSV = MustNewCSV(filepath.Join(m.OutputDir, "requests.csv"))
	m.timesCSV = MustNewCSV(filepath.Join(m.OutputDir, "times.csv"))
	m.bucketsCSV = MustNewCSV(filepath.Join(m.OutputDir, "time_buckets.csv"))
}

func (m *Measure) Measure() error {
	m.MustOpen()
	defer m.Close()

	if protocolSecure {
		log.Infof("Enabling Choria protocol security")
		protocol.Secure = "true"
	} else {
		log.Infof("Disabling Choria protocol security")
		protocol.Secure = "false"
	}

	log.Infof("Performing discovery of nodes with the emulator0 agent")
	nodes, err := m.discover()
	if err != nil {
		return fmt.Errorf("could not perform discovery: %w", err)
	}

	if len(nodes) == 0 {
		return fmt.Errorf("did not discover any nodes")
	}

	log.Infof("Performing %d tests against %d nodes with a payload of %d bytes, reports in %s", m.Count, len(nodes), m.PayloadSize, m.OutputDir)

	for i := 1; i <= testCount; i++ {
		err = m.runTest(i, nodes)
		if err != nil {
			log.Errorf("Test %d failed: %s", i, err)
		}
	}

	return nil
}

func (m *Measure) saveTimes(times []time.Duration) error {
	stimes := make([]string, len(times))
	for i, t := range times {
		stimes[i] = strconv.Itoa(int(t))
	}

	return m.timesCSV.Write(stimes)
}

func (m *Measure) saveTimeBuckets(times []time.Duration) error {
	return m.bucketsCSV.Write(m.timeBuckets(times, 0.01))
}

func (m *Measure) saveRequestData(c int, nodes []string, stats *mco.Stats) error {
	dd := time.Duration(0)
	if !forceDirect {
		dd, err = stats.DiscoveryDuration()
		if err != nil {
			return fmt.Errorf("could not determine discovery time: %s", err)
		}
	}

	pd, err := stats.PublishDuration()
	if err != nil {
		return fmt.Errorf("could not determine publish time: %s", err)
	}

	rd, err := stats.RequestDuration()
	if err != nil {
		return fmt.Errorf("could not determine total duration: %s", err)
	}

	if c == 1 {
		m.requestCSV.Write([]string{"test", "description", "payload_bytes", "expected", "discovered", "discovery_duration", "publish_duration", "request_duration", "failed", "ok", "noresponses", "unexpected", "responded"})
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

	err = m.requestCSV.Write(rstats)
	if err != nil {
		return fmt.Errorf("could not write measurement: %s", err)
	}

	logline := fmt.Sprintf("REQUEST #: %-4d OK: %d / %d ELAPSED TIME: %v", c, stats.OKCount(), len(nodes), dd+pd+rd)
	if stats.OKCount() == len(nodes) {
		log.Infof(logline)
	} else {
		log.Errorf(logline)
	}

	return nil
}

func (m *Measure) recordStats(c int, nodes []string, stats *mco.Stats, times []time.Duration) error {
	err := m.saveRequestData(c, nodes, stats)
	if err != nil {
		return fmt.Errorf("could not write request data: %s", err)
	}

	err = m.saveTimes(times)
	if err != nil {
		return fmt.Errorf("could not write times: %s", err)
	}

	err = m.saveTimeBuckets(times)
	if err != nil {
		return fmt.Errorf("could not write time buckets: %s", err)
	}

	return nil
}

func (m *Measure) runTest(c int, nodes []string) error {
	log.Debugf("Starting test %d", c)

	type genRequest struct {
		Size int `json:"size"`
	}

	rpc, err := mco.New(fw, "emulated0")
	if err != nil {
		return err
	}

	times := []time.Duration{}
	var startTime time.Time

	f := func(pr protocol.Reply, rpcr *mco.RPCReply) {
		m.Lock()
		times = append(times, time.Now().Sub(startTime))
		m.Unlock()

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

		opts = append(opts, mco.BroadcastRequest())
		opts = append(opts, mco.Filter(f))
		opts = append(opts, mco.DiscoveryTimeout(dt))
	}

	log.Debugf("...starting request")

	startTime = time.Now()
	result, err := rpc.Do(ctx, "generate", &genRequest{payloadSize}, opts...)
	if err != nil {
		return err
	}

	err = m.recordStats(c, nodes, result.Stats(), times)
	if err != nil {
		log.Errorf("could not save stats: %s", err)
	}

	log.Debugf("...finished request")

	return nil
}

func (m *Measure) discover() ([]string, error) {
	b := broadcast.New(m.fw)

	f, err := client.NewFilter(client.AgentFilter("emulated0"))
	if err != nil {
		return []string{}, err
	}

	return b.Discover(ctx, broadcast.Filter(f), broadcast.Timeout(dt))
}

func (m *Measure) timeBuckets(times []time.Duration, interval float64) []string {
	max := times[len(times)-1]
	nbuckets := math.Round(float64(max.Seconds())/interval) + 1
	buckets := make([]int, int(nbuckets))

	for i := 0; i < len(buckets); i++ {
		buckets[i] = 0
	}

	for _, t := range times {
		buckets[int(math.Round(float64(t.Seconds())/interval))]++
	}

	out := make([]string, len(buckets))
	for i, c := range buckets {
		out[i] = strconv.Itoa(c)
	}

	return out
}
