package statsd

import (
	"errors"
	"log"
	"time"

	"github.com/tampajohn/statsd/memstatsd"
)

// test comment for codeship
var client Statter
var config *MetricsConfig
var errorReporter ErrorReporter
var noStatterErr = errors.New("No Statter provided, if youy wish to use a mockStatter, please set config.MockingEnabled = true")

type MetricsConfig struct {
	EnvName              string
	StuckFunctionTimeout time.Duration
	MockingEnabled       bool
}

type ErrorReporter interface {
	Errorf(parts ...interface{})
}

func (m *MetricsConfig) BaseTags() string {
	return ",env=" + config.EnvName
}

// Statter
type Statter interface {
	Count(bucket string, n interface{})
	Increment(bucket string)
	Gauge(bucket string, value interface{})
	Timing(bucket string, value interface{})
	Histogram(bucket string, value interface{})
	Unique(bucket string, value string)
	Close()
}

func Close() {
	if client == nil {
		return
	}
	client.Close()
}

func Disable() {
	config = checkConfig(nil)
	client = newMockStatter(true)
}

func Init(s Statter, r ErrorReporter, cfg *MetricsConfig) error {
	config = checkConfig(cfg)
	if config.MockingEnabled {
		// init a mock statter instead of real statsd client
		client = newMockStatter(false)
		return nil
	} else if s == nil {
		return noStatterErr
	}
	client = s
	return nil
}

func checkConfig(cfg *MetricsConfig) *MetricsConfig {
	if cfg == nil {
		cfg = &MetricsConfig{}
	}
	if cfg.StuckFunctionTimeout < time.Second {
		cfg.StuckFunctionTimeout = 5 * time.Minute
	}
	if len(cfg.EnvName) == 0 {
		cfg.EnvName = "local"
	}
	return cfg
}

func RunMemstatsd(envName string, d time.Duration) {
	if client == nil {
		return
	}
	m := memstatsd.New("memstatsd.", envName, proxy{client})
	m.Run(d)
}

type proxy struct {
	client Statter
}

func (p proxy) Timing(bucket string, d time.Duration) {
	if d < 0 {
		return
	}
	p.client.Timing(bucket, int(d/time.Millisecond))
}

func (p proxy) Gauge(bucket string, value int) {
	p.client.Gauge(bucket, value)
}

func newMockStatter(noop bool) Statter {
	return &mockStatter{
		noop: noop,
	}
}

type mockStatter struct {
	noop bool
}

func (s *mockStatter) Count(bucket string, n interface{}) {
	if s.noop {
		return
	}
	log.Printf("[STATTER] Count %s: %v", bucket, n)
}

func (s *mockStatter) Increment(bucket string) {
	if s.noop {
		return
	}
	log.Printf("[STATTER] Increment %s", bucket)
}

func (s *mockStatter) Gauge(bucket string, value interface{}) {
	if s.noop {
		return
	}
	log.Printf("[STATTER] Gauge %s: %v", bucket, value)
}

func (s *mockStatter) Timing(bucket string, value interface{}) {
	if s.noop {
		return
	}
	log.Printf("[STATTER] Timing %s: %v", bucket, value)
}

func (s *mockStatter) Histogram(bucket string, value interface{}) {
	if s.noop {
		return
	}
	log.Printf("[STATTER] Histogram %s: %v", bucket, value)
}

func (s *mockStatter) Unique(bucket string, value string) {
	if s.noop {
		return
	}
	log.Printf("[STATTER] Unique %s: %v", bucket, value)
}

func (s *mockStatter) Close() {
	if s.noop {
		return
	}
	log.Printf("[STATTER] Closed at %s", time.Now())
}
