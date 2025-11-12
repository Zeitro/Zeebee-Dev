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

type PageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PageLogic) Page(in *borrower.BorrowerApplicantPageRequest) (resp *borrower.BorrowerApplicantPageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	pageResult, err := l.svcCtx.BorrowerApplicantRepo.Page(l.ctx, in.Data)
	if err != nil {
		return nil, err
	}

	respData := make([]*borrower.BorrowerApplicantPageEntry, 0, len(pageResult.List))

	err = copier.CopyWithOption(&respData, &pageResult.List, copier.Option{
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

	return &borrower.BorrowerApplicantPageResponse{
		Data: &borrower.BorrowerApplicantPageResponseData{
			Page:      pageResult.Page,
			Limit:     pageResult.Limit,
			Total:     pageResult.Total,
			TotalPage: pageResult.TotalPage,
			List:      respData,
		},
	}, nil
}
