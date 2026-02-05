package resourcelogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *UpdateLogic) Update(in *admin.ResourceUpdateRequest) (resp *admin.ResourceDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	var result *model.Resource
	err = l.svcCtx.ResourceRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		var isExisted bool

		result, innerErr = l.svcCtx.ResourceRepo.GetResourceWithTx(l.ctx, tx, in.Data.Id, constants.UserType(in.Data.Type))
		if innerErr != nil {
			return innerErr
		}
		innerErr = copier.CopyWithOption(result, in.Data, copier.Option{DeepCopy: true})
		if innerErr != nil {
			innerErr = errorx.GetFailedTransRespErr(innerErr)
			return innerErr
		}

		// check existed
		isExisted, innerErr = l.svcCtx.ResourceRepo.IsExistedWithTx(l.ctx, tx, result)
		if innerErr != nil {
			return innerErr
		}
		if isExisted {
			innerErr = errorx.New("Resource is existed")
			return innerErr
		}
		// handle apis
		oldApis := result.Apis
		currentApis, innerErr := l.svcCtx.ApiRepo.GetApisWithTx(l.ctx, tx, in.Data.ApiIds, constants.UserType(in.Data.Type))
		if innerErr != nil {
			return innerErr
		}
		removedApis := slice.DifferenceWith(oldApis, currentApis, func(item1, item2 *model.Api) bool {
			return item1.Id == item2.Id
		})
		newApis := slice.DifferenceWith(currentApis, oldApis, func(item1, item2 *model.Api) bool {
			return item1.Id == item2.Id
		})
		result.RemoveRelationships()
		innerErr = l.svcCtx.ResourceRepo.RemoveApisWithTx(l.ctx, tx, result, removedApis)
		if innerErr != nil {
			return innerErr
		}
		innerErr = l.svcCtx.ResourceRepo.AddApisWithTx(l.ctx, tx, result, newApis)
		if innerErr != nil {
			return innerErr
		}
		// todo update user permissions

		return innerErr
	})

	respData := &admin.ResourceEntity{}
	err = copier.CopyWithOption(respData, result, copier.Option{
		DeepCopy: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &admin.ResourceDetailResponse{
		Data: respData,
	}, nil
}
