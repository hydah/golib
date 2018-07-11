package logger

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hydah/golib/config"
)

// Parse a number with K/M/G suffixes based on thousands (1000) or 2^10 (1024)
func strToNumSuffix(str string, mult int) int {
	num := 1
	if len(str) > 1 {
		switch str[len(str)-1] {
		case 'G', 'g':
			num *= mult
			fallthrough
		case 'M', 'm':
			num *= mult
			fallthrough
		case 'K', 'k':
			num *= mult
			str = str[0 : len(str)-1]
		}
	}
	parsed, _ := strconv.Atoi(str)
	return parsed * num
}
func GetCurrPath() (string, string) {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	p1 := path[:index]
	p2 := path[index:]
	return p1, p2
}
func (log Logger) LoadConfigurationConf(cfg *config.Config) {
	sec, err := cfg.GetSection("logtype")
	if err != nil {
		return
	}
	typ, err := sec.GetValue("type")
	if err != nil || len(typ) == 0 {
		return
	}
	typeList := strings.Split(typ, ",")
	if len(typeList) == 0 {
		return
	}
	log.LoadConfigurationDetails(typeList, sec)
}
func (log Logger) LoadConfigurationDetails(typeList []string, sec config.Section) (err error) {
	log.Close()
	for _, typ := range typeList {
		typ = strings.TrimSpace(typ)
		var filt, errFilt LogWriter
		var lvl level
		var debugLevel string
		bad, good := false, true
		debugLevel, err = sec.GetValue(typ + ".level")
		if err != nil {
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required conf failed")
		}
		if len(strings.TrimSpace(debugLevel)) == 0 {
			debugLevel = "DEBUG"
		}
		switch debugLevel {
		case "DEBUG":
			lvl = DEBUG
		case "TRACE":
			lvl = TRACE
		case "INFO":
			lvl = INFO
		case "WARNING":
			lvl = WARNING
		case "ERROR":
			lvl = ERROR
		default:
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required conf has unknown value of %s,%s\n", "level", debugLevel)
			bad = true
		}
		// Just so all of the required attributes are errored at the same time if missing
		if bad {
			os.Exit(1)
		}
		switch typ {
		case "stdout":
			filt, good = newStdoutLogWriter(typ, sec)
		case "file":
			filt, good = newFileLogWriter(typ, sec)
		default:
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required conf has unknown value of %s,%s\n", "level", typ)
			os.Exit(1)
		}
		// Just so all of the required params are errored at the same time if wrong
		if !good {
			os.Exit(1)
		}
		// If we're disabled (syntax and correctness checks only), don't add to logger
		enabled, err := sec.GetValue(typ + ".enabled")
		if err != nil {
			os.Exit(1)
		}
		enable, err := strconv.ParseBool(enabled)
		if !enable || err != nil {
			continue
		}
		log[typ] = &Filter{lvl, filt}
		if typ == "scribe" {
			log[typ+"-err"] = &Filter{WARNING, errFilt}
		}
	}
	return
}
func newStdoutLogWriter(prefix string, sec config.Section) (StdoutLogWriter, bool) {
	enabled, err := sec.GetValue(prefix + ".enabled")
	if err != nil {
		return nil, true
	}
	enable, err := strconv.ParseBool(enabled)
	if !enable || err != nil {
		return nil, true
	}
	return NewStdoutLogWriter(), true
}
func newFileLogWriter(prefix string, sec config.Section) (*FileLogWriter, bool) {
	// If it's disabled, we're just checking syntax
	enabled, err := sec.GetValue(prefix + ".enabled")
	if err != nil {
		return nil, true
	}
	enable, err := strconv.ParseBool(enabled)
	if !enable || err != nil {
		return nil, true
	}
	for key, val := range sec {
		_ = key
		_ = val
		//fmt.Println(key, val)
	}
	var file string
	format := "[%D %T] [%L] (%S) %M"
	maxlines := 0
	maxsize := 0
	daily := true
	rotate := true
	// Parse properties
	for key, val := range sec {
		switch key {
		case prefix + ".filename":
			file = strings.Trim(val, " \r\n")
		case prefix + ".format":
			format = strings.Trim(val, " \r\n")
		case prefix + ".maxlines":
			maxlines = strToNumSuffix(strings.Trim(val, " \r\n"), 1000)
		case prefix + ".maxsize":
			maxsize = strToNumSuffix(strings.Trim(val, " \r\n"), 1024)
		case prefix + ".daily":
			daily = strings.Trim(val, " \r\n") != "false"
		case prefix + ".rotate":
			rotate = strings.Trim(val, " \r\n") != "false"
		default:
		}
	}
	if file == "" {
		p1, p2 := GetCurrPath()
		index := strings.LastIndex(p1, string(os.PathSeparator))
		path := p1[:index] + "/log"
		os.MkdirAll(path, 0777)
		file = path + p2
	} else {
		index := strings.LastIndex(file, string(os.PathSeparator))
		if index >= 0 {
			path := file[:index] + "/"
			if _, err := os.Stat(path); err != nil {
				if os.IsNotExist(err) {
					os.MkdirAll(path, 0777)
				}
			}
		}
	}
	flw := NewFileLogWriter(file, rotate)
	flw.SetFormat(format)
	flw.SetRotateLines(maxlines)
	flw.SetRotateSize(maxsize)
	flw.SetRotateDaily(daily)
	return flw, true
}