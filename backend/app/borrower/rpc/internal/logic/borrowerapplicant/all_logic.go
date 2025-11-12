package borrowerapplicantlogic

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"
)

type AllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLogic {
	return &AllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AllLogic) All(in *borrower.BorrowerApplicantAllRequest) (resp *borrower.BorrowerApplicantAllResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	results, err := l.svcCtx.BorrowerApplicantRepo.All(l.ctx, in.Data.BorrowerId)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	respData := make([]*borrower.BorrowerApplicantPageEntry, 0, len(results))

	err = copier.CopyWithOption(&respData, &results, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
			copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &borrower.BorrowerApplicantAllResponse{
		Data: respData,
	}, nil
}
