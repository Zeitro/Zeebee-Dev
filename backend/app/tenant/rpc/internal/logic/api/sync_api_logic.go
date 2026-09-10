package apilogic

import (
	"context"
	"encoding/json"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/swagger"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type SyncApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewSyncApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncApiLogic {
	return &SyncApiLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// SyncApi upserts the caller platform's api registry from a swagger manifest.
func (l *SyncApiLogic) SyncApi(in *tenant.SyncApiRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}

	var apiGroups swagger.ApiGroup
	if err = json.Unmarshal([]byte(in.Data.ApiJson), &apiGroups); err != nil {
		return nil, errorx.Wrap(err, "failed to unmarshal api json")
	}

	syncDto := &dto.ApiSyncRequestDto{
		Type: userType,
		Apis: make(map[string][]dto.ApiSyncItemDto),
	}
	for tag, apis := range apiGroups {
		items := make([]dto.ApiSyncItemDto, 0, len(apis))
		for _, api := range apis {
			items = append(items, dto.ApiSyncItemDto{
				Path:        api.Path,
				Method:      api.Method,
				Name:        api.Name,
				Description: api.Description,
			})
		}
		syncDto.Apis[tag] = items
	}

	if err = l.baseLogic.ApiLogic.SyncApi(l.ctx, syncDto); err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
