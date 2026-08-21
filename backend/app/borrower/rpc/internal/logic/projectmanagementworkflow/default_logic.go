package projectmanagementworkflowlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type DefaultLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDefaultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DefaultLogic {
	return &DefaultLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DefaultLogic) Default(in *proto.PathIdRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.New("request data is nil")
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err = base.ProjectManagementWorkflowLogic.SetDefault(l.ctx, claim, in.Data.Id); err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &proto.EmptyResponse{}, nil
}
