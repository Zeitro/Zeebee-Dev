package svc

import (
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/config"
	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file"
	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/common/llm"
	"zeebee-backend/domains/zeebee/repo"
	"zeebee-backend/domains/zeebee/zeitro"
)

type ServiceContext struct {
	Config config.RpcConf
	Logger logx.Logger
	*repo.AllRepos
	ZeitroClient *zeitro.Client

	// LLM is the shared single-instance client for file-parse extraction +
	// arbitration. Nil if unconfigured (parse dispatch becomes a no-op).
	LLM      llm.Client
	// Progress owns the bpt:parse:* Redis keys (locks, staged progress, done,
	// reqmap, queue, job) for async parse + restart recovery.
	Progress *parse_file.ProgressManager
	// FileParser orchestrates async file extraction (escrow OnUpload now; rpa
	// manual_tree/select with US2/US3). Nil-shared deps make it a no-op when
	// LLM is unconfigured.
	FileParser *parse_file.Parser
}

func NewServiceContext(c config.RpcConf, logger logx.Logger) *ServiceContext {
	// Create HTTP client for Zeitro API calls with 10 second timeout
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	// Create Zeitro client
	zeitroClient, err := zeitro.NewClient(zeitro.Config{
		URL:      c.ZeitroConf.Url,
		Username: c.ZeitroConf.Username,
		Password: c.ZeitroConf.Password,
	}, httpClient)
	if err != nil {
		panic(err)
	}

	// Shared LLM client (single instance) for file-parse extraction. Optional:
	// if unconfigured the parse dispatcher no-ops rather than blocking startup.
	var llmClient llm.Client
	if lc, lerr := llm.NewClient(c.LlmConf.Providers, c.LlmConf.TransactionFileParseConf); lerr != nil {
		logger.Errorf("init llm client for file parse failed (parse disabled): %v", lerr)
	} else {
		llmClient = lc
	}

	ctx := &ServiceContext{
		Config:       c,
		Logger:       logger,
		AllRepos:     repo.NewAllRepos(global.GetGormDB()),
		ZeitroClient: zeitroClient,
		LLM:          llmClient,
		Progress:     parse_file.NewProgressManager(global.GetRedis(), logger),
	}
	ctx.FileParser = parse_file.NewParser(ctx.AllRepos, llmClient, ctx.Progress, logger)
	return ctx
}
