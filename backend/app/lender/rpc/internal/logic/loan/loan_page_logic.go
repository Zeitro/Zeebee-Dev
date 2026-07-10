package loanlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
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

	baseLogic *logic.BaseLogic
}

func NewLoanPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoanPageLogic {
	return &LoanPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *LoanPageLogic) LoanPage(in *lender.LoanPageRequest) (resp *lender.LoanPageResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	pageRep := &dto.LoanLenderPageRequestDto{}
	pageRep, err = copierx.Trans[lender.LoanPageRequestData, dto.LoanLenderPageRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	pageRep.Filter.OrganizationId = claim.OrganizationId

	pageResult, err := l.svcCtx.LoanRepo.PageForLender(l.ctx, *pageRep)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	resp = &lender.LoanPageResponse{
		Data: &lender.LoanPageResponseData{},
	}

	resp.Data, err = copierx.Trans[gormx.PageResponse[*model.Loan], lender.LoanPageResponseData](&pageResult,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
