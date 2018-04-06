# check-prometheusexporter

[![Build Status](https://travis-ci.org/nxadm/check-prometheusexporter.svg?branch=master)](https://travis-ci.org/nxadm/check-prometheusexporter)

Nagios/Icinga check for Prometheus exporter endpoints. This program designed
to work with
[Prometheus exporters](https://prometheus.io/docs/instrumenting/exporters/),
like [script_exporter](https://github.com/adhocteam/script_exporter). A
Prometheus setup is not needed to use this monitoring script, just an
endpoint exposing the metrics.

# Usage

```bash
$ ./check-prometheusexporter -h
check-prometheusexporter, 0.2.0.
Nagios/Icinga check to query metric endpoint of Prometheus exporters.
Author: Claudio Ramirez <pub.claudio@gmail.com>.
Repo: https://github.com/nxadm/check-prometheusexporter.git.
   _       _       _       _       _       _       _       _
_-(_)-  _-(_)-  _-(_)-  _-(")-  _-(_)-  _-(_)-  _-(_)-  _-(_)-
*(___)  *(___)  *(___)  *%%%%%  *(___)  *(___)  *(___)  *(___)
// \\   // \\   // \\   // \\   // \\   // \\   // \\   // \\

Usage:
  check-prometheusexporter
    -u <URL> -s <metric> -d <metric>
    -w <seconds> -c <seconds>
    [-t <seconds>]
  check-prometheusexporter -v
  check-prometheusexporter -h

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
  check-prometheusexporter -u http://somehost:9172/probe?name=success \
    -s 'script_success{script="success"}' \
    -d 'script_duration_seconds{script="success"}' \
    -w 6 -c 8 -t 15

```


