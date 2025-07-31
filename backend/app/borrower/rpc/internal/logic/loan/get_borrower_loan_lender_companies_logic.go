package loanlogic

import (
	"context"

	"github.com/pkg/errors"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type GetBorrowerLoanLenderCompaniesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBorrowerLoanLenderCompaniesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBorrowerLoanLenderCompaniesLogic {
	return &GetBorrowerLoanLenderCompaniesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBorrowerLoanLenderCompaniesLogic) GetBorrowerLoanLenderCompanies(in *borrower.GetLoanLenderCompaniesRequest) (resp *borrower.GetLoanLenderCompaniesResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var lenderCompanies []string
	lenderCompanies, err = l.svcCtx.LoanRepo.GetLenderCompany(l.ctx, in.Data.BorrowerId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get lender companies")
	}

	return &borrower.GetLoanLenderCompaniesResponse{
		Data: lenderCompanies,
	}, nil
}
