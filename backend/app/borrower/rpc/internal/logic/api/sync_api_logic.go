package apilogic

import (
	"context"
	"encoding/json"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/swagger"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewSyncApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncApiLogic {
	return &SyncApiLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SyncApiLogic) SyncApi(in *borrower.SyncApiRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	var apiGroups swagger.ApiGroup
	if err = json.Unmarshal([]byte(in.Data.ApiJson), &apiGroups); err != nil {
		return nil, errorx.Wrap(err, "failed to unmarshal api json")
	}

	syncDto := &dto.ApiSyncRequestDto{
		Type: constants.UserTypeBorrower,
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

	return &borrower.EmptyResponse{}, nil
}
