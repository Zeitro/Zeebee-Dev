package projectmanagementworkflowlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type ProjectWorkflowDeleteLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewProjectWorkflowDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectWorkflowDeleteLogic {
	return &ProjectWorkflowDeleteLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ProjectWorkflowDeleteLogic) ProjectWorkflowDelete(in *borrower.ProjectManagementProjectWorkflowDeleteRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err = l.baseLogic.ProjectManagementWorkflowLogic.ProjectWorkflowDelete(l.ctx, claim, in.Data.ProjectId, in.Data.WorkflowId); err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
