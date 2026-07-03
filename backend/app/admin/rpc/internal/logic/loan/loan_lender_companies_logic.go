package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type LoanLenderCompaniesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewLoanLenderCompaniesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoanLenderCompaniesLogic {
	return &LoanLenderCompaniesLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *LoanLenderCompaniesLogic) LoanLenderCompanies(in *proto.EmptyRequest) (resp *admin.LoanLenderCompaniesResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

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
