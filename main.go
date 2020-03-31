package main

import (
	"fmt"
	"os"
)

const appname = "check-prometheusexporter"
const version = "0.3.0"
const author = "Claudio Ramirez <pub.claudio@gmail.com>"
const repo = "https://github.com/nxadm/" + appname + ".git"

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
	cfg := handleCLI()
	fmt.Printf("%#v\n", cfg)


	/* Query the probe */
	answer, err := getMetrics(&cfg)
	fmt.Printf("ANSWER: %#v\n%#v\n", answer, err) //debug
	if err != nil {
		if answer != nil && answer.TimedOut {
			fmt.Printf("[CRITICAL] Timeout (%d) reached: %v\n", cfg.TimeoutSec, err)
			os.Exit(CRITICAL)
		} else {
			fmt.Printf("[UNKNOWN] Can not decode the server's answer: %v\n", err)
			os.Exit(UNKNOWN)
		}
	}
	//
	//if answer.Success == false {
	//	fmt.Println("[CRITICAL] The service reports a failure.")
	//	os.Exit(CRITICAL)
	//
	//}
	//
	//switch {
	//case answer.Duration >= float64(cfg.CriticalSec):
	//	fmt.Printf(
	//		"[CRITICAL] Check duration (%f) was higher than critical threshold (%d).\n",
	//		answer.Duration, cfg.CriticalSec)
	//	os.Exit(CRITICAL)
	//case answer.Duration >= float64(cfg.WarningSec):
	//	fmt.Printf(
	//		"[WARNING] Check duration (%f) was higher than warning threshold (%d).\n",
	//		answer.Duration, cfg.WarningSec)
	//	os.Exit(WARNING)
	//default:
	//	fmt.Printf(
	//		"[OK] Check duration (%f) lower than thresholds (critical: %d, warning: %d).\n",
	//		answer.Duration, cfg.CriticalSec, cfg.WarningSec)
	//	os.Exit(OK)
	//}

}
