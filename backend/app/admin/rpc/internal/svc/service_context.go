package svc

import (
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/config"
	"zeebee-backend/app/admin/rpc/internal/global"
	"zeebee-backend/common/llm"
	"zeebee-backend/domains/zeebee/repo"
)

type ServiceContext struct {
	Logger logx.Logger
	Config config.RpcConf
	*repo.AllRepos
	// DiscoverLLM is the lender-discover LLM client (built from LlmConf).
	// nil if the provider is unconfigured; discover calls then fail at runtime.
	DiscoverLLM llm.Client
}

func NewServiceContext(c config.RpcConf, logger logx.Logger) *ServiceContext {
	discoverLLM, err := llm.NewClient(c.LlmConf.Providers, c.LlmConf.LenderDiscoverConf)
	if err != nil {
		logger.Errorf("lender discover llm client not initialized: %v", err)
	}
	return &ServiceContext{
		Config:      c,
		Logger:      logger,
		AllRepos:    repo.NewAllRepos(global.GetGormDB()),
		DiscoverLLM: discoverLLM,
	}
}
