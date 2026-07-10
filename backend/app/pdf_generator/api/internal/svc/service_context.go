package svc

import (
	"zeebee-backend/app/pdf_generator/api/internal/global"
	"zeebee-backend/app/pdf_generator/config"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config config.ApiConf
	Logger logx.Logger
}

func NewServiceContext(c config.ApiConf, logger logx.Logger) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Logger: logger,
	}
}

// Close cleanup resources
// Should be called when service is stopping
func (s *ServiceContext) Close() {
	global.CloseChromePool()
	logx.Info("Service context closed")
}
