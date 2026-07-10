package logic

import (
	"context"
	"math"
	"strconv"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/payment"
)

// metadata keys attached to checkout sessions and Stripe subscriptions so
// webhook handlers can resolve the local entities.
const (
	stripeMetaOrganizationId = "organization_id"
	stripeMetaPlanId         = "pricing_plan_id"
	stripeMetaPolicyId       = "pricing_plan_policy_id"
	stripeMetaUserType       = "user_type"
	stripeMetaOperatorUserId = "operator_user_id"
)

// StripePaymentLogic implements the self-serve payment flows on top of the
// existing subscription state machine. Stripe is the billing source of truth;
// the local organization_subscriptions table stays the authorization source
// of truth and is synced by webhook events.
type StripePaymentLogic struct {
	*BaseLogic
	authTopicMap      map[constants.UserType]string
	subscriptionLogic *OrganizationSubscriptionLogic
}

func NewStripePaymentLogic(baseLogic *BaseLogic, authTopicMap map[constants.UserType]string, subscriptionLogic *OrganizationSubscriptionLogic) *StripePaymentLogic {
	return &StripePaymentLogic{
		BaseLogic:         baseLogic,
		authTopicMap:      authTopicMap,
		subscriptionLogic: subscriptionLogic,
	}
}

// Current returns the organization's newest effective subscription with its
// plan and policy attached, or nil when the org has no effective subscription.
func (l *StripePaymentLogic) Current(ctx context.Context, organizationId int64) (*dto.OrganizationSubscriptionEntity, error) {
	subscriptions, err := l.repos.OrganizationSubscriptionRepo.FindActiveByOrganizationId(ctx, organizationId)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to find organization subscriptions")
	}
	if len(subscriptions) == 0 {
		return nil, nil
	}

	subscription := subscriptions[0]
	result, err := l.subscriptionLogic.transToDetailDto(ctx, subscription)
	if err != nil {
		return nil, err
	}
	if subscription.PricingPlanPolicyId > 0 {
		policy, policyErr := l.repos.PricingPlanPolicyRepo.GetOne(ctx, subscription.PricingPlanPolicyId)
		if policyErr == nil && policy != nil {
			result.PricingPlanPolicy, err = copierx.Trans[model.PricingPlanPolicy, dto.PricingPlanPolicyEntityDto](policy)
			if err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}

// CreateCheckoutSession validates the request, lazily provisions the Stripe
// customer/product/price, and returns a hosted Checkout Session URL. The
// subscription itself is created later by the checkout.session.completed webhook.
func (l *StripePaymentLogic) CreateCheckoutSession(ctx context.Context, input dto.StripeCheckoutRequestDto) (string, error) {
	gateway, err := l.requirePaymentGateway()
	if err != nil {
		return "", err
	}

	organization, err := l.loadOwnedOrganization(ctx, input.OrganizationId, input.UserId)
	if err != nil {
		return "", err
	}
	if organization.Type != input.UserType {
		return "", errorx.New("organization type not match")
	}

	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, input.PricingPlanId, nil)
	if err != nil {
		return "", errorx.Wrap(err, "fail to find pricing plan")
	}
	if plan.Status != constants.CommonStatusActive {
		return "", errorx.New("pricing plan not found or disabled")
	}
	if plan.Type != organization.Type {
		return "", errorx.New("organization type not match pricing plan type")
	}

	policy, err := l.repos.PricingPlanPolicyRepo.GetOne(ctx, input.PricingPlanPolicyId)
	if err != nil {
		return "", errorx.Wrap(err, "fail to find pricing plan policy")
	}
	if policy.PricingPlanId != plan.Id {
		return "", errorx.New("pricing plan policy not match")
	}
	if policy.ActualAmount <= 0 {
		return "", errorx.New("free pricing policy does not require checkout")
	}

	// v1 keeps upgrades simple: an org with a live Stripe subscription must let
	// it lapse (cancel auto-renew) before checking out a different plan.
	if err = l.rejectWhenStripeSubscriptionEffective(ctx, input.OrganizationId); err != nil {
		return "", err
	}

	customerId, err := l.ensureCustomer(ctx, gateway, organization)
	if err != nil {
		return "", err
	}
	priceId, err := l.ensurePrice(ctx, gateway, plan, policy)
	if err != nil {
		return "", err
	}

	metadata := map[string]string{
		stripeMetaOrganizationId: strconv.FormatInt(organization.Id, 10),
		stripeMetaPlanId:         strconv.FormatInt(plan.Id, 10),
		stripeMetaPolicyId:       strconv.FormatInt(policy.Id, 10),
		stripeMetaUserType:       string(plan.Type),
		stripeMetaOperatorUserId: strconv.FormatInt(input.UserId, 10),
	}
	checkoutUrl, err := gateway.CreateCheckoutSession(ctx, payment.CheckoutParams{
		CustomerId: customerId,
		PriceId:    priceId,
		Metadata:   metadata,
	})
	if err != nil {
		return "", err
	}
	return checkoutUrl, nil
}

