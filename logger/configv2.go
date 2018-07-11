package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/hydah/golib/logger/cfgstruct"
)

func (log Logger) LoadConfigurationConfV2(cfg cfgstruct.LogTypeSt) {
	typ := cfg.Type
	if len(typ) == 0 {
		return
	}
	typeList := strings.Split(typ, ",")
	if len(typeList) == 0 {
		return
	}
	log.LoadConfigurationDetailsV2(typeList, cfg)
}

func (log Logger) LoadConfigurationDetailsV2(typeList []string, cfg cfgstruct.LogTypeSt) (err error) {
	log.Close()
	for _, typ := range typeList {
		typ = strings.TrimSpace(typ)
		var filt, errFilt LogWriter
		var lvl level
		var enabled bool

		switch typ {
		case "stdout":
			filt, lvl, enabled, err = newStdoutLogWriterV2(cfg)
		case "file":
			filt, lvl, enabled, err = newFileLogWriterV2(cfg)
		default:
			err = fmt.Errorf("LoadConfiguration: Error: Required conf has unknown value of %s,%s\n", "type", typ)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		if !enabled {
			continue
		}

		log[typ] = &Filter{lvl, filt}

		if typ == "scribe" {
			log[typ+"-err"] = &Filter{WARNING, errFilt}
		}
	}
	return
}

func getLevel(levelStr string) (lvl level, err error) {
	switch levelStr {
	case "", "DEBUG":
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
		err = fmt.Errorf("LoadConfiguration: Error: Required conf has unknown value of %s,%s\n", "level", levelStr)
	}
	return
}

func newStdoutLogWriterV2(cfg cfgstruct.LogTypeSt) (w StdoutLogWriter, lvl level, enabled bool, err error) {
	if !cfg.StdoutEnabled {
		return
	}
	enabled = true
	lvl, err = getLevel(cfg.StdoutLevel)
	if err != nil {
		return
	}
	w = NewStdoutLogWriter()
	return
}

func newFileLogWriterV2(cfg cfgstruct.LogTypeSt) (w *FileLogWriter, lvl level, enabled bool, err error) {
	if !cfg.FileEnable {
		return
	}
	enabled = true
	lvl, err = getLevel(cfg.FileLevel)
	if err != nil {
		return
	}

	file := strings.Trim(cfg.FileFileName, " \r\n")
	if file == "" {
		p1, p2 := GetCurrPath()
		index := strings.LastIndex(p1, string(os.PathSeparator))
		path := p1[:index] + "/logs"
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

	format := strings.Trim(cfg.FileFormat, " \r\n")
	if len(format) == 0 {
		format = "[%D %T] [%L] (%S) %M"
	}

	w = NewFileLogWriter(file, !cfg.FileNoRotate)
	w.SetFormat(format)
	w.SetRotateLines(cfg.FileMaxLines)
	w.SetRotateSize(cfg.FileMaxSize)
	w.SetRotateDaily(!cfg.FileNoDaily)
	return
}
