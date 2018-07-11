package elapsed

import (
	"time"
)

type ElapsedTime struct {
	start int64
}

func NewElapsedTime() *ElapsedTime {
	return &ElapsedTime{start: time.Now().UnixNano()}
}

func (e *ElapsedTime) Elapsed() int64 {
	stop := time.Now().UnixNano()
	return int64(stop-e.start) / 1000000
}

func (e *ElapsedTime) Reset() {
	e.start = time.Now().UnixNano()
}
