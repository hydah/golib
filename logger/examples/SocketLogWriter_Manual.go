package main

import (
	"time"

	"github.com/hydah/golib/logger"
)

func main() {
	log := logger.NewLogger()
	log.AddFilter("network", logger.FINEST, logger.NewSocketLogWriter("udp", "192.168.1.255:12124"))

	// Run `nc -u -l -p 12124` or similar before you run this to see the following message
	log.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))

	// This makes sure the output stream buffer is written
	log.Close()
}
