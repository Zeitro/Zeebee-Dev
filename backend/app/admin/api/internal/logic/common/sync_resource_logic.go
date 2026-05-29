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
	var resourceList []dto.ResourceSyncItemDto
	err = json.Unmarshal(resourceBytes, &resourceList)
	if err != nil {
		return nil, errorx.Wrap(err, "failed to unmarshal resources")
	}
	if len(resourceList) == 0 {
		return nil, errorx.New("no resources found")
	}
	resourceContent, err := json.Marshal(resourceList)
	if err != nil {
		return nil, errorx.Wrap(err, "failed to marshal resource content")
	}

	_, err = l.svcCtx.ResourceRpcClient.SyncResource(l.ctx, &admin.SyncResourceRequest{
		Data: &admin.SyncResourceRequestData{
			ResourceJson: string(resourceContent),
		},
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}

var resourceBytes []byte = []byte(`[
  {
    "resource_type": "Menu",
    "name": "Test Menu",
    "path": "/TestMenu",
    "icon": "",
    "tag": "TestMenu",
    "is_show": true,
    "status": "Active",
    "enable_data_scope": false,
    "description": "test menu",
    "is_public": false,
    "children": [
      {
        "resource_type": "menu",
        "name": "sub test menu",
        "path": "/SubTestMenu",
        "icon": "",
        "tag": "TestMenu:Sub",
        "is_show": true,
        "status": "Active",
        "enable_data_scope": false,
        "description": "sub test menu",
        "is_public": false,
        "children": [
          {
            "resource_type": "Button",
            "name": "create",
            "path": "",
            "icon": "",
            "tag": "TestMenu:Sub:Create",
            "is_show": true,
            "status": "Active",
            "enable_data_scope": false,
            "description": "create button",
            "is_public": false,
            "apis": [
              {
                "path": "/api/v1/loan/:id",
                "method": "GET"
              },
              {
                "path": "/api/v1/loan",
                "method": "POST"
              }
            ]
          }
        ]
      }
    ],
    "apis": [

    ]
  }
]`)
