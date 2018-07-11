package main

import (
	"time"

	"github.com/hydah/golib/logger"
)

func main() {
	log := logger.NewLogger()
	log.AddFilter("stdout", logger.DEBUG, logger.NewStdoutLogWriter())
	log.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))
}
