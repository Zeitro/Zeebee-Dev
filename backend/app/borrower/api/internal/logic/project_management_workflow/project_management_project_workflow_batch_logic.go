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

type ProjectManagementProjectWorkflowBatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// copy workflows into project
func NewProjectManagementProjectWorkflowBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectManagementProjectWorkflowBatchLogic {
	return &ProjectManagementProjectWorkflowBatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectManagementProjectWorkflowBatchLogic) ProjectManagementProjectWorkflowBatch(req *types.ProjectManagementProjectWorkflowBatchRequest) (resp *types.ProjectManagementWorkflowBatchResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementProjectWorkflowBatchRequest, borrower.ProjectManagementProjectWorkflowBatchRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ProjectManagementWorkflowRpcClient.ProjectWorkflowBatch(l.ctx, &borrower.ProjectManagementProjectWorkflowBatchRequest{Data: data})
	if err != nil {
		return nil, err
	}
	list, err := copierx.Trans[[]*proto.ProjectManagementWorkflowEntity, []*types.ProjectManagementWorkflowEntity](&rpcResp.List)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementWorkflowBatchResponse{List: *list}, nil
}
