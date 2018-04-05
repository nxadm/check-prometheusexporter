package main

import "fmt"

const appname = "check-prometheusprobe"
const version = "0.1.0"
const author = "Claudio Ramirez <pub.claudio@gmail.com>"
const repo = "https://github.com/nxadm/check-prometheusprobe.git"

type Config struct {
	Url, SuccessMetric, DurationMetric  string
	WarningSec, CriticalSec, TimeoutSec int
}

func main() {
	config := Config{}
	config.readParams()
	answer, err := getMetrics(&config)
	fmt.Printf("%+v\n", err)
	fmt.Printf("%+v\n", answer)
}
