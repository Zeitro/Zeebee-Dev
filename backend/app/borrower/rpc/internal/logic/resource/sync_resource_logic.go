package resourcelogic

import (
	"context"
	"encoding/json"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
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
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SyncResourceLogic) SyncResource(in *borrower.SyncResourceRequest) (resp *proto.EmptyResponse, err error) {

	var resources []dto.ResourceSyncItemDto
	if err = json.Unmarshal([]byte(in.Data.ResourceJson), &resources); err != nil {
		return nil, errorx.Wrap(err, "failed to unmarshal resource json")
	}

	syncDto := &dto.ResourceSyncRequestDto{
		Type:      constants.UserTypeBorrower,
		Resources: resources,
	}

	err = l.baseLogic.ResourceLogic.SyncResource(l.ctx, syncDto)
	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}
