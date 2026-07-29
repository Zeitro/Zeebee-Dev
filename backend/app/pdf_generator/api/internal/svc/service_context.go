package svc

import (
	"time"

	"zeebee-backend/app/pdf_generator/api/internal/global"
	"zeebee-backend/app/pdf_generator/api/internal/libreoffice"
	"zeebee-backend/app/pdf_generator/config"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config        config.ApiConf
	Logger        logx.Logger
	Converter     *libreoffice.Converter
	MaxInputBytes int
}

func NewServiceContext(c config.ApiConf, logger logx.Logger) *ServiceContext {
	maxInput := c.LibreOffice.MaxInputBytes
	if maxInput <= 0 {
		maxInput = 10 * 1024 * 1024
	}
	// Converter init failing (soffice absent) is non-fatal: the service stays up
	// and doc-to-pdf returns 500 so callers fail-soft. PDF-from-html is unaffected.
	converter, err := libreoffice.NewConverter(c.LibreOffice.ConcurrentCap,
		time.Duration(c.LibreOffice.TimeoutSec)*time.Second)
	if err != nil {
		logger.Errorf("init libreoffice converter failed (doc-to-pdf disabled): %v", err)
	}

	return &ServiceContext{
		Config:        c,
		Logger:        logger,
		Converter:     converter,
		MaxInputBytes: maxInput,
	}
}

// Close cleanup resources
// Should be called when service is stopping
func (s *ServiceContext) Close() {
	if s.Converter != nil {
		s.Converter.Close()
	}
	global.CloseChromePool()
	logx.Info("Service context closed")
}
