package svc

import (
	"net/http"
	"time"

	"zeebee-backend/app/borrower/config"
	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/domains/zeebee/repo"
	"zeebee-backend/domains/zeebee/zeitro"
)

type ServiceContext struct {
	Config config.RpcConf
	*repo.AllRepos
	ZeitroClient *zeitro.Client
}

func NewServiceContext(c config.RpcConf) *ServiceContext {
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

	return &ServiceContext{
		Config:       c,
		AllRepos:     repo.NewAllRepos(global.GetGormDB()),
		ZeitroClient: zeitroClient,
	}
}
