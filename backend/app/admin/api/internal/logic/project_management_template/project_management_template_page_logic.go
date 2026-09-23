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

type ProjectManagementTemplatePageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// page
func NewProjectManagementTemplatePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementTemplatePageLogic {
	return &ProjectManagementTemplatePageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementTemplatePageLogic) ProjectManagementTemplatePage(req *types.ProjectManagementTemplatePageRequest) (resp *types.ProjectManagementTemplatePageResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementTemplatePageRequest, admin.ProjectManagementTemplatePageRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ProjectManagementTemplateRpcClient.Page(l.ctx, &admin.ProjectManagementTemplatePageRequest{Data: data})
	if err != nil {
		return nil, err
	}
	list, err := copierx.Trans[[]*proto.ProjectManagementTemplateEntity, []types.ProjectManagementTemplateEntity](&rpcResp.List)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementTemplatePageResponse{Total: rpcResp.Total, Page: rpcResp.Page, Limit: rpcResp.Limit, TotalPage: rpcResp.TotalPage, List: *list}, nil
}
