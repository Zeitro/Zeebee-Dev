package resourcelogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResourceTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewResourceTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceTreeLogic {
	return &ResourceTreeLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *ResourceTreeLogic) ResourceTree(in *admin.ResourceTreeRequest) (resp *admin.ResourceTreeResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	trees, err := l.baseLogic.ResourceLogic.ResourceTree(l.ctx, constants.UserType(in.Data.Type))
	if err != nil {
		return nil, err
	}

	data := []*admin.ResourceTreeEntity{}
	err = copier.CopyWithOption(&data, &trees, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &admin.ResourceTreeResponse{
		Data: data,
	}, nil
}
