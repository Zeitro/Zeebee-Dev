package project_management_workflow_transition

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create
func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.ProjectManagementWorkflowTransitionCreateRequest) (resp *types.ProjectManagementWorkflowTransitionDetailResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementWorkflowTransitionCreateRequest, borrower.ProjectManagementWorkflowTransitionCreateRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.WorkflowTransitionRpcClient.Create(l.ctx, &borrower.ProjectManagementWorkflowTransitionCreateRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	entity, err := copierx.Trans[proto.ProjectManagementWorkflowTransitionEntity, types.ProjectManagementWorkflowTransitionEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementWorkflowTransitionDetailResponse{ProjectManagementWorkflowTransitionEntity: *entity}, nil
}
