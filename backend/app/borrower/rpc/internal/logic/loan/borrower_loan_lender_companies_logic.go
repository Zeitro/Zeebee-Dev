package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/errorx"
)

type BorrowerLoanLenderCompaniesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewBorrowerLoanLenderCompaniesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerLoanLenderCompaniesLogic {
	return &BorrowerLoanLenderCompaniesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *BorrowerLoanLenderCompaniesLogic) BorrowerLoanLenderCompanies(in *borrower.LoanLenderCompaniesRequest) (resp *borrower.LoanLenderCompaniesResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	var lenderCompanies []string
	lenderCompanies, err = l.svcCtx.LoanRepo.GetLenderCompanyForBorrower(l.ctx, claim.OrganizationUserId, claim.OrganizationId)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to get lender companies")
	}

	lenderCompanies = slice.Filter(lenderCompanies, func(index int, item string) bool {
		return len(item) > 0
	})

	return &borrower.LoanLenderCompaniesResponse{
		Data: lenderCompanies,
	}, nil
}
