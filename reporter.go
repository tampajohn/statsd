package statsd

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

type Reporter interface {
	ReportError(tags ...MetricTags)
	ReportCall(tags ...MetricTags) StopTimerFunc
}

type ReporterConfig struct {
	FuncName     string
	ReportTiming bool
}

type defaultFuncReporter struct {
	reportTiming bool
	funcName     string
}

func NewReporter() Reporter {
	return defaultFuncReporter{
		reportTiming: true,
		funcName:     funcName(),
	}
}

func NewReporterFromConfig(options ReporterConfig) Reporter {
	return defaultFuncReporter{
		reportTiming: options.ReportTiming,
		funcName:     options.FuncName,
	}
}

func (f defaultFuncReporter) ReportCall(tags ...MetricTags) StopTimerFunc {
	reportFunc(f.funcName, "called", tags...)
	if f.reportTiming {
		fmt.Println("timing")
		return f.reportFuncTiming(tags...)
	}
	return mockTiming
}

func (f defaultFuncReporter) ReportError(tags ...MetricTags) {
	reportFunc(f.funcName, "error", tags...)
}

func mockTiming() {
	// This really does nothing...
}

func reportFunc(funcName, action string, tags ...MetricTags) {
	if client == nil {
		return
	}
	tagSpec := config.BaseTags() + joinTags(tags...)
	tagSpec += ",func_name=" + funcName
	client.Increment(fmt.Sprintf("func.%v", action) + tagSpec)
}

type StopTimerFunc func()

func (f defaultFuncReporter) reportFuncTiming(tags ...MetricTags) StopTimerFunc {
	if client == nil {
		return func() {}
	}
	t := time.Now()
	tagSpec := config.BaseTags() + joinTags(tags...)
	tagSpec += ",func_name=" + f.funcName

	doneC := make(chan struct{})
	go func(name string, start time.Time) {
		ticker := time.NewTicker(config.StuckFunctionTimeout)
		select {
		case <-doneC:
			ticker.Stop()
			return
		case <-ticker.C:
			if reporter != nil {
				reporter.Errorf("detected stuck function: %s stuck for %v\nspec:%s", name, time.Since(start), tagSpec)
			}
			client.Increment("func.stuck" + tagSpec)
			ticker.Stop()
			return
		}
	}(f.funcName, t)

	return func() {
		d := time.Since(t)
		close(doneC)
		client.Timing("func.timing"+tagSpec, int(d/time.Millisecond))
	}
}

func funcName() string {
	pc, _, _, _ := runtime.Caller(2)
	fullName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fullName, "/")
	nameParts := strings.Split(parts[len(parts)-1], ".")
	return nameParts[len(nameParts)-1]
}

type MetricTags map[string]string

func (t MetricTags) With(k, v string) MetricTags {
	if t == nil || len(t) == 0 {
		return map[string]string{
			k: v,
		}
	}
	t[k] = v
	return t
}

func joinTags(tags ...MetricTags) string {
	if len(tags) == 0 {
		return ""
	}
	var str string
	for k, v := range tags[0] {
		str += fmt.Sprintf(",%s=%s", k, v)
	}
	return str
}
