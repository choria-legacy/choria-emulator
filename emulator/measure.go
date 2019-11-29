package emulator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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
	forceDirect       bool
	testCount         int
	payloadSize       int
	outDir            string
	description       string
	rpcTimeout        time.Duration
	expectedNodeCount int
	measureWorkers    int

	dt time.Duration
)

type Measure struct {
	TestDescription string        `json:"description"`
	OutputDir       string        `json:"output_dir"`
	PayloadSize     int           `json:"payload_size"`
	Count           int           `json:"count"`
	TimeOut         time.Duration `json:"rpc_timeout"`
	ExpectedNodes   int           `json:"expected_nodes"`
	Workers         int           `json:"workers"`

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
		TimeOut:         rpcTimeout,
		ExpectedNodes:   expectedNodeCount,
		Workers:         measureWorkers,
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

	if len(nodes) != m.ExpectedNodes {
		return fmt.Errorf("discovered %d nodes, expected %d", len(nodes), m.ExpectedNodes)
	}

	log.Infof("Performing %d tests against %d nodes with a payload of %d bytes and timeout %v, reports in %s", m.Count, len(nodes), m.PayloadSize, m.TimeOut, m.OutputDir)

	sort.Strings(nodes)

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
	j, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("could not json marshal test config: %s", err)
	}

	err = ioutil.WriteFile(filepath.Join(m.OutputDir, "suite.json"), j, 0644)
	if err != nil {
		return fmt.Errorf("could not save test suit config: %s", err)
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
		m.requestCSV.Write([]string{"test", "description", "payload_bytes", "expected", "discovered", "publish_duration", "request_duration", "failed", "ok", "noresponses", "unexpected", "responded"})
	}

	rstats := []string{
		strconv.Itoa(c),
		description,
		strconv.Itoa(payloadSize),
		strconv.Itoa(len(nodes)),
		strconv.Itoa(stats.DiscoveredCount()),
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

	if len(stats.NoResponseFrom()) > 0 {
		log.Errorf("Did not receive responses from %d nodes:", len(stats.NoResponseFrom()))
		choria.SliceGroups(stats.NoResponseFrom(), 5, func(nodes []string) {
			nr := []string{}
			for _, n := range nodes {
				if n != "" {
					nr = append(nr, n)
				}
			}

			log.Errorf(strings.Join(nr, ", "))
		})
	}

	logline := fmt.Sprintf("REQUEST %s #: %-4d OK: %d / %d ELAPSED TIME: %v", stats.RequestID, c, stats.OKCount(), len(nodes), pd+rd)
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
		mco.Timeout(m.TimeOut),
		mco.Workers(m.Workers),
	}

	if forceDirect {
		log.Debugf("...using %d discovered nodes in directed mode", len(nodes))
		opts = append(opts, mco.Targets(nodes), mco.DirectRequest())
	} else {
		log.Debugf("...using %d discovered nodes in directed mode", len(nodes))
		opts = append(opts, mco.Targets(nodes), mco.BroadcastRequest())
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
