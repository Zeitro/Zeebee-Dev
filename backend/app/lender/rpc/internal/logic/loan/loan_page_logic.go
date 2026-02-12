package loanlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/data/dto"

	"zeebee-backend/common/errorx"
)

type LoanPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewLoanPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoanPageLogic {
	return &LoanPageLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *LoanPageLogic) LoanPage(in *lender.LoanPageRequest) (resp *lender.LoanPageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	pageRep := dto.LoanLenderPageRequestDto{}
	_ = copier.CopyWithOption(&pageRep, in.Data, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
	})
	pageRep.Filter.OrganizationId = claim.OrganizationId

	pageResult, err := l.svcCtx.LoanRepo.PageForLender(l.ctx, pageRep)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	resp = &lender.LoanPageResponse{
		Data: &lender.LoanPageResponseData{},
	}

	err = copier.CopyWithOption(&resp.Data, pageResult, copier.Option{
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
			copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
