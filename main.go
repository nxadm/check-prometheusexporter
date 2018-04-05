package main

import (
	"fmt"
	"os"
)

const appname = "check-prometheusprobe"
const version = "0.1.1"
const author = "Claudio Ramirez <pub.claudio@gmail.com>"
const repo = "https://github.com/nxadm/check-prometheusprobe.git"

/* Nagios exit status */
const (
	OK = iota
	WARNING
	CRITICAL
	UNKNOWN
)

type Answer struct {
	Success, TimedOut bool
	Duration          float64
}

type Config struct {
	Url, SuccessMetric, DurationMetric  string
	WarningSec, CriticalSec, TimeoutSec int
}

func main() {
	/* Read the CLI parameters */
	config := Config{}
	config.readParams()

	/* Query the probe */
	answer, err := getMetrics(&config)
	// fmt.Printf("ANSWER: %+v\n", answer) //debug
	if err != nil {
		if answer.TimedOut {
			fmt.Printf("[CRITICAL] Timeout (%d) reached: %v\n", config.TimeoutSec, err)
			os.Exit(CRITICAL)
		} else {
			fmt.Printf("[UNKNOWN] Can not decode the server's answer: %v\n", err)
			os.Exit(UNKNOWN)
		}
	}

	if answer.Success == false {
		fmt.Println("[CRITICAL] The service reports a failure.")
		os.Exit(CRITICAL)

	}

	switch {
	case answer.Duration >= float64(config.CriticalSec):
		fmt.Printf(
			"[CRITICAL] Check duration (%f) was higher than critical threshold (%d).\n",
			answer.Duration, config.CriticalSec)
		os.Exit(CRITICAL)
	case answer.Duration >= float64(config.WarningSec):
		fmt.Printf(
			"[WARNING] Check duration (%f) was higher than warning threshold (%d).\n",
			answer.Duration, config.WarningSec)
		os.Exit(WARNING)
	default:
		fmt.Printf(
			"[OK] Check duration (%f) lower than thresholds (critical: %d, warning: %d).\n",
			answer.Duration, config.CriticalSec, config.WarningSec)
		os.Exit(OK)
	}

}
