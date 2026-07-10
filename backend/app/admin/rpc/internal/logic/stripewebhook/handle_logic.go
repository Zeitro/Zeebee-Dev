package stripewebhooklogic

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandleLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
	logx.Logger
}

func NewHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleLogic {
	return &HandleLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

// Handle verifies and processes a Stripe webhook event. Signature failures are
// answered with InvalidArgument so the http layer returns 400 (Stripe must not
// retry them); any other failure maps to 5xx so Stripe retries.
func (l *HandleLogic) Handle(in *admin.StripeWebhookRequest) (resp *admin.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	err = l.baseLogic.StripePaymentLogic.HandleWebhookEvent(l.ctx, in.Payload, in.Signature)
	if err != nil {
		if errors.Is(err, payment.ErrSignatureVerification) {
			return nil, status.Error(codes.InvalidArgument, "invalid stripe signature")
		}
		return nil, err
	}
	return &admin.EmptyResponse{}, nil
}
