package main

import (
	"fmt"
	"os"
)

const (
	appname = "check-prometheusexporter"
	author  = "Claudio Ramirez <pub.claudio@gmail.com>"
	repo    = "https://github.com/nxadm/" + appname + ".git"
)

/* Nagios exit status */
const (
	OK = iota
	WARNING
	CRITICAL
	UNKNOWN
)

var version = "development" // to be overwritten by tag at buildtime with ld flags "-X main.version=tag"

type Answer struct {
	Success, TimedOut bool
	Duration          float64
}

type Config struct {
	URL, SuccessMetric, DurationMetric  string
	WarningSec, CriticalSec, TimeoutSec int
}

func main() {
	config := handleCLI()
	msg, exitCode := config.getExitInfo(config.getMetrics())
	fmt.Println(msg)
	os.Exit(exitCode)
}

func (config *Config) getExitInfo(answer *Answer, err error) (string, int) {
	if err != nil {
		if answer != nil && answer.TimedOut {
			return fmt.Sprintf("[CRITICAL] Timeout (%d) reached: %v\n", config.TimeoutSec, err), CRITICAL
		}

		return fmt.Sprintf("[UNKNOWN] Can not decode the server's answer: %v\n", err), UNKNOWN
	}

	if !answer.Success {
		return fmt.Sprintln("[CRITICAL] The service reports a failure."), CRITICAL
	}

	switch {
	case answer.Duration >= float64(config.CriticalSec):
		return fmt.Sprintf(
			"[CRITICAL] Check duration (%f) was higher than critical threshold (%d).\n",
			answer.Duration, config.CriticalSec), CRITICAL
	case answer.Duration >= float64(config.WarningSec):
		return fmt.Sprintf(
			"[WARNING] Check duration (%f) was higher than warning threshold (%d).\n",
			answer.Duration, config.WarningSec), WARNING
	default:
		return fmt.Sprintf(
			"[OK] Check duration (%f) lower than thresholds (critical: %d, warning: %d).\n",
			answer.Duration, config.CriticalSec, config.WarningSec), OK
	}
}
