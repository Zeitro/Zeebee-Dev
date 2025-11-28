package resourcelogic

import (
	"context"

	"github.com/jinzhu/copier"

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
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *admin.ResourceCreateRequest) (resp *admin.ResourceCreateResponse, err error) {
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

	// check existed
	isExisted, err := l.svcCtx.ResourceRepo.IsExisted(l.ctx, resourceModel)
	if err != nil {
		return nil, err
	}
	if isExisted {
		err = errorx.New("Resource is Existed")
		return nil, err
	}

	result, err := l.svcCtx.ResourceRepo.Create(l.ctx, resourceModel)
	if err != nil {
		return nil, err
	}

	respData := &admin.ResourceTreeEntity{}
	err = copier.CopyWithOption(respData, result, copier.Option{
		DeepCopy: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &admin.ResourceCreateResponse{
		Data: respData,
	}, nil
}
