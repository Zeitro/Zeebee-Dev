package subscription

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubscriptionCancelRenewalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// cancel auto renewal
func NewSubscriptionCancelRenewalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscriptionCancelRenewalLogic {
	return &SubscriptionCancelRenewalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubscriptionCancelRenewalLogic) SubscriptionCancelRenewal(req *types.EmptyRequest) (resp *types.EmptyResponse, err error) {
	_, err = l.svcCtx.OrganizationSubscriptionRpcClient.SetAutoRenew(l.ctx, &lender.SubscriptionAutoRenewRequest{
		Data: &lender.SubscriptionAutoRenewRequestData{Enabled: false},
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}
