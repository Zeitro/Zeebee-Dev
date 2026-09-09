package projectmanagementworkflowstagelogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type SwapOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSwapOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwapOrderLogic {
	return &SwapOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SwapOrderLogic) SwapOrder(in *borrower.ProjectManagementWorkflowStageSwapOrderRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err = base.ProjectManagementWorkflowLogic.StageSwapOrder(l.ctx, claim, in.Data.Id, in.Data.FromStageId, in.Data.ToStageId); err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
