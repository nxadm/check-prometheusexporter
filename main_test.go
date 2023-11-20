package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

// func (config *Config) getMetrics() (*Answer, error) {
func TestConfig_getExitInfo(t *testing.T) {
	config := &Config{
		SuccessMetric:  `script_success{script="probename"}`,
		DurationMetric: `script_duration_seconds{script="probename"}`,
		WarningSec:     5,
		CriticalSec:    10,
		TimeoutSec:     15,
	}
	msg, exitCode := config.getExitInfo(nil, errors.New("test error"))
	assert.Contains(t, msg, "Can not decode the server's answer")
	assert.Equal(t, UNKNOWN, exitCode)

	msg, exitCode = config.getExitInfo(&Answer{TimedOut: true}, errors.New("test error"))
	assert.Contains(t, msg, "Timeout")
	assert.Equal(t, CRITICAL, exitCode)

	msg, exitCode = config.getExitInfo(&Answer{Success: false}, nil)
	assert.Contains(t, msg, "failure")
	assert.Equal(t, CRITICAL, exitCode)

	msg, exitCode = config.getExitInfo(&Answer{Success: true, Duration: float64(1)}, nil)
	assert.Contains(t, msg, "lower")
	assert.Equal(t, OK, exitCode)

	msg, exitCode = config.getExitInfo(&Answer{Success: true, Duration: float64(6)}, nil)
	assert.Contains(t, msg, "higher than warning threshold")
	assert.Equal(t, WARNING, exitCode)

	msg, exitCode = config.getExitInfo(&Answer{Success: true, Duration: float64(11)}, nil)
	assert.Contains(t, msg, "higher than critical threshold")
	assert.Equal(t, CRITICAL, exitCode)
}
