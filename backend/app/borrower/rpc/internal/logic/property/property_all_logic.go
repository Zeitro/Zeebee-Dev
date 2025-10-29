package propertylogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/timex"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyAllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPropertyAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyAllLogic {
	return &PropertyAllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PropertyAllLogic) PropertyAll(in *borrower.PropertyAllRequest) (resp *borrower.PropertyAllResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	results, err := l.svcCtx.PropertyBaseRepo.All(l.ctx, in.BorrowerId)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	respData := make([]*borrower.PropertyBaseEntry, 0, len(results))

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

	return &borrower.PropertyAllResponse{
		Data: respData,
	}, nil
}
