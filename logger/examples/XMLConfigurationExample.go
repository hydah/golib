package main

import (
	"github.com/hydah/golib/logger"
)

func main() {
	// Load the configuration (isn't this easy?)
	logger.LoadConfiguration("example.xml")

	// And now we're ready!
	logger.Finest("This will only go to those of you really cool UDP kids!  If you change enabled=true.")
	logger.Debug("Oh no!  %d + %d = %d!", 2, 2, 2+2)
	logger.Info("About that time, eh chaps?")
}
