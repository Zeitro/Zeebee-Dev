package commonlogic

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/static"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type SyncDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewSyncDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncDataLogic {
	return &SyncDataLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SyncDataLogic) SyncData(in *admin.EmptyRequest) (resp *admin.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	err = l.syncApi()
	if err != nil {
		return nil, err
	}

	return &admin.EmptyResponse{}, nil
}

type ApiGroup struct {
	Group string     `json:"group"`
	Apis  []*ApiItem `json:"apis"`
}

type ApiItem struct {
	Name               string `json:"name"`
	Path               string `json:"path"`
	Method             string `json:"method"`
	Description        string `json:"description"`
	AuthRequired       bool   `json:"auth_required"`
	PermissionRequired bool   `json:"permission_required"`
}

func (l *SyncDataLogic) syncApi() (err error) {
	apiMap := map[constants.UserType][]ApiGroup{}
	apiContent := static.ApiStatic
	err = json.Unmarshal(apiContent, &apiMap)
	if err != nil {
		return errorx.Wrap(err, "Failed to unmarshal api content")
	}
	apiRepo := l.svcCtx.ApiRepo
	for apiType, apiGroup := range apiMap {
		err = apiRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
			var innerError error
			var parentApiSort int32 = 1
			for _, group := range apiGroup {
				// get or create parent api
				parentApi := &model.Api{
					Name:        group.Group,
					Path:        "",
					Method:      "",
					Description: "",
					Type:        apiType,
					ParentId:    0,
					SortOrder:   parentApiSort,
				}
				parentApi, innerError = apiRepo.GetOrCreateApiWithTx(l.ctx, tx, parentApi)
				if innerError != nil {
					return innerError
				}
				// get or create child api
				childApiSort := int32(1)
				for _, api := range group.Apis {
					childApi := &model.Api{
						Name:               api.Name,
						Path:               api.Path,
						Method:             api.Method,
						Description:        api.Description,
						Type:               apiType,
						ParentId:           parentApi.Id,
						AuthRequired:       api.AuthRequired,
						PermissionRequired: api.PermissionRequired,
						SortOrder:          childApiSort,
					}
					_, innerError = apiRepo.GetOrCreateApiWithTx(l.ctx, tx, childApi)
					if innerError != nil {
						return innerError
					}
					childApiSort++
				}
				parentApiSort++
			}

			return innerError
		})
		if err != nil {
			return err
		}
	}

	return err
}
