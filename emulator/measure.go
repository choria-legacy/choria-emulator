package emulator

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/choria-io/go-choria/choria"
	"github.com/choria-io/go-client/client"
	"github.com/choria-io/go-client/discovery/broadcast"
	"github.com/choria-io/go-protocol/protocol"
	"github.com/choria-io/mcorpc-agent-provider/mcorpc"
	mco "github.com/choria-io/mcorpc-agent-provider/mcorpc/client"
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/nats-io/nuid"
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

type EventPublisher interface {
	PublishRaw(target string, data []byte) error
	Close()
}

type Measure struct {
	TestDescription string        `json:"description"`
	OutputDir       string        `json:"output_dir"`
	PayloadSize     int           `json:"payload_size"`
	Count           int           `json:"count"`
	TimeOut         time.Duration `json:"rpc_timeout"`
	ExpectedNodes   int           `json:"expected_nodes"`
	Workers         int           `json:"workers"`
	ID              string        `json:"id"`

	eventsConn EventPublisher

	sync.Mutex

	fw         *choria.Framework
	requestCSV *CSV
	timesCSV   *CSV
	bucketsCSV *CSV
}

func NewMeasure() (*Measure, error) {
	id, _ := fw.NewRequestID()
	conn, err := fw.NewConnector(context.Background(), fw.MiddlewareServers, id, fw.Logger("cloudevents"))
	if err != nil {
		return nil, err
	}

	return &Measure{
		TestDescription: description,
		OutputDir:       outDir,
		PayloadSize:     payloadSize,
		Count:           testCount,
		TimeOut:         rpcTimeout,
		ExpectedNodes:   expectedNodeCount,
		Workers:         measureWorkers,
		ID:              nuid.New().Next(),
		eventsConn:      conn,
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

	err := m.PublishEvent()
	if err != nil {
		log.Warnf("Cloud not publish cloud event: %s", err)
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

func (m *Measure) newEvent(etype string) cloudevents.Event {
	event := cloudevents.NewEvent("1.0")
	event.SetType("io.choria.event." + etype)
	event.SetSource("io.choria.choria-emulator.measure")
	event.SetID(nuid.New().Next())

	return event
}

func (m *Measure) PublishEvent() error {
	event := m.newEvent("emulator.test_suite")
	event.SetData(m)

	return m.publishEvent(event)
}

func (m *Measure) publishEvent(event cloudevents.Event) error {
	j, err := event.MarshalJSON()
	if err != nil {
		return err
	}

	return m.eventsConn.PublishRaw(event.Type(), j)
}

func (m *Measure) recordStats(mr *MeasureResult) error {
	if mr.NoResponses > 0 {
		log.Errorf("Did not receive responses from %d nodes:", mr.NoResponses)
		choria.SliceGroups(mr.requestStats.NoResponseFrom(), 5, func(nodes []string) {
			nr := []string{}
			for _, n := range nodes {
				if n != "" {
					nr = append(nr, n)
				}
			}

			log.Errorf(strings.Join(nr, ", "))
		})
	}

	logline := fmt.Sprintf("REQUEST %s #: %-4d OK: %d / %d ELAPSED TIME: %v", mr.requestStats.RequestID, mr.Instance, mr.OKCount, m.ExpectedNodes, mr.RequestDuration)

	if mr.OKCount == m.ExpectedNodes {
		log.Infof(logline)
	} else {
		log.Errorf(logline)
	}

	errs := mr.SaveAll()
	if len(errs) > 0 {
		log.Errorf("Errors encountered while saving data")
		for _, err := range errs {
			log.Error(err)
		}
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

	result := NewMeasureResult(m)

	var startTime time.Time

	f := func(pr protocol.Reply, rpcr *mco.RPCReply) {
		result.RecordTime(time.Now().Sub(startTime))

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
	rpcres, err := rpc.Do(ctx, "generate", &genRequest{payloadSize}, opts...)
	if err != nil {
		return err
	}

	result.SetMCOResult(c, rpcres.Stats())

	err = m.recordStats(result)
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