// SetAutoRenew toggles cancel_at_period_end on the org's live Stripe
// subscription and mirrors the flag locally. Access is not affected until the
// period actually ends, so no policy rebuild is needed.
func (l *StripePaymentLogic) SetAutoRenew(ctx context.Context, input dto.StripeAutoRenewRequestDto) error {
	gateway, err := l.requirePaymentGateway()
	if err != nil {
		return err
	}

	if _, err = l.loadOwnedOrganization(ctx, input.OrganizationId, input.UserId); err != nil {
		return err
	}

	subscription, err := l.findEffectiveStripeSubscription(ctx, input.OrganizationId)
	if err != nil {
		return err
	}
	if subscription == nil {
		return errorx.New("no active stripe subscription found")
	}
	if subscription.IsAutoRenew == input.Enabled {
		return nil
	}

	if err = gateway.SetCancelAtPeriodEnd(ctx, subscription.StripeSubscriptionId, !input.Enabled); err != nil {
		return err
	}
	return l.repos.OrganizationSubscriptionRepo.UpdateColumnsById(ctx, subscription.Id, map[string]interface{}{
		"is_auto_renew": input.Enabled,
	})
}

// loadOwnedOrganization loads the org and verifies the operator is its owner.
func (l *StripePaymentLogic) loadOwnedOrganization(ctx context.Context, organizationId, userId int64) (*model.Organization, error) {
	organization, err := l.repos.OrganizationRepo.GetOrganizationById(ctx, organizationId, false)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to find organization")
	}
	if organization == nil {
		return nil, errorx.New("organization not found")
	}
	if organization.OwnerId != userId {
		return nil, errorx.New("only the organization owner can manage subscription payment")
	}
	return organization, nil
}

// findEffectiveStripeSubscription returns the org's newest effective
// Stripe-sourced subscription, or nil when there is none.
func (l *StripePaymentLogic) findEffectiveStripeSubscription(ctx context.Context, organizationId int64) (*model.OrganizationSubscription, error) {
	subscriptions, err := l.repos.OrganizationSubscriptionRepo.FindActiveByOrganizationId(ctx, organizationId)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to find organization subscriptions")
	}
	for _, subscription := range subscriptions {
		if subscription.Source == constants.OrganizationSubscriptionSourceStripe &&
			subscription.StripeSubscriptionId != "" {
			return subscription, nil
		}
	}
	return nil, nil
}

func (l *StripePaymentLogic) rejectWhenStripeSubscriptionEffective(ctx context.Context, organizationId int64) error {
	subscription, err := l.findEffectiveStripeSubscription(ctx, organizationId)
	if err != nil {
		return err
	}
	if subscription != nil {
		return errorx.New("organization already has an active paid subscription; cancel auto-renew and wait for it to expire before switching plans")
	}
	return nil
}

// ensureCustomer lazily creates the Stripe customer for the org and stores its id.
func (l *StripePaymentLogic) ensureCustomer(ctx context.Context, gateway payment.Gateway, organization *model.Organization) (string, error) {
	if organization.StripeCustomerId != "" {
		return organization.StripeCustomerId, nil
	}
	email := organization.ContactEmail
	if email == "" {
		if owner, err := l.repos.UserRepo.FindOneById(ctx, organization.OwnerId); err == nil && owner != nil {
			email = owner.Email
		}
	}
	customerId, err := gateway.CreateCustomer(ctx, payment.CustomerParams{
		Email:          email,
		Name:           organization.Name,
		OrganizationId: organization.Id,
	})
	if err != nil {
		return "", err
	}
	if err = l.repos.OrganizationRepo.UpdateColumnsById(ctx, organization.Id, map[string]interface{}{
		"stripe_customer_id": customerId,
	}); err != nil {
		return "", err
	}
	organization.StripeCustomerId = customerId
	return customerId, nil
}

// ensurePrice lazily creates the Stripe product and recurring price for the
// plan/policy pair and stores their ids.
func (l *StripePaymentLogic) ensurePrice(ctx context.Context, gateway payment.Gateway, plan *model.PricingPlan, policy *model.PricingPlanPolicy) (string, error) {
	if policy.StripePriceId != "" {
		return policy.StripePriceId, nil
	}

	productId := plan.StripeProductId
	if productId == "" {
		var err error
		productId, err = gateway.CreateProduct(ctx, payment.ProductParams{
			Name:          plan.Name,
			PricingPlanId: plan.Id,
		})
		if err != nil {
			return "", err
		}
		if err = l.repos.PricingPlanRepo.UpdateColumnsById(ctx, plan.Id, map[string]interface{}{
			"stripe_product_id": productId,
		}); err != nil {
			return "", err
		}
		plan.StripeProductId = productId
	}

	interval, intervalCount := policyRecurring(policy.Type)
	priceId, err := gateway.CreatePrice(ctx, payment.PriceParams{
		ProductId:       productId,
		UnitAmountCents: amountToCents(policy.ActualAmount),
		Interval:        interval,
		IntervalCount:   intervalCount,
		PolicyId:        policy.Id,
	})
	if err != nil {
		return "", err
	}
	if err = l.repos.PricingPlanPolicyRepo.UpdateColumnsById(ctx, policy.Id, map[string]interface{}{
		"stripe_price_id": priceId,
	}); err != nil {
		return "", err
	}
	policy.StripePriceId = priceId
	return priceId, nil
}

// policyRecurring maps a billing cycle to Stripe recurring interval params.
func policyRecurring(policyType constants.PricingPolicyType) (interval string, intervalCount int64) {
	switch policyType {
	case constants.PricingPolicyTypeQuarterly:
		return "month", 3
	case constants.PricingPolicyTypeYearly:
		return "year", 1
	default:
		return "month", 1
	}
}

// amountToCents converts a decimal(10,2) amount to integer cents, rounding to
// avoid float truncation (e.g. 19.99 * 100 = 1998.9999...).
func amountToCents(amount float64) int64 {
	return int64(math.Round(amount * 100))
}

// planPreloadResources is a shorthand used by webhook handlers.
var planPreloadResources = []gormx.Preload{{Query: "Resources"}}
