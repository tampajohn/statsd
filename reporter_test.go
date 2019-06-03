package statsd

import (
	"testing"
	"time"
)

func init() {
	client = newMockStatter(false)
	config = checkConfig(&MetricsConfig{
		EnvName:              "testing",
		StuckFunctionTimeout: time.Second,
	})
	Init(client, nil, config)
}

var testTags = MetricTags{
	"layer":   "service",
	"service": "users",
}

var testErrTags = MetricTags{
	"layer":   "service",
	"service": "users",
	"error":   "an-error",
}

func TestReportFuncCall(t *testing.T) {
	var testReporter = NewReporter()
	testReporter.ReportCall(testTags)
}

func TestReportFuncError(t *testing.T) {
	var testReporter = NewReporter()
	testReporter.ReportError(testErrTags)
}

func TestReportFuncTiming(t *testing.T) {
	var testReporter = NewReporter()
	stopFn := testReporter.ReportCall(testTags)
	time.Sleep(500 * time.Millisecond)
	stopFn()
}

func TestReportFuncTimingStuck(t *testing.T) {
	var testReporter = NewReporter()
	stopFn := testReporter.ReportCall(testTags)
	time.Sleep(2 * time.Second)
	stopFn()
}
