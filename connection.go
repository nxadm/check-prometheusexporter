package main

import (
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
	client := &http.Client{Timeout: time.Second * time.Duration(config.TimeoutSec)}

	/* Retrieve the metrics */
	request, err := http.NewRequest("GET", config.Url, nil)
	if err != nil {
		return checkTimeout(answer, err)
	}
	//request.SetBasicAuth(config.User, config.Pass)

	response, err := client.Do(request)
	if err != nil {
		return checkTimeout(answer, err)
	}
	defer response.Body.Close()
	var bodyStr string
	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		bodyStr = string(bodyBytes)
	} else {
		errorStr :=
			fmt.Sprintf("%s: %d", "Received error HTTP status", response.StatusCode)
		return nil, errors.New(errorStr)
	}

	/* Convert the metrics */
	info := strings.Split(bodyStr, "\n")
	for _, result := range info {
		switch {
		case strings.HasPrefix(result, config.SuccessMetric):
			valStrSlice := strings.Split(result, config.SuccessMetric)
			valStr := strings.Trim(valStrSlice[1], " \n")
			converted, err := strconv.ParseInt(valStr, 10, 2)
			if err != nil {
				return nil, errors.New("Invalid answer: " + err.Error())
			}
			if converted == 1 {
				answer.Success = true
			}
		case strings.HasPrefix(result, config.DurationMetric):
			valStrSlice := strings.Split(result, config.DurationMetric)
			valStr := strings.Trim(valStrSlice[1], " \n")
			converted, err := strconv.ParseFloat(valStr, 8)
			if err != nil {
				return nil, errors.New("Invalid answer: " + err.Error())
			}
			answer.Duration = converted
		}
	}

	return answer, err
}

func checkTimeout(answer *Answer, err error) (*Answer, error) {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		answer.TimedOut = true
	}
	return answer, err
}
