package common

import (
	"context"
	"encoding/json"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
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

	_, err = l.svcCtx.ResourceRpcClient.SyncResource(l.ctx, &borrower.SyncResourceRequest{
		Data: &borrower.SyncResourceRequestData{
			ResourceJson: string(resourceJson),
		},
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}

var resources = []dto.ResourceSyncItemDto{
	{
		ResourceType:    constants.ResourceTypeMenu,
		Name:            "home",
		Path:            "/home",
		Icon:            "Home",
		Tag:             "",
		IsShow:          true,
		Status:          constants.CommonStatusActive,
		EnableDataScope: false,
		Description:     "home menu",
		IsPublic:        true,
		Children: []dto.ResourceSyncItemDto{
			{
				ResourceType:    constants.ResourceTypeMenu,
				Name:            "dashboard",
				Path:            "/home/dashboard",
				Icon:            "Dashboard",
				Tag:             "",
				IsShow:          true,
				Status:          constants.CommonStatusActive,
				EnableDataScope: false,
				Description:     "dashboard menu",
				IsPublic:        true,
				Children:        nil,
				Apis: []dto.ApiSyncItemDto{
					{
						Path:   "/home/dashboard",
						Method: "GET",
					},
				},
			},
		},
		Apis: nil,
	},
}
