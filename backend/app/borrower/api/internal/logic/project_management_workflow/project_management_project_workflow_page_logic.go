package project_management_workflow

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectManagementProjectWorkflowPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// project workflow page
func NewProjectManagementProjectWorkflowPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementProjectWorkflowPageLogic {
	return &ProjectManagementProjectWorkflowPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementProjectWorkflowPageLogic) ProjectManagementProjectWorkflowPage(req *types.ProjectManagementProjectWorkflowPageRequest) (resp *types.ProjectManagementWorkflowPageResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementProjectWorkflowPageRequest, borrower.ProjectManagementProjectWorkflowPageRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ProjectManagementWorkflowRpcClient.ProjectWorkflowPage(l.ctx, &borrower.ProjectManagementProjectWorkflowPageRequest{Data: data})
	if err != nil {
		return nil, err
	}
	list, err := copierx.Trans[[]*proto.ProjectManagementWorkflowEntity, []*types.ProjectManagementWorkflowEntity](&rpcResp.List)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementWorkflowPageResponse{Total: rpcResp.Total, Page: rpcResp.Page, Limit: rpcResp.Limit, TotalPage: rpcResp.TotalPage, List: *list}, nil
}
