package resourcelogic

import (
	"context"
	"encoding/json"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type SyncResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewSyncResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncResourceLogic {
	return &SyncResourceLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// SyncResource upserts a backend manifest into the resource registry.
func (l *SyncResourceLogic) SyncResource(in *tenant.SyncResourceRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}

	var resources []dto.ResourceSyncItemDto
	if err = json.Unmarshal([]byte(in.Data.ResourceJson), &resources); err != nil {
		return nil, errorx.Wrap(err, "failed to unmarshal resource json")
	}

	err = l.baseLogic.ResourceLogic.SyncResource(l.ctx, &dto.ResourceSyncRequestDto{
		Type:      userType,
		Resources: resources,
	})
	if err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
