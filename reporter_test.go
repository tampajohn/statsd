package statsd

import (
	"testing"
	"time"
)

func init() {

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
	client = newMockStatter(false)
	config = checkConfig(&MetricsConfig{
		EnvName:              "testing",
		StuckFunctionTimeout: time.Second,
	})
	var testReporter = NewReporter()
	testReporter.ReportCall(testTags)
}

func TestReportFuncCallWithoutTiming(t *testing.T) {
	client = newMockStatter(false)
	config = checkConfig(&MetricsConfig{
		EnvName:              "testing",
		StuckFunctionTimeout: time.Second,
	})
	var testReporter = NewReporterFromConfig(&ReporterConfig{
		ReportTiming: false,
		FuncName:     "some-test",
	})
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

func TestNewReporterFromConfig(t *testing.T) {
	var reporter = NewReporterFromConfig(&ReporterConfig{
		ReportTiming: true,
		FuncName:     "something-arbitrary",
	})
	reporter.ReportError(testTags)
}

func TestNewReporterWithEmptyTags(t *testing.T) {
	var reporter = NewReporterFromConfig(&ReporterConfig{
		ReportTiming: true,
		FuncName:     "something-arbitrary",
	})
	reporter.ReportCall()
}

func TestNewReporterWithExtraTags(t *testing.T) {
	var reporter = NewReporterFromConfig(&ReporterConfig{
		ReportTiming: true,
		FuncName:     "something-arbitrary",
	})
	reporter.ReportCall(testTags.With("extra", "thing"))
}

func TestNewReporterWithEmptyExtraTags(t *testing.T) {
	var reporter = NewReporterFromConfig(&ReporterConfig{
		ReportTiming: true,
		FuncName:     "something-arbitrary",
	})
	reporter.ReportCall(MetricTags{}.With("extra", "thing"))
}
