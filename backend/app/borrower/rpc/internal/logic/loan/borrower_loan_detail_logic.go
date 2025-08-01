package loanlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type BorrowerLoanDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBorrowerLoanDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerLoanDetailLogic {
	return &BorrowerLoanDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BorrowerLoanDetailLogic) BorrowerLoanDetail(in *borrower.LoanDetailRequest) (resp *borrower.LoanDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	loan, err := l.svcCtx.LoanRepo.GetBorrowerLoan(l.ctx, in.Data.Id, in.Data.BorrowerId)
	if err != nil {
		return nil, err
	}
	resp = &borrower.LoanDetailResponse{
		Data: &borrower.LoanEntry{},
	}
	err = copier.CopyWithOption(&resp.Data, loan, copier.Option{
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
			copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
		},
	})

	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, err
}
