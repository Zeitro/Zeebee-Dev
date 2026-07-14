package subscription

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
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
	_, err = l.svcCtx.OrganizationSubscriptionRpcClient.SetAutoRenew(l.ctx, &borrower.SubscriptionAutoRenewRequest{
		Data: &borrower.SubscriptionAutoRenewRequestData{Enabled: false},
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}
