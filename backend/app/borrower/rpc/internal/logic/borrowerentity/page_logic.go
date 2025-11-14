package borrowerentitylogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/repo"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
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

func (l *PageLogic) Page(in *borrower.BorrowerEntityPageRequest) (resp *borrower.BorrowerEntityPageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	pageReq := &repo.BorrowerEntityBorrowerPageReq{}
	_ = copier.CopyWithOption(pageReq, in.Data, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
	})
	pageResult, err := l.svcCtx.BorrowerEntityRepo.PageForBorrower(l.ctx, pageReq)
	if err != nil {
		return nil, err
	}

	respData := make([]*borrower.BorrowerEntityPageEntry, 0, len(pageResult.List))

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

	return &borrower.BorrowerEntityPageResponse{
		Data: &borrower.BorrowerEntityPageResponseData{
			Page:      pageResult.Page,
			Limit:     pageResult.Limit,
			Total:     pageResult.Total,
			TotalPage: pageResult.TotalPage,
			List:      respData,
		},
	}, nil
}
