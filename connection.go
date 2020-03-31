package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Answer example
//script_duration_seconds{script="success"} 5.003527
//script_success{script="success"} 1
func getMetrics(config *Config) (*Answer, error) {
	answer := &Answer{}

	req, err := http.NewRequest("GET", config.Url, nil)
	if err != nil {
		return checkTimeout(answer, err)
	}

	timeout := time.Second * time.Duration(config.TimeoutSec)
	ctx, cancel := context.WithTimeout(req.Context(), timeout)
	defer cancel()

	req = req.WithContext(ctx)

	//req.SetBasicAuth(config.User, config.Pass)

	c := http.DefaultClient
	resp, err := c.Do(req)
	if err != nil {
		return checkTimeout(answer, err)
	}

	defer resp.Body.Close()
	var bodyStr string
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		bodyStr = string(bodyBytes)
	} else {
		errorStr :=
			fmt.Sprintf("%s: %d", "HTTP status", resp.StatusCode)
		return nil, errors.New(errorStr)
	}

	/* Convert the metrics */
	var expectedResults int

	info := strings.Split(bodyStr, "\n")

	for _, result := range info {
		switch {
		case strings.HasPrefix(fmt.Sprintf(result), config.SuccessMetric):
			expectedResults++
			valStrSlice := strings.Split(result, config.SuccessMetric)
			valStr := strings.Trim(valStrSlice[1], " \n")
			converted, err := strconv.ParseInt(valStr, 10, 2)
			if err != nil {
				return nil, errors.New("invalid answer: " + err.Error())
			}
			if converted == 1 {
				answer.Success = true
			}
		case strings.HasPrefix(result, config.DurationMetric):
			expectedResults++
			valStrSlice := strings.Split(result, config.DurationMetric)
			valStr := strings.Trim(valStrSlice[1], " \n")
			converted, err := strconv.ParseFloat(valStr, 8)
			if err != nil {
				return nil, errors.New("invalid answer: " + err.Error())
			}
			answer.Duration = converted
		}
	}

	fmt.Printf("expected: %#v\n", expectedResults)
	if expectedResults != 2 {
		return nil, errors.New("invalid body")
	}

	return answer, nil
}

func checkTimeout(answer *Answer, err error) (*Answer, error) {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		answer.TimedOut = true
	}
	return answer, err
}
