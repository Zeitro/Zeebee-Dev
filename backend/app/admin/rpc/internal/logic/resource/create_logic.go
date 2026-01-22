package resourcelogic

import (
	"context"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

func (l *CreateLogic) Create(in *admin.ResourceCreateRequest) (resp *admin.ResourceDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	resourceModel := &model.Resource{}
	err = copier.CopyWithOption(resourceModel, in.Data, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	var result *model.Resource

	err = l.svcCtx.ResourceRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		var isExisted bool
		// check existed
		isExisted, innerErr = l.svcCtx.ResourceRepo.IsExistedWithTx(l.ctx, tx, resourceModel)
		if innerErr != nil {
			return innerErr
		}
		if isExisted {
			innerErr = errorx.New("Resource is Existed")
			return innerErr
		}

		apis, innerErr := l.svcCtx.ApiRepo.GetApisWithTx(l.ctx, tx, in.Data.ApiIds, resourceModel.Type)
		if innerErr != nil {
			return innerErr
		}
		resourceModel.Apis = apis

		result, innerErr = l.svcCtx.ResourceRepo.CreateWithTx(l.ctx, tx, resourceModel, true)

		return innerErr
	})

	// result, err := l.svcCtx.ResourceRepo.Create(l.ctx, resourceModel)
	if err != nil {
		return nil, err
	}

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
