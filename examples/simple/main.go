package main

import (
	"fmt"
	"time"

	"github.com/tampajohn/statsd"
)

var tags = map[string]string{
	"application": "simple-example",
}

func main() {
	statsd.Init(nil, nil, &statsd.MetricsConfig{
		EnvName:              "prod",
		MockingEnabled:       true,
		StuckFunctionTimeout: 2 * time.Second,
	})

	fmt.Println("doing something quickly")
	doSomethingQuickly()
	fmt.Println("done something quickly")
	go func() {
		fmt.Println("doing something slowly async")
		doSomethingSlowly()
		fmt.Println("done something slowly async")
	}()
	fmt.Println("doing something slowly")
	doSomethingSlowly()
	fmt.Println("done something slowly")
}

func doSomethingQuickly() {
	defer statsd.NewReporter().ReportCall(tags)()
}

func doSomethingSlowly() {
	defer statsd.NewReporter().ReportCall(tags)()
	time.Sleep(10 * time.Second)
}
