package common

import (
	"context"
	"encoding/json"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// sync resource
func NewSyncResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncResourceLogic {
	return &SyncResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncResourceLogic) SyncResource(req *types.EmptyRequest) (resp *types.EmptyResponse, err error) {
	resourceJson, err := json.Marshal(resources)
	if err != nil {
		return nil, errorx.Wrap(err, "failed to marshal resources")
	}

	_, err = l.svcCtx.ResourceRpcClient.SyncResource(l.ctx, &admin.SyncResourceRequest{
		Data: &admin.SyncResourceRequestData{
			ResourceJson: string(resourceJson),
		},
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}

var resources = []dto.ResourceSyncItemDto{}
