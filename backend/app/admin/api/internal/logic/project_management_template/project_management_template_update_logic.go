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

type ProjectManagementTemplateUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update
func NewProjectManagementTemplateUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTemplateUpdateLogic {
	return &ProjectManagementTemplateUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTemplateUpdateLogic) ProjectManagementTemplateUpdate(req *types.ProjectManagementTemplateUpdateRequest) (resp *types.ProjectManagementTemplateEntity, err error) {
	data, err := copierx.Trans[types.ProjectManagementTemplateUpdateRequest, admin.ProjectManagementTemplateUpdateRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ProjectManagementTemplateRpcClient.Update(l.ctx, &admin.ProjectManagementTemplateUpdateRequest{Data: data})
	if err != nil {
		return nil, err
	}
	return copierx.Trans[proto.ProjectManagementTemplateEntity, types.ProjectManagementTemplateEntity](rpcResp.Data)
}
