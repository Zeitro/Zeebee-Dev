package subscription

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubscriptionResumeRenewalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// resume auto renewal
func NewSubscriptionResumeRenewalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscriptionResumeRenewalLogic {
	return &SubscriptionResumeRenewalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubscriptionResumeRenewalLogic) SubscriptionResumeRenewal(req *types.EmptyRequest) (resp *types.EmptyResponse, err error) {
	_, err = l.svcCtx.OrganizationSubscriptionRpcClient.SetAutoRenew(l.ctx, &lender.SubscriptionAutoRenewRequest{
		Data: &lender.SubscriptionAutoRenewRequestData{Enabled: true},
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}
