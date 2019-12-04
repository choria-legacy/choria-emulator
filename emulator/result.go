package emulator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	mco "github.com/choria-io/mcorpc-agent-provider/mcorpc/client"
)

type MeasureResult struct {
	Suite               *Measure        `json:"-"`
	SuiteID             string          `json:"suite_id"`
	Instance            int             `json:"instance"`
	TimeBuckets         []int           `json:"time_buckets"`
	Times               []time.Duration `json:"-"`
	Discovered          int             `json:"discovered"`
	PublishDuration     time.Duration   `json:"publish_duration"`
	RequestDuration     time.Duration   `json:"request_duration"`
	FailedCount         int             `json:"failed"`
	OKCount             int             `json:"ok"`
	NoResponses         int             `json:"noresponses"`
	UnexpectedResponses int             `json:"unexpected"`
	ResponsesReceived   int             `json:"responded"`

	sync.Mutex
	requestStats *mco.Stats
}

func NewMeasureResult(m *Measure) *MeasureResult {
	return &MeasureResult{Suite: m, SuiteID: m.ID}
}

func (r *MeasureResult) SaveAll() []error {
	errs := []error{}

	err := r.saveSuite()
	if err != nil {
		errs = append(errs, fmt.Errorf("could not save suite: %s", err))
	}

	err = r.saveTimes(r.Suite.timesCSV)
	if err != nil {
		errs = append(errs, fmt.Errorf("could not save times: %s", err))
	}

	err = r.saveTimeBuckets(r.Suite.bucketsCSV)
	if err != nil {
		errs = append(errs, fmt.Errorf("could not save time buckets: %s", err))
	}

	err = r.saveRequestData(r.Suite.requestCSV)
	if err != nil {
		errs = append(errs, fmt.Errorf("could not save request data: %s", err))
	}

	err = r.PublishCloudEvent()
	if err != nil {
		errs = append(errs, fmt.Errorf("could not publish cloud event: %s", err))
	}

	return errs
}

func (r *MeasureResult) PublishCloudEvent() error {
	event := r.Suite.newEvent("emulator.test_result")
	event.SetData(r)

	return r.Suite.publishEvent(event)
}

func (r *MeasureResult) RecordTime(t time.Duration) {
	r.Lock()
	r.Times = append(r.Times, t)
	r.Unlock()
}

func (r *MeasureResult) SetMCOResult(instance int, stats *mco.Stats) error {
	r.Instance = instance
	r.requestStats = stats

	var err error

	r.PublishDuration, err = stats.PublishDuration()
	if err != nil {
		return fmt.Errorf("could not determine publish time: %s", err)
	}

	r.RequestDuration, err = stats.RequestDuration()
	if err != nil {
		return fmt.Errorf("could not determine total duration: %s", err)
	}

	r.Discovered = stats.DiscoveredCount()
	r.FailedCount = stats.FailCount()
	r.OKCount = stats.OKCount()
	r.NoResponses = len(stats.NoResponseFrom())
	r.UnexpectedResponses = len(stats.UnexpectedResponseFrom())
	r.ResponsesReceived = stats.ResponsesCount()

	return nil
}

func (r *MeasureResult) saveRequestData(csv *CSV) error {
	if r.Instance == 1 {
		csv.Write([]string{"test", "description", "payload_bytes", "expected", "discovered", "publish_duration", "request_duration", "failed", "ok", "noresponses", "unexpected", "responded"})
	}

	rstats := []string{
		strconv.Itoa(r.Instance),
		description,
		strconv.Itoa(payloadSize),
		strconv.Itoa(r.Suite.ExpectedNodes),
		strconv.Itoa(r.Discovered),
		strconv.Itoa(int(r.PublishDuration)),
		strconv.Itoa(int(r.RequestDuration)),
		strconv.Itoa(r.FailedCount),
		strconv.Itoa(r.OKCount),
		strconv.Itoa(r.NoResponses),
		strconv.Itoa(r.UnexpectedResponses),
		strconv.Itoa(r.ResponsesReceived),
	}

	return csv.Write(rstats)
}

func (r *MeasureResult) saveSuite() error {
	j, err := json.Marshal(r.Suite)
	if err != nil {
		return fmt.Errorf("could not json marshal test config: %s", err)
	}

	err = ioutil.WriteFile(filepath.Join(r.Suite.OutputDir, "suite.json"), j, 0644)
	if err != nil {
		return fmt.Errorf("could not save test suit config: %s", err)
	}

	return nil
}

func (r *MeasureResult) saveTimeBuckets(csv *CSV) error {
	return csv.Write(r.timeBuckets(0.01))
}

func (r *MeasureResult) timeBuckets(interval float64) []string {
	r.Lock()
	defer r.Unlock()

	if len(r.TimeBuckets) == 0 && len(r.Times) > 0 {
		max := r.Times[len(r.Times)-1]
		nbuckets := math.Round(float64(max.Seconds())/interval) + 1
		r.TimeBuckets = make([]int, int(nbuckets))

		for i := 0; i < len(r.TimeBuckets); i++ {
			r.TimeBuckets[i] = 0
		}

		for _, t := range r.Times {
			r.TimeBuckets[int(math.Round(float64(t.Seconds())/interval))]++
		}
	}

	out := make([]string, len(r.TimeBuckets))
	for i, c := range r.TimeBuckets {
		out[i] = strconv.Itoa(c)
	}

	return out
}

func (r *MeasureResult) saveTimes(csv *CSV) error {
	stimes := make([]string, len(r.Times))
	for i, t := range r.Times {
		stimes[i] = strconv.Itoa(int(t))
	}

	return csv.Write(stimes)
}
