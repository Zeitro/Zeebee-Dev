package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"
)

type TaskDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskDeleteLogic {
	return &TaskDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskDeleteLogic) TaskDelete(in *borrower.ProjectManagementTaskDeleteRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	actor, err := requireProjectManagementActor(l.ctx, base)
	if err != nil {
		return nil, err
	}
	if err = base.ProjectManagementLogic.DeleteTask(l.ctx, actor, in.Data.Id, in.Data.Confirm); err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &proto.EmptyResponse{}, nil
}
