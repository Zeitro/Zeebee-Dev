package rolelogic

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/repo"
)

type PageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *PageLogic) Page(in *borrower.RolePageRequest) (resp *borrower.RolePageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	pageReq := &repo.RolePageReq{}
	_ = copier.CopyWithOption(pageReq, in.Data, copier.Option{
		CaseSensitive: true,
		DeepCopy:      true,
	})
	if pageReq.Filter == nil {
		pageReq.Filter = &repo.RolePageReqFilter{}
	}
	pageReq.Filter.RoleType = constants.UserTypeBorrower

	result, err := l.baseLogic.RoleLogic.Page(l.ctx, pageReq)
	if err != nil {
		return nil, err
	}

	data := &borrower.RolePageResponseData{}
	err = copier.CopyWithOption(data, &result, copier.Option{
		DeepCopy: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &borrower.RolePageResponse{
		Data: data,
	}, nil
}
