package resourcelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
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
		baseLogic: logic.GetBaseLogic(svcCtx),
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

	data, err := copierx.Trans[[]*dto.ResourceTreeDto, []*admin.ResourceTreeEntity](&trees, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &admin.ResourceTreeResponse{
		Data: *data,
	}, nil
}
