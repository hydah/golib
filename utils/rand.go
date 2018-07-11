package utils

import (
	"math/rand"
	"time"
)

var RandChan = make(chan int, 1024)

func init() {
	rand.Seed(time.Now().UnixNano())
	go func() {
		for {
			RandChan <- rand.Int()
		}
	}()
}

func RandInt64(min, max int64) int64 {
	if min > max {
		return max
	}
	return min + rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(max-min)
}
