# check-prometheusprobe

[![Build Status](https://travis-ci.org/nxadm/check-prometheusprobe.svg?branch=master)](https://travis-ci.org/nxadm/check-prometheusprobe)

Nagios/Icinga check for Prometheus exporter endpoints. This program designed to
work with [Prometheus exporters](https://prometheus.io/docs/instrumenting/exporters/),
like [script_exporter](https://github.com/adhocteam/script_exporter). A Prometheus
setup is not needed to use this monitoring script, just and endpoint exposing
the metrics.

# Usage

```bash
$ ./check-prometheusprobe -h
check-prometheusprobe, 0.1.0.
Nagios/Icinga check to query metric endpoint of Prometheus exporters.
Author: Claudio Ramirez <pub.claudio@gmail.com>.
Repo: https://github.com/nxadm/check-prometheusprobe.git.
   _       _       _       _       _       _       _       _
_-(_)-  _-(_)-  _-(_)-  _-(")-  _-(_)-  _-(_)-  _-(_)-  _-(_)-
*(___)  *(___)  *(___)  *%%%%%  *(___)  *(___)  *(___)  *(___)
// \\   // \\   // \\   // \\   // \\   // \\   // \\   // \\

Usage:
  check-prometheusprobe
    -u <URL> -s <metric> -d <metric>
    -w <seconds> -c <seconds>
    [-t <seconds>]
  check-prometheusprobe -v
  check-prometheusprobe -h

Options:
  -u <URL>       Prometheus probe endpoint.
  -s <metric>    Metric (int) for success (1) or failure (0).
  -d <duration>  Metric (int) with the duration of the check.
  -w <seconds>   Duration after which the check is marked as WARNING.
  -c <seconds>   Duration after which the check is marked as CRITICAL.
  -t <seconds>   Duration after which the check will time out (CRITICAL).
  -v             Show the version of this program.
  -h             Show this screen.

Example:
  check-prometheusprobe -u http://somehost:9172/probe?name=success \
    -s 'script_success{script="success"}' \
    -d 'script_duration_seconds{script="success"}' \
    -w 6 -c 8 -t 15

```


