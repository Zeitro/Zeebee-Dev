package loanlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"
	"zeebee-backend/common/errorx"
)

type GetBorrowerLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBorrowerLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBorrowerLoanLogic {
	return &GetBorrowerLoanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBorrowerLoanLogic) GetBorrowerLoan(in *borrower.GetLoanRequest) (resp *borrower.GetLoanResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	// todo: add your logic here and delete this line

	return &borrower.GetLoanResponse{}, nil
}
