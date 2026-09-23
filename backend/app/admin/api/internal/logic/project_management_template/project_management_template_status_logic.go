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

type ProjectManagementTemplateStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// status
func NewProjectManagementTemplateStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTemplateStatusLogic {
	return &ProjectManagementTemplateStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTemplateStatusLogic) ProjectManagementTemplateStatus(req *types.ProjectManagementTemplateStatusRequest) (resp *types.ProjectManagementTemplateEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementTemplateStatusRequest, admin.ProjectManagementTemplateStatusRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ProjectManagementTemplateRpcClient.Status(l.ctx, &admin.ProjectManagementTemplateStatusRequest{Data: data})
	if err != nil {
		return nil, err
	}
	return copierx.Trans[proto.ProjectManagementTemplateEntity, types.ProjectManagementTemplateEntity](rpcResp.Data)
}
