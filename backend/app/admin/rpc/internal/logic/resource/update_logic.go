package resourcelogic

import (
	"context"

	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLogic) Update(in *admin.ResourceUpdateRequest) (resp *admin.ResourceCreateResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// First get the existing record
	result, err := l.svcCtx.ResourceRepo.GetResource(l.ctx, in.Data.Id, constants.UserType(in.Data.Type))
	if err != nil {
		return nil, err
	}

	// Update the fields
	err = copier.CopyWithOption(result, in.Data, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	// check existed
	isExisted, err := l.svcCtx.ResourceRepo.IsExisted(l.ctx, result)
	if err != nil {
		return nil, err
	}
	if isExisted {
		err = errorx.New("Resource is Existed")
		return nil, err
	}

	// Save the updated record
	_, err = l.svcCtx.ResourceRepo.UpdateOneById(l.ctx, in.Data.Id, result)
	if err != nil {
		return nil, err
	}

	respData := &admin.ResourceTreeEntity{}
	err = copier.CopyWithOption(respData, result, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return &admin.ResourceCreateResponse{
		Data: respData,
	}, nil
}
