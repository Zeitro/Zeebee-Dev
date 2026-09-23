package project_management_template

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementTemplateDefaultsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// project defaults
func NewProjectManagementTemplateDefaultsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTemplateDefaultsLogic {
	return &ProjectManagementTemplateDefaultsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTemplateDefaultsLogic) ProjectManagementTemplateDefaults(req *types.ProjectManagementTemplateDefaultsRequest) (resp *types.EmptyResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementTemplateDefaultsRequest, admin.ProjectManagementTemplateDefaultsRequestData](req)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.ProjectManagementTemplateRpcClient.Defaults(l.ctx, &admin.ProjectManagementTemplateDefaultsRequest{Data: data})
	if err != nil {
		return nil, err
	}
	return &types.EmptyResponse{}, nil
}
