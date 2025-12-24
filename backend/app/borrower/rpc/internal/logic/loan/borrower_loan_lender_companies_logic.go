package loanlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type BorrowerLoanLenderCompaniesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewBorrowerLoanLenderCompaniesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerLoanLenderCompaniesLogic {
	return &BorrowerLoanLenderCompaniesLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *BorrowerLoanLenderCompaniesLogic) BorrowerLoanLenderCompanies(in *borrower.LoanLenderCompaniesRequest) (resp *borrower.LoanLenderCompaniesResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	var lenderCompanies []string
	lenderCompanies, err = l.svcCtx.LoanRepo.GetLenderCompanyForBorrower(l.ctx, claim.OrganizationUserId, claim.OrganizationId)
	if err != nil {
		return nil, errorx.Wrap(err, "Failed To Get Lender Companies")
	}

	return &borrower.LoanLenderCompaniesResponse{
		Data: lenderCompanies,
	}, nil
}
