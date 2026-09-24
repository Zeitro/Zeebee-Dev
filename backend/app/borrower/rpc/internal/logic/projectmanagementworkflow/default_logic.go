package projectmanagementworkflowlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type DefaultLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewDefaultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DefaultLogic {
	return &DefaultLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DefaultLogic) Default(in *proto.PathIdRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		return nil, errorx.New("request data is nil")
	}
	base := l.baseLogic
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err = base.ProjectManagementWorkflowLogic.SetDefault(l.ctx, claim, in.Data.Id); err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
