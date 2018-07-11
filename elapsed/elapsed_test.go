package elapsed

import (
	"testing"
	"time"

	"github.com/bmizerany/assert"

	"github.com/hydah/golib/logger"
)

func TestElapsedTime(t *testing.T) {
	var et *ElapsedTime
	et = NewElapsedTime()
	time.Sleep(1 * time.Second)
	elapsed := et.Elapsed()

	printable := false
	if printable {
		logger.Debug("elapsed time: %d milliseconds", elapsed)
	}

	assert.Equal(t, elapsed, (int64)(1000))
}
