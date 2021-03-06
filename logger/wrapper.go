package logger

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/hydah/golib/config"
	"github.com/hydah/golib/logger/cfgstruct"
)

// logger define
var (
	Global Logger
)

func init() {
	Global = NewDefaultLogger(DEBUG)
}

// LoadConfiguration : Wrapper for (*Logger).LoadConfiguration
func LoadConfigurationConf(cfg *config.Config) {
	Global.LoadConfigurationConf(cfg)
}

func LoadConfigurationConfV2(cfg cfgstruct.LogTypeSt) {
	Global.LoadConfigurationConfV2(cfg)
}

// AddFilter : Wrapper for (*Logger).AddFilter
func AddFilter(name string, lvl level, writer LogWriter) {
	Global.AddFilter(name, lvl, writer)
}

// Close : Wrapper for (*Logger).Close (closes and removes all logwriters)
func Close() {
	Global.Close()
}

// Crash _
func Crash(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(ERROR, strings.Repeat(" %v", len(args))[1:], args...)
	}
	panic(args)
}

// Crashf : Logs the given message and crashes the program
func Crashf(format string, args ...interface{}) {
	Global.intLogf(ERROR, format, args...)
	Global.Close() // so that hopefully the messages get logged
	panic(fmt.Sprintf(format, args...))
}

// Exit : Compatibility with `log`
func Exit(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(ERROR, strings.Repeat(" %v", len(args))[1:], args...)
	}
	Global.Close() // so that hopefully the messages get logged
	os.Exit(0)
}

// Exitf : Compatibility with `log`
func Exitf(format string, args ...interface{}) {
	Global.intLogf(ERROR, format, args...)
	Global.Close() // so that hopefully the messages get logged
	os.Exit(0)
}

// Stderr : Compatibility with `log`
func Stderr(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(ERROR, strings.Repeat(" %v", len(args))[1:], args...)
	}
}

// Stderrf : Compatibility with `log`
func Stderrf(format string, args ...interface{}) {
	Global.intLogf(ERROR, format, args...)
}

// Stdout : Compatibility with `log`
func Stdout(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(INFO, strings.Repeat(" %v", len(args))[1:], args...)
	}
}

// Stdoutf : Compatibility with `log`
func Stdoutf(format string, args ...interface{}) {
	Global.intLogf(INFO, format, args...)
}

// Log : Send a log message manually
// Wrapper for (*Logger).Log
func Log(lvl level, source, message string) {
	Global.Log(lvl, source, message)
}

// Logf : Send a formatted log message easily
// Wrapper for (*Logger).Logf
func Logf(lvl level, format string, args ...interface{}) {
	Global.intLogf(lvl, format, args...)
}

// Logc : Send a closure log message
// Wrapper for (*Logger).Logc
func Logc(lvl level, closure func() string) {
	Global.intLogc(lvl, closure)
}

// Debug : Utility for debug log messages
// When given a string as the first argument, this behaves like Logf but with the DEBUG log level (e.g. the first argument is interpreted as a format for the latter arguments)
// When given a closure of type func()string, this logs the string returned by the closure iff it will be logged.  The closure runs at most one time.
// When given anything else, the log message will be each of the arguments formatted with %v and separated by spaces (ala Sprint).
// Wrapper for (*Logger).Debug
func Debug(arg0 interface{}, args ...interface{}) {
	var (
		lvl = DEBUG
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Trace : Utility for trace log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Trace
func Trace(arg0 interface{}, args ...interface{}) {
	var (
		lvl = TRACE
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Info : Utility for info log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Info
func Info(arg0 interface{}, args ...interface{}) {
	var (
		lvl = INFO
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Warn : Utility for warn log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Warn
func Warn(arg0 interface{}, args ...interface{}) error {
	var (
		lvl = WARNING
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
		return fmt.Errorf(first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLogf(lvl, "%s", str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
}

// Error : Utility for error log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Error
func Error(arg0 interface{}, args ...interface{}) error {
	var (
		lvl = ERROR
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
		return fmt.Errorf(first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLogf(lvl, "%s", str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
}

func Fatal(args ...interface{}) {
	Crash(args)
}

func Fatalf(format string, args ...interface{}) {
	Crashf(format, args)
}

func Fatalln(args ...interface{}) {
	Crash(args)
}

func Infoln(args ...interface{}) {
	var (
		lvl = INFO
	)
	Global.intLogf(lvl, strings.Repeat(" %v", len(args)), args...)
}

func Traceln(args ...interface{}) {
	var (
		lvl = TRACE
	)
	Global.intLogf(lvl, strings.Repeat(" %v", len(args)), args...)
}

func Recover() (paniced bool) {
	err := recover()
	if err != nil {
		buf := make([]byte, 10240)
		runtime.Stack(buf, false)
		Error("panic: %v,\n%s", err, string(buf))
		return true
	} else {
		return false
	}
}

func PrintStack(err interface{}) {
	buf := make([]byte, 10240)
	runtime.Stack(buf, false)
	Error("stack: %v,\n%s", err, string(buf))
}
