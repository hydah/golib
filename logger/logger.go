// Package logger provides level-based and highly configurable logging.
//
// Enhanced Logging
//
// This is inspired by the logging functionality in Java.  Essentially, you create a Logger
// object and create output filters for it.  You can send whatever you want to the Logger,
// and it will filter that based on your settings and send it to the outputs.  This way, you
// can put as much debug code in your program as you want, and when you're done you can filter
// out the mundane messages so only the important ones show up.
//
// Utility functions are provided to make life easier. Here is some example code to get started:
//
// log := logger.NewLogger()
// log.AddFilter("stdout", logger.DEBUG, logger.NewStdoutLogWriter())
// log.AddFilter("log",    logger.FINE,  logger.NewFileLogWriter("example.log", true))
// log.Info("The time is now: %s", time.LocalTime().Format("15:04:05 MST 2006/01/02"))
//
// The first two lines can be combined with the utility NewDefaultLogger:
//
// log := logger.NewDefaultLogger(logger.DEBUG)
// log.AddFilter("log",    logger.FINE,  logger.NewFileLogWriter("example.log", true))
// log.Info("The time is now: %s", time.LocalTime().Format("15:04:05 MST 2006/01/02"))
//
// Usage notes:
// - The StdoutLogWriter does not display the source of the message to standard
//   output, but the FileLogWriter does.
// - The utility functions (Info, Debug, Warn, etc) derive their source from the
//   calling function, and this incurs extra overhead.
//
// Changes from 2.0:
// - The external interface has remained mostly stable, but a lot of the
//   internals have been changed, so if you depended on any of this or created
//   your own LogWriter, then you will probably have to update your code.  In
//   particular, Logger is now a map and StdoutLogWriter is now a channel
//   behind-the-scenes, and the LogWrite method no longer has return values.
//

package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// These are the integer logging levels used by the logger
type level int

// LEVEL define
const (
	DEBUG level = iota
	TRACE
	INFO
	WARNING
	ERROR
)

// Logging level strings
var (
	levelStrings = [...]string{"DEBG", "TRAC", "INFO", "WARN", "EROR"}
)

func (l level) String() string {
	if l < 0 || int(l) > len(levelStrings) {
		return "UNKNOWN"
	}
	return levelStrings[int(l)]
}

var (
	// LogBufferLength specifies how many log messages a particular logger
	// logger can buffer at a time before writing them.
	LogBufferLength = 1024 * 4
)

// LogRecord contains all of the pertinent information for each message
type LogRecord struct {
	Level   level     // The log level
	Created time.Time // The time at which the log message was created (nanoseconds)
	Source  string    // The message source
	Message string    // The log message
}

// LogWriter : This is an interface for anything that should be able to write logs
type LogWriter interface {
	// This will be called to log a LogRecord message.
	LogWrite(rec *LogRecord)

	// This should clean up anything lingering about the LogWriter, as it is called before
	// the LogWriter is removed.  LogWrite should not be called after Close.
	Close()
}

// Filter represents the log level below which no log records are written to
// the associated LogWriter.
type Filter struct {
	Level level
	LogWriter
}

// Logger represents a collection of Filters through which log messages are
// written.
type Logger map[string]*Filter

// NewLogger : Create a new logger.
func NewLogger() Logger {
	os.Stderr.WriteString("warning: use of deprecated NewLogger\n")
	return make(Logger)
}

// NewStdoutLogger : Create a new logger with a "stdout" filter configured to send log messages at
// or above lvl to standard output.
//
// DEPRECATED: use NewDefaultLogger instead.
func NewStdoutLogger(lvl level) Logger {
	os.Stderr.WriteString("warning: use of deprecated NewStdoutLogger\n")
	return Logger{
		"stdout": &Filter{lvl, NewStdoutLogWriter()},
	}
}

// NewDefaultLogger : Create a new logger with a "stdout" filter configured to send log messages at
// or above lvl to standard output.
func NewDefaultLogger(lvl level) Logger {
	return Logger{
		"stdout": &Filter{lvl, NewStdoutLogWriter()},
	}
}

// Close : Closes all log writers in preparation for exiting the program or a
// reconfiguration of logging.  Calling this is not really imperative, unless
// you want to guarantee that all log messages are written.  Close removes
// all filters (and thus all LogWriters) from the logger.
func (log Logger) Close() {
	// Close all open loggers
	for name, filt := range log {
		delete(log, name)
		filt.Close()
	}
}

// AddFilter : Add a new LogWriter to the Logger which will only log messages at lvl or
// higher.  This function should not be called from multiple goroutines.
// Returns the logger for chaining.
func (log Logger) AddFilter(name string, lvl level, writer LogWriter) Logger {
	log[name] = &Filter{lvl, writer}
	return log
}

