package projectmanagementworkflowlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type ProjectWorkflowDefaultLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewProjectWorkflowDefaultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectWorkflowDefaultLogic {
	return &ProjectWorkflowDefaultLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ProjectWorkflowDefaultLogic) ProjectWorkflowDefault(in *borrower.ProjectManagementProjectWorkflowDefaultRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err = l.baseLogic.ProjectManagementWorkflowLogic.ProjectWorkflowSetDefault(l.ctx, claim, in.Data.ProjectId, in.Data.WorkflowId); err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
