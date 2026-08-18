package projectmanagementlogic

import (
	"context"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type WorkflowDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkflowDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkflowDeleteLogic {
	return &WorkflowDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkflowDeleteLogic) WorkflowDelete(in *borrower.ProjectManagementWorkflowDeleteRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	actor, err := requireProjectManagementActor(l.ctx, base)
	if err != nil {
		return nil, err
	}
	if err := base.ProjectManagementLogic.DeleteWorkflow(l.ctx, actor, in.Data.Id); err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &proto.EmptyResponse{}, nil
}