// Send a formatted log message internally
func (log Logger) intLogf(lvl level, format string, args ...interface{}) {
	skip := true

	// Determine if any logging will be done
	for _, filt := range log {
		if lvl >= filt.Level {
			skip = false
			break
		}
	}
	if skip {
		return
	}

	// Determine caller func
	pc, file, line, ok := runtime.Caller(2)
	file = filepath.Base(file)
	src := ""
	if ok {
		src = fmt.Sprintf("%s %s:%d", runtime.FuncForPC(pc).Name(), file, line)
	}

	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	if len(args) > 10240 {
		msg = msg[:1024] + "......"
	}

	// Make the log record
	rec := &LogRecord{
		Level:   lvl,
		Created: time.Now(),
		Source:  src,
		Message: msg,
	}

	// Dispatch the logs
	for _, filt := range log {
		if lvl < filt.Level {
			continue
		}
		filt.LogWrite(rec)
	}
}

// Send a closure log message internally
func (log Logger) intLogc(lvl level, closure func() string) {
	skip := true

	// Determine if any logging will be done
	for _, filt := range log {
		if lvl >= filt.Level {
			skip = false
			break
		}
	}
	if skip {
		return
	}

	// Determine caller func
	pc, file, line, ok := runtime.Caller(2)
	file = filepath.Base(file)
	src := ""
	if ok {
		src = fmt.Sprintf("%s %s:%d", runtime.FuncForPC(pc).Name(), file, line)
	}

	// Make the log record
	rec := &LogRecord{
		Level:   lvl,
		Created: time.Now(),
		Source:  src,
		Message: closure(),
	}

	// Dispatch the logs
	for _, filt := range log {
		if lvl < filt.Level {
			continue
		}
		filt.LogWrite(rec)
	}
}

// Log : Send a log message with manual level, source, and message.
func (log Logger) Log(lvl level, source, message string) {
	skip := true

	// Determine if any logging will be done
	for _, filt := range log {
		if lvl >= filt.Level {
			skip = false
			break
		}
	}
	if skip {
		return
	}

	// Make the log record
	rec := &LogRecord{
		Level:   lvl,
		Created: time.Now(),
		Source:  source,
		Message: message,
	}

	// Dispatch the logs
	for _, filt := range log {
		if lvl < filt.Level {
			continue
		}
		filt.LogWrite(rec)
	}
}

// Logf logs a formatted log message at the given log level, using the caller as
// its source.
func (log Logger) Logf(lvl level, format string, args ...interface{}) {
	log.intLogf(lvl, format, args...)
}

// Logc logs a string returned by the closure at the given log level, using the caller as
// its source.  If no log message would be written, the closure is never called.
func (log Logger) Logc(lvl level, closure func() string) {
	log.intLogc(lvl, closure)
}

// Debug is a utility method for debug log messages.
// The behavior of Debug depends on the first argument:
// - arg0 is a string
//   When given a string as the first argument, this behaves like Logf but with
//   the DEBUG log level: the first argument is interpreted as a format for the
//   latter arguments.
// - arg0 is a func()string
//   When given a closure of type func()string, this logs the string returned by
//   the closure iff it will be logged.  The closure runs at most one time.
// - arg0 is interface{}
//   When given anything else, the log message will be each of the arguments
//   formatted with %v and separated by spaces (ala Sprint).
func (log Logger) Debug(arg0 interface{}, args ...interface{}) {
	var (
		lvl = DEBUG
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		log.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		log.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		log.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Trace logs a message at the trace log level.
// See Debug for an explanation of the arguments.
func (log Logger) Trace(arg0 interface{}, args ...interface{}) {
	var (
		lvl = TRACE
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		log.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		log.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		log.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Info logs a message at the info log level.
// See Debug for an explanation of the arguments.
func (log Logger) Info(arg0 interface{}, args ...interface{}) {
	var (
		lvl = INFO
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		log.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		log.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		log.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Warn logs a message at the warning log level and returns the formatted error.
// At the warning level and higher, there is no performance benefit if the
// message is not actually logged, because all formats are processed and all
// closures are executed to format the error message.
// See Debug for further explanation of the arguments.
func (log Logger) Warn(arg0 interface{}, args ...interface{}) error {
	var (
		lvl = WARNING
	)
	var msg string
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg = fmt.Sprintf(first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		msg = first()
	default:
		// Build a format string so that it will be similar to Sprint
		msg = fmt.Sprintf(fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
	}
	log.intLogf(lvl, msg)
	return errors.New(msg)
}

// Error logs a message at the error log level and returns the formatted error,
// See Warn for an explanation of the performance and Debug for an explanation
// of the parameters.
func (log Logger) Error(arg0 interface{}, args ...interface{}) error {
	var (
		lvl = ERROR
	)
	var msg string
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg = fmt.Sprintf(first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		msg = first()
	default:
		// Build a format string so that it will be similar to Sprint
		msg = fmt.Sprintf(fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
	}
	log.intLogf(lvl, msg)
	return errors.New(msg)
}
