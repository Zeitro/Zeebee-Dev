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

type ProjectManagementTemplateDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewProjectManagementTemplateDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTemplateDetailLogic {
	return &ProjectManagementTemplateDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTemplateDetailLogic) ProjectManagementTemplateDetail(req *types.PathIdRequest) (resp *types.ProjectManagementTemplateEntity, err error) {
	rpcResp, err := l.svcCtx.ProjectManagementTemplateRpcClient.Detail(l.ctx, &admin.ProjectManagementTemplateDetailRequest{Data: &admin.ProjectManagementTemplateDetailRequestData{Id: req.Id}})
	if err != nil {
		return nil, err
	}
	return copierx.Trans[proto.ProjectManagementTemplateEntity, types.ProjectManagementTemplateEntity](rpcResp.Data)
}
