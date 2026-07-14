package subscription

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

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
	rpcResp, err := l.svcCtx.OrganizationSubscriptionRpcClient.Current(l.ctx, &proto.EmptyRequest{})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if rpcResp.Data == nil {
		return &types.SubscriptionCurrentResponse{}, nil
	}

	subscription, err := copierx.Trans[borrower.OrganizationSubscriptionEntity, types.SubscriptionCurrentEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.SubscriptionCurrentResponse{Subscription: subscription}, nil
}
