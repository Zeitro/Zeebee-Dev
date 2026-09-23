package projectmanagementstagelogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/app/admin/rpc/internal/logic"
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

func (l *SwapOrderLogic) SwapOrder(in *admin.ProjectManagementStageSwapOrderRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	base := logic.GetBaseLogic(l.svcCtx)
	claim, err := base.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err = base.ProjectManagementTemplateLogic.SwapStageOrder(l.ctx, claim, in.Data.TemplateId, in.Data.StageId, in.Data.TargetStageId); err != nil {
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}
