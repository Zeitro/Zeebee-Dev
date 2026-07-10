package subscription

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
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
	_, err = l.svcCtx.OrganizationSubscriptionRpcClient.SetAutoRenew(l.ctx, &borrower.SubscriptionAutoRenewRequest{
		Data: &borrower.SubscriptionAutoRenewRequestData{Enabled: true},
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return &types.EmptyResponse{}, nil
}
