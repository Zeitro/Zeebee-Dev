package loanlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/repo"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type LoanPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoanPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoanPageLogic {
	return &LoanPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoanPageLogic) LoanPage(in *lender.LoanPageRequest) (resp *lender.LoanPageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	pageRep := &repo.LoanLenderCompanyPageReq{}
	_ = copier.CopyWithOption(pageRep, in.Data, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
	})

	pageResult, err := l.svcCtx.LoanRepo.PageForLenderCompany(l.ctx, pageRep)
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
