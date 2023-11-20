package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testResultSuccess = "1"
	testResultTime    = "0.123456789"
	testBodyStr       = "script_success{script=\"probename\"}" + testResultSuccess + "\n" +
		"script_duration_seconds{script=\"probename\"} " + testResultTime + "\n"
)

// func (config *Config) getMetrics() (*Answer, error) {
func TestConfig_getMetrics(t *testing.T) {
	config := &Config{
		SuccessMetric:  `script_success{script="probename"}`,
		DurationMetric: `script_duration_seconds{script="probename"}`,
	}

	_, err := config.getMetrics()
	assert.Error(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(testBodyStr))
	}))

	defer server.Close()

	config.URL = server.URL
	config.TimeoutSec = 10

	answer, err := config.getMetrics()
	assert.NoError(t, err)
	assert.True(t, answer.Success)

	converted, _ := strconv.ParseFloat(testResultTime, 64)
	assert.Equal(t, converted, answer.Duration)
}

func TestConfig_convertBody(t *testing.T) {
	config := &Config{
		SuccessMetric:  `script_success{script="probename"}`,
		DurationMetric: `script_duration_seconds{script="probename"}`,
	}
	answer := &Answer{}
	assert.Error(t, config.convertBody("", answer))

	bodyStr := "script_success{script=\"probename\"} 1\nscript_duration_seconds{script=\"probename\"} 0.1\nscript_success{script=\"probename\"} 0\n"
	assert.Error(t, config.convertBody(bodyStr, answer))
	assert.NoError(t, config.convertBody(testBodyStr, answer))
}
