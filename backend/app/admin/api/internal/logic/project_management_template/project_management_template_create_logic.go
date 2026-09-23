package project_management_template

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementTemplateCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create
func NewProjectManagementTemplateCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTemplateCreateLogic {
	return &ProjectManagementTemplateCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTemplateCreateLogic) ProjectManagementTemplateCreate(req *types.ProjectManagementTemplateCreateRequest) (resp *types.ProjectManagementTemplateEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementTemplateCreateRequest, admin.ProjectManagementTemplateCreateRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ProjectManagementTemplateRpcClient.Create(l.ctx, &admin.ProjectManagementTemplateCreateRequest{Data: data})
	if err != nil {
		return nil, err
	}
	return copierx.Trans[proto.ProjectManagementTemplateEntity, types.ProjectManagementTemplateEntity](rpcResp.Data)
}
