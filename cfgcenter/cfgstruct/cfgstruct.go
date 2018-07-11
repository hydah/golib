package cfgstruct

import (
	loggercfg "github.com/hydah/golib/logger/cfgstruct"
)

type BaseCfgSt struct {
	LogConf     LogConfSt                    `json:"log_conf" ini:"log_conf"`
	LogType     loggercfg.LogTypeSt          `json:"logtype" ini:"logtype"`
}

type LogConfSt struct {
	LogDir string `json:"logdir" ini:"logdir"`
	Prefix string `json:"prefix" ini:"prefix"`
}