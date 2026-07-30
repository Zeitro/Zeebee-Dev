package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type LoanLenderCompaniesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewLoanLenderCompaniesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoanLenderCompaniesLogic {
	return &LoanLenderCompaniesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *LoanLenderCompaniesLogic) LoanLenderCompanies(in *proto.EmptyRequest) (resp *admin.LoanLenderCompaniesResponse, err error) {

	var lenderCompanies []string
	lenderCompanies, err = l.svcCtx.LoanRepo.GetLenderCompanyForAdmin(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to get lender companies")
	}

	lenderCompanies = slice.Filter(lenderCompanies, func(index int, item string) bool {
		return len(item) > 0
	})

	return &admin.LoanLenderCompaniesResponse{
		Data: lenderCompanies,
	}, nil
}
