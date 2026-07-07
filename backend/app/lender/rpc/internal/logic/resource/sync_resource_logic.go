package resourcelogic

import (
	"context"
	"encoding/json"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewSyncResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncResourceLogic {
	return &SyncResourceLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SyncResourceLogic) SyncResource(in *lender.SyncResourceRequest) (resp *proto.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	var resources []dto.ResourceSyncItemDto
	if err = json.Unmarshal([]byte(in.Data.ResourceJson), &resources); err != nil {
		return nil, errorx.Wrap(err, "failed to unmarshal resource json")
	}

	syncDto := &dto.ResourceSyncRequestDto{
		Type:      constants.UserTypeLender,
		Resources: resources,
	}

	err = l.baseLogic.ResourceLogic.SyncResource(l.ctx, syncDto)
	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}
