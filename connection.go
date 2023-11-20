package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Answer example
// script_duration_seconds{script="success"} 5.003527
// script_success{script="success"} 1
func (config *Config) getMetrics() (*Answer, error) {
	answer := &Answer{}

	req, err := http.NewRequest(http.MethodGet, config.URL, nil)
	if err != nil {
		return answer, err
	}

	timeout := time.Second * time.Duration(config.TimeoutSec)

	ctx, cancel := context.WithTimeout(req.Context(), timeout)
	defer cancel()

	req = req.WithContext(ctx)
	c := http.DefaultClient

	resp, err := c.Do(req)
	if err != nil {
		var netErr net.Error

		if errors.As(err, &netErr) && err.(net.Error).Timeout() { //nolint:forcetypeassert
			answer.TimedOut = true
		}

		return answer, err
	}

	defer resp.Body.Close()

	var (
		bodyStr   string
		bodyBytes []byte
	)

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		bodyStr = string(bodyBytes)
	} else {
		errorStr := fmt.Sprintf("%s: %d", "HTTP status", resp.StatusCode)
		return nil, errors.New(errorStr)
	}

	err = config.convertBody(bodyStr, answer)
	if err != nil {
		return nil, err
	}

	return answer, nil
}

func (config *Config) convertBody(bodyStr string, answer *Answer) error {
	var expectedResults int

	for _, result := range strings.Split(bodyStr, "\n") {
		switch {
		case strings.HasPrefix(result, config.SuccessMetric):
			expectedResults++

			valStrSlice := strings.Split(result, config.SuccessMetric)
			valStr := strings.Trim(valStrSlice[1], " \n")

			converted, err := strconv.ParseInt(valStr, 10, 2)
			if err != nil {
				return errors.New("invalid answer: " + err.Error())
			}

			if converted == 1 {
				answer.Success = true
			}
		case strings.HasPrefix(result, config.DurationMetric):
			expectedResults++

			valStrSlice := strings.Split(result, config.DurationMetric)
			valStr := strings.Trim(valStrSlice[1], " \n")

			converted, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				return errors.New("invalid answer: " + err.Error())
			}

			answer.Duration = converted
		}
	}

	if expectedResults != 2 {
		return errors.New("invalid body")
	}

	return nil
}
