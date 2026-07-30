package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	ApiConf ApiConf
}

type ApiConf struct {
	rest.RestConf
	// LibreOffice drives doc/docx → pdf conversion. Optional: if soffice is
	// absent the service still starts (doc-to-pdf returns 500, callers fail-soft).
	LibreOffice LibreOfficeConf `json:",optional"`
}

// LibreOfficeConf tunes the headless doc → pdf converter.
type LibreOfficeConf struct {
	ConcurrentCap int `json:",default=4"`        // profile pool size + concurrency gate
	TimeoutSec    int `json:",default=90"`       // per-file soffice ctx timeout
	MaxInputBytes int `json:",default=10485760"` // input cap, 10MB default
}
