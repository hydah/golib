package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

// stdout _
var stdout = os.Stdout
var stdoutCloseSignal chan bool

// StdoutLogWriter : This is the standard writer that prints to standard output.
type StdoutLogWriter chan *LogRecord

// NewStdoutLogWriter : This creates a new StdoutLogWriter
func NewStdoutLogWriter() StdoutLogWriter {
	records := make(StdoutLogWriter, LogBufferLength)
	stdoutCloseSignal = make(chan bool)
	go records.run(stdout)
	return records
}

func (w StdoutLogWriter) run(out io.Writer) {
	var timestr string
	var timestrAt int64
	var secondstr string

	for rec := range w {
		if at := rec.Created.UnixNano() / 1e9; at != timestrAt {
			timestr, timestrAt = rec.Created.Format("01/02/06 15:04:05"), at
		}
		secondstr = fmt.Sprintf("%s.%03d", timestr, rec.Created.Nanosecond()/1000000)

		index := 46
		source := rec.Source
		if len(rec.Source) > index {
			source = ".." + rec.Source[(len(rec.Source)-index):]
		}
		fmt.Fprint(out, "[", secondstr, "] [", levelStrings[rec.Level], "] [", source, "] ", rec.Message, "\n")
	}
	stdoutCloseSignal <- true
}

// LogWrite : This is the StdoutLogWriter's output method.
func (w StdoutLogWriter) LogWrite(rec *LogRecord) {
	select {
	case w <- rec:
	default:
	}
}

// Close : stops the logger from sending messages to standard output.  Attempts to
// send log messages to this logger after a Close have undefined behavior.
func (w StdoutLogWriter) Close() {
	close(w)
	t := time.NewTimer(1 * time.Second)
	select {
	case <-stdoutCloseSignal:
		return
	case <-t.C:
		return
	}
}
