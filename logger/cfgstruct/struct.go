package cfgstruct

type LogTypeSt struct {
	Type string `json:"type" ini:"type"`

	StdoutLevel   string `json:"stdout.level" ini:"stdout.level"`
	StdoutEnabled bool   `json:"stdout.enabled" ini:"stdout.enabled"`

	FileLevel    string `json:"file.level" ini:"file.level"`
	FileEnable   bool   `json:"file.enabled" ini:"file.enabled"`
	FileFileName string `json:"file.filename" ini:"file.filename"`
	FileFormat   string `json:"file.format,omitempty" ini:"file.format,omitempty"`
	FileMaxLines int    `json:"file.maxlines,omitempty" ini:"file.maxlines,omitempty"`
	FileMaxSize  int    `json:"file.maxsize,omitempty" ini:"file.maxsize,omitempty"`
	FileNoDaily  bool   `json:"file.nodaily,omitempty" ini:"file.nodaily,omitempty"`
	FileNoRotate bool   `json:"file.norotate,omitempty" ini:"file.norotate,omitempty"`

	ScribeLevel    string `json:"scribe.level" ini:"scribe.level"`
	ScribeEnabled  bool   `json:"scribe.enabled" ini:"scribe.enabled"`
	ScribeEndpoint string `json:"scribe.endpoint" ini:"scribe.endpoint"`
	ScribeCategory string `json:"scribe.category" ini:"scribe.category"`
	ScribeFormat   string `json:"scribe.format" ini:"scribe.format"`

	CLogLevel   string `json:"clog.level" ini:"clog.level"`
	CLogEnabled bool   `json:"clog.enabled" ini:"clog.enabled"`
	CLogHost    string `json:"clog.host" ini:"clog.host"`
	CLogModule  string `json:"clog.module" ini:"clog.module"`
}
