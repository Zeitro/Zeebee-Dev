package apilogic

import (
	"context"
	"encoding/json"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/swagger"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type SyncApiLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewSyncApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncApiLogic {
	return &SyncApiLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SyncApiLogic) SyncApi(in *borrower.SyncApiRequest) (resp *proto.EmptyResponse, err error) {

	var apiGroups swagger.ApiGroup
	if err = json.Unmarshal([]byte(in.Data.ApiJson), &apiGroups); err != nil {
		return nil, errorx.Wrap(err, "failed to unmarshal api json")
	}

	syncDto := &dto.ApiSyncRequestDto{
		Type: constants.UserTypeAdmin,
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

	err = l.baseLogic.ApiLogic.SyncApi(l.ctx, syncDto)
	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}
