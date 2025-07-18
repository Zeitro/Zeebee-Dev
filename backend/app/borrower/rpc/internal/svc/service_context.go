package svc

import (
	"zeebee-backend/app/borrower/rpc/internal/config"
	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/repo"
)

type ServiceContext struct {
	Config config.Config

	// repo
	UserRepo     *repo.UserRepo
	BorrowerRepo *repo.BorrowerRepo
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		// repo
		UserRepo:     repo.NewUserRepo(global.GetGormDB()),
		BorrowerRepo: repo.NewBorrowerRepo(global.GetGormDB()),
	}
}
