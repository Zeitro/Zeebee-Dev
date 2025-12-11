package apilogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewApiTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiTreeLogic {
	return &ApiTreeLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *ApiTreeLogic) ApiTree(in *admin.ApiTreeRequest) (resp *admin.ApiTreeResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	trees, err := l.baseLogic.ApiLogic.ApiTree(l.ctx, constants.UserType(in.Data.Type))
	if err != nil {
		return nil, err
	}
	data := []*admin.ApiTreeEntity{}
	err = copier.CopyWithOption(&data, &trees, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &admin.ApiTreeResponse{
		Data: data,
	}, nil
}
