package loanlogic

import (
	"context"

	"github.com/pkg/errors"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type BorrowerLoanLenderCompaniesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBorrowerLoanLenderCompaniesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerLoanLenderCompaniesLogic {
	return &BorrowerLoanLenderCompaniesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BorrowerLoanLenderCompaniesLogic) BorrowerLoanLenderCompanies(in *borrower.LoanLenderCompaniesRequest) (resp *borrower.LoanLenderCompaniesResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	var lenderCompanies []string
	lenderCompanies, err = l.svcCtx.LoanRepo.GetLenderCompany(l.ctx, in.Data.BorrowerId)
	if err != nil {
		return nil, errors.Wrap(err, "Failed To Get Lender Companies")
	}

	return &borrower.LoanLenderCompaniesResponse{
		Data: lenderCompanies,
	}, nil
}
