package statsd

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

func ReportFuncError(tags ...MetricTags) {
	fn := funcName()
	reportFunc(fn, "error", tags...)
}

func ReportFuncStatus(tags ...MetricTags) {
	fn := funcName()
	reportFunc(fn, "status", tags...)
}

func ReportFuncCall(tags ...MetricTags) {
	fn := funcName()
	reportFunc(fn, "called", tags...)
}

func reportFunc(fn, action string, tags ...MetricTags) {
	if client == nil {
		return
	}
	tagSpec := config.BaseTags() + joinTags(tags...)
	tagSpec += ",func_name=" + fn
	client.Increment(fmt.Sprintf("func.%v", action) + tagSpec)
}

type StopTimerFunc func()

func ReportFuncTiming(tags ...MetricTags) StopTimerFunc {
	if client == nil {
		return func() {}
	}
	t := time.Now()
	fn := funcName()
	tagSpec := config.BaseTags() + joinTags(tags...)
	tagSpec += ",func_name=" + fn

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
	}(fn, t)

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
