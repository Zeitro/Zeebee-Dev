package subscription

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubscriptionCurrentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// current organization subscription
func NewSubscriptionCurrentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscriptionCurrentLogic {
	return &SubscriptionCurrentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubscriptionCurrentLogic) SubscriptionCurrent(req *types.EmptyRequest) (resp *types.SubscriptionCurrentResponse, err error) {
	rpcResp, err := l.svcCtx.OrganizationSubscriptionRpcClient.Current(l.ctx, &lender.EmptyRequest{})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if rpcResp.Data == nil {
		return &types.SubscriptionCurrentResponse{}, nil
	}

	subscription, err := copierx.Trans[lender.OrganizationSubscriptionEntity, types.SubscriptionCurrentEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.SubscriptionCurrentResponse{Subscription: subscription}, nil
}
