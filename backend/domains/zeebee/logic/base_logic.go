package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/redisx"
	"zeebee-backend/domains/zeebee/auth"
	"zeebee-backend/domains/zeebee/payment"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/repo"
)

// BaseLogic provides base business logic for all logic types
type BaseLogic struct {
	repos          *repo.AllRepos
	redisCli       *redisx.Client
	logger         logx.Logger
	paymentGateway payment.Gateway // nil when Stripe is not configured
	jwtConfig dto.JwtConfig
}

// NewBaseLogic creates a new base logic instance
func NewBaseLogic(repos *repo.AllRepos, redisCli *redisx.Client, logger logx.Logger, paymentGateway payment.Gateway, jwtConfig dto.JwtConfig) *BaseLogic {
	return &BaseLogic{
		repos:          repos,
		redisCli:       redisCli,
		logger:         logger,
		jwtConfig: jwtConfig,
		paymentGateway: paymentGateway,
	}
}

// requirePaymentGateway returns the gateway or an error when Stripe is not
// configured for the current service.
func (l *BaseLogic) requirePaymentGateway() (payment.Gateway, error) {
	if l.paymentGateway == nil {
		return nil, errorx.New("payment is not configured")
	}
	return l.paymentGateway, nil
}

// cancelStripeSubscriptionQuietly cancels a Stripe subscription best-effort.
// Failures are logged loudly for manual follow-up but never block the caller:
// admin operations must not be held hostage by Stripe availability, at the
// cost of a manual-refund window when the remote cancel fails.
func (l *BaseLogic) cancelStripeSubscriptionQuietly(ctx context.Context, stripeSubscriptionId string) {
	if stripeSubscriptionId == "" || l.paymentGateway == nil {
		return
	}
	if err := l.paymentGateway.CancelSubscription(ctx, stripeSubscriptionId); err != nil {
		logx.WithContext(ctx).WithFields(
			logx.Field("stacks", errorx.GetErrorStack(err)),
			logx.Field("stripe_subscription_id", stripeSubscriptionId),
		).Error("failed to cancel stripe subscription; manual follow-up required to stop billing")
	}
}

// archiveStripePriceQuietly archives a Stripe price best-effort: failures are
// logged but never block the caller (admin operations must not depend on
// Stripe availability).
func (l *BaseLogic) archiveStripePriceQuietly(ctx context.Context, priceId string) {
	if priceId == "" || l.paymentGateway == nil {
		return
	}
	if err := l.paymentGateway.ArchivePrice(ctx, priceId); err != nil {
		logx.WithContext(ctx).WithFields(
			logx.Field("stacks", errorx.GetErrorStack(err)),
			logx.Field("stripe_price_id", priceId),
		).Error("failed to archive stripe price")
	}
}

func (l *BaseLogic) publishAuthEvent(ctx context.Context, topic string, event auth.Msg) (err error) {
	defer func() {
		if err != nil {
			logx.WithContext(ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("topic", topic),
				logx.Field("event", event),
			).Error("failed to publish auth event")
		}
	}()

	if l.redisCli == nil {
		err = errorx.New("redis client is nil")
		return
	}

	if len(topic) == 0 {
		err = errorx.New("topic is empty")
		return
	}
	err = l.redisCli.GetClient().Publish(ctx, topic, event.String()).Err()
	return errorx.WithStack(err)
}
