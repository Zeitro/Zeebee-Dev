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

type ProjectManagementProjectWorkflowDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// project workflow detail
func NewProjectManagementProjectWorkflowDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementProjectWorkflowDetailLogic {
	return &ProjectManagementProjectWorkflowDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementProjectWorkflowDetailLogic) ProjectManagementProjectWorkflowDetail(req *types.ProjectManagementProjectWorkflowDetailRequest) (resp *types.ProjectManagementWorkflowDetailResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementProjectWorkflowDetailRequest, borrower.ProjectManagementProjectWorkflowDetailRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ProjectManagementWorkflowRpcClient.ProjectWorkflowDetail(l.ctx, &borrower.ProjectManagementProjectWorkflowDetailRequest{Data: data})
	if err != nil {
		return nil, err
	}
	entity, err := copierx.Trans[proto.ProjectManagementWorkflowEntity, types.ProjectManagementWorkflowEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementWorkflowDetailResponse{ProjectManagementWorkflowEntity: *entity}, nil
}
