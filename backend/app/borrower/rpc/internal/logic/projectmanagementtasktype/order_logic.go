package projectmanagementtasktypelogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type OrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderLogic {
	return &OrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *OrderLogic) Order(in *borrower.ProjectManagementTaskTypeOrderRequest) (resp *proto.EmptyResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	claimValue, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err = l.baseLogic.ProjectManagementTaskTypeLogic.Order(l.ctx, claimValue, in.Data.Id, int(in.Data.SortOrder)); err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return &proto.EmptyResponse{}, nil
}
