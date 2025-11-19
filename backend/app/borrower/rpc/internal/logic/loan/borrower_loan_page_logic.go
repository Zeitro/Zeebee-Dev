package loanlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/repo"

	"github.com/zeromicro/go-zero/core/logx"
)

type BorrowerLoanPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBorrowerLoanPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerLoanPageLogic {
	return &BorrowerLoanPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BorrowerLoanPageLogic) BorrowerLoanPage(in *borrower.LoanPageRequest) (resp *borrower.LoanPageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	pageReq := &repo.LoanBorrowerPageReq{}
	_ = copier.CopyWithOption(pageReq, in.Data, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
	})

	pageResult, err := l.svcCtx.LoanRepo.PageForBorrower(l.ctx, pageReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	resp = &borrower.LoanPageResponse{
		Data: &borrower.LoanPageResponseData{},
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
