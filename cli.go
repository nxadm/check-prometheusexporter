package main

import (
	"flag"
	"fmt"
	"os"
)

const usage = appname + ", " + version + ".\n" +
	`Nagios/Icinga check to query metric endpoint of Prometheus exporters.
Author: ` + author + `.
Repo: ` + repo + `.
   _       _       _       _       _       _       _       _
_-(_)-  _-(_)-  _-(_)-  _-(")-  _-(_)-  _-(_)-  _-(_)-  _-(_)-
*(___)  *(___)  *(___)  *%%%%%  *(___)  *(___)  *(___)  *(___)
// \\   // \\   // \\   // \\   // \\   // \\   // \\   // \\

Usage:
  check-prometheusprobe
    -u <URL> -s <metric> -d <metric>
    -w <seconds> -c <seconds>
    [-t <seconds>]
  check-netscaler-activeservices -v
  check-netscaler-activeservices -h

Options:
  -u <URL>       Prometheus probe endpoint.
  -s <metric>    Metric (int) for success (1) or failure (0).
  -d <duration>  Metric (int) with the duration of the check.
  -w <seconds>   Duration after which the check is marked as WARNING.
  -c <seconds>   Duration after which the check is marked as CRITICAL.
  -t <seconds>   Duration after which the check will time out (CRITICAL).
  -v             Show the version of this program.
  -h             Show this screen.
`

func (config *Config) readParams() {
	help := flag.Bool("h", false, "")
	url := flag.String("u", "", "")
	progVersion := flag.Bool("v", false, "")
	successMetric := flag.String("s", "", "")
	durationMetric := flag.String("d", "", "")
	warningSec := flag.Int("w", 0, "")
	criticalSec := flag.Int("c", 0, "")
	timeoutSec := flag.Int("t", 0, "")

	// Set a custom usage message & parse it
	flag.Usage = func() { fmt.Println(usage) }
	flag.Parse()

	// Handle early exits
	switch {
	case *help == true:
		fmt.Println(usage)
		os.Exit(UNKNOWN)
	case *progVersion == true:
		fmt.Println(appname + ", " + version + ".")
		os.Exit(UNKNOWN)
	case *url == "" || *successMetric == "" || *durationMetric == "" ||
		*warningSec == 0 || *criticalSec == 0 || *timeoutSec == 0:
		fmt.Println("[UNKNOWN] Invalid of missing values for parameters.")
		os.Exit(UNKNOWN)
	}

	// Import the values
	config.Url = *url
	config.SuccessMetric = *successMetric
	config.DurationMetric = *durationMetric
	config.WarningSec = *warningSec
	config.CriticalSec = *criticalSec
	config.TimeoutSec = *timeoutSec
}
