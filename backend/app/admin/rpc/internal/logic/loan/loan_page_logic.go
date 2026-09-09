package loanlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"

	"zeebee-backend/common/errorx"
)

type LoanPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoanPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoanPageLogic {
	return &LoanPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoanPageLogic) LoanPage(in *admin.LoanPageRequest) (resp *admin.LoanPageResponse, err error) {

	pageRep, err := copierx.Trans[admin.LoanPageRequestData, dto.LoanAdminPageRequestDto](in.Data)
	if err != nil {
		return nil, err
	}

	pageResult, err := l.svcCtx.LoanRepo.PageForAdmin(l.ctx, *pageRep)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	data, err := copierx.Trans[gormx.PageResponse[*model.Loan], admin.LoanPageResponseData](&pageResult,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return &admin.LoanPageResponse{
		Data: data,
	}, nil
}
