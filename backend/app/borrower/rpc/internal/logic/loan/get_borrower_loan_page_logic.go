package loanlogic

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type GetBorrowerLoanPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBorrowerLoanPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBorrowerLoanPageLogic {
	return &GetBorrowerLoanPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBorrowerLoanPageLogic) GetBorrowerLoanPage(in *borrower.PageLoanRequest) (resp *borrower.PageLoanResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	pageResult, err := l.svcCtx.LoanRepo.Page(l.ctx, in.Data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp = &borrower.PageLoanResponse{
		Data: &borrower.PageLoanResponseData{},
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
