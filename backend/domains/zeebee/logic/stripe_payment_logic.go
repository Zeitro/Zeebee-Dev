package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"math"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v82"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

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

// HandleWebhookEvent verifies the signature, applies the idempotency gate, and
// dispatches the event to its handler. Returning an error makes the transport
// answer 5xx so Stripe retries; signature failures carry
// payment.ErrSignatureVerification and must map to 400 (no retry).
func (l *StripePaymentLogic) HandleWebhookEvent(ctx context.Context, payload []byte, signature string) error {
	gateway, err := l.requirePaymentGateway()
	if err != nil {
		return err
	}
	event, err := gateway.ConstructEvent(payload, signature)
	if err != nil {
		return err
	}

	// idempotency gate: the unique event_id insert wins exactly once; replays
	// of Processed/Skipped events return success without touching anything.
	created, err := l.repos.StripeEventRepo.CreateIgnoreConflict(ctx, &model.StripeEvent{
		EventId:   event.ID,
		EventType: string(event.Type),
		Status:    constants.StripeEventStatusProcessing,
		Payload:   string(event.Data.Raw),
	})
	if err != nil {
		return err
	}
	if !created {
		existing, findErr := l.repos.StripeEventRepo.FindOneByEventId(ctx, event.ID)
		if findErr != nil {
			return findErr
		}
		if existing.Status != constants.StripeEventStatusFailed {
			return nil
		}
		// Failed events are retried by Stripe: reopen and process again
		if err = l.repos.StripeEventRepo.UpdateStatusByEventId(ctx, event.ID, constants.StripeEventStatusProcessing, ""); err != nil {
			return err
		}
	}

	handled, err := l.dispatchEvent(ctx, event)
	if err != nil {
		_ = l.repos.StripeEventRepo.UpdateStatusByEventId(ctx, event.ID, constants.StripeEventStatusFailed, err.Error())
		return err
	}
	status := constants.StripeEventStatusProcessed
	if !handled {
		status = constants.StripeEventStatusSkipped
	}
	return l.repos.StripeEventRepo.UpdateStatusByEventId(ctx, event.ID, status, "")
}

// dispatchEvent routes a verified event; handled=false means the event type is
// not consumed (recorded as Skipped, answered 200 so Stripe stops retrying).
func (l *StripePaymentLogic) dispatchEvent(ctx context.Context, event stripe.Event) (handled bool, err error) {
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err = json.Unmarshal(event.Data.Raw, &session); err != nil {
			return true, errorx.Wrap(err, "fail to parse checkout session event")
		}
		return true, l.activateStripeSubscription(ctx, &session)
	case "invoice.paid":
		var invoice stripe.Invoice
		if err = json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			return true, errorx.Wrap(err, "fail to parse invoice event")
		}
		return true, l.handleInvoicePaid(ctx, &invoice)
	case "invoice.payment_failed":
		var invoice stripe.Invoice
		if err = json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			return true, errorx.Wrap(err, "fail to parse invoice event")
		}
		return true, l.handleInvoicePaymentFailed(ctx, &invoice)
	case "customer.subscription.updated":
		var subscription stripe.Subscription
		if err = json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			return true, errorx.Wrap(err, "fail to parse subscription event")
		}
		return true, l.handleSubscriptionUpdated(ctx, &subscription)
	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		if err = json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			return true, errorx.Wrap(err, "fail to parse subscription event")
		}
		return true, l.handleSubscriptionDeleted(ctx, &subscription)
	default:
		return false, nil
	}
}

// activateStripeSubscription is the main activation path, driven by
// checkout.session.completed. It replaces the org's currently effective
// subscriptions (mirroring OrganizationSubscriptionLogic.Create) with a new
// Stripe-sourced Active subscription.
func (l *StripePaymentLogic) activateStripeSubscription(ctx context.Context, session *stripe.CheckoutSession) error {
	if session.Subscription == nil || session.Subscription.ID == "" {
		return errorx.New("checkout session has no subscription")
	}
	stripeSubscriptionId := session.Subscription.ID

	// replayed or duplicated events: the subscription is already activated
	existing, err := l.repos.OrganizationSubscriptionRepo.FindOneByStripeSubscriptionId(ctx, stripeSubscriptionId)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}

	organizationId := metadataInt64(session.Metadata, stripeMetaOrganizationId)
	planId := metadataInt64(session.Metadata, stripeMetaPlanId)
	policyId := metadataInt64(session.Metadata, stripeMetaPolicyId)
	if organizationId <= 0 || planId <= 0 || policyId <= 0 {
		return errorx.New("checkout session metadata is incomplete")
	}

	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, planId, planPreloadResources)
	if err != nil {
		return errorx.Wrap(err, "fail to find pricing plan")
	}
	policy, err := l.repos.PricingPlanPolicyRepo.GetOne(ctx, policyId)
	if err != nil {
		return errorx.Wrap(err, "fail to find pricing plan policy")
	}

	expireDate := l.stripePeriodEndOrFallback(ctx, stripeSubscriptionId, policy.Type)

	code, err := l.subscriptionLogic.generateCode(ctx, nil)
	if err != nil {
		return err
	}
	subscription := &model.OrganizationSubscription{
		Code:                 code,
		OrganizationId:       organizationId,
		PricingPlanId:        plan.Id,
		PricingPlanPolicyId:  policy.Id,
		Status:               constants.OrganizationSubscriptionStatusActive,
		ExpireDate:           expireDate,
		OriginalAmount:       policy.OriginalAmount,
		ActualAmount:         policy.ActualAmount,
		IsAutoRenew:          true,
		Source:               constants.OrganizationSubscriptionSourceStripe,
		StripeSubscriptionId: stripeSubscriptionId,
	}

	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		oldSubs, innerErr := l.repos.OrganizationSubscriptionRepo.FindActiveByOrgIdWithTx(ctx, tx, organizationId)
		if innerErr != nil {
			return innerErr
		}
		// remove-side: unbind resources of currently effective old subscriptions
		if _, innerErr = unbindEffectiveSubscriptions(ctx, l.repos, tx, oldSubs); innerErr != nil {
			return innerErr
		}
		innerErr = l.repos.OrganizationSubscriptionRepo.UpdateColumnsWithTx(ctx, tx,
			map[string]interface{}{"status": constants.OrganizationSubscriptionStatusExpired},
			map[string]interface{}{"organization_id": organizationId},
			map[string]interface{}{"status": []constants.OrganizationSubscriptionStatus{
				constants.OrganizationSubscriptionStatusExpired,
				constants.OrganizationSubscriptionStatusSuspended,
			}},
		)
		if innerErr != nil {
			return innerErr
		}
		subscription, innerErr = l.repos.OrganizationSubscriptionRepo.CreateWithTx(ctx, tx, subscription, false)
		return innerErr
	})
	if err != nil {
		return err
	}

	// backfill the customer id in case the org record predates checkout
	if session.Customer != nil && session.Customer.ID != "" {
		_ = l.repos.OrganizationRepo.UpdateColumnsById(ctx, organizationId, map[string]interface{}{
			"stripe_customer_id": session.Customer.ID,
		})
	}

	l.subscriptionLogic.publishOrgPolicyRebuild(ctx, []int64{organizationId}, plan.Type)
	return nil
}

// handleInvoicePaid extends the local subscription to the invoice's period end
// and records the payment. Renewal invoices arrive on every billing cycle; the
// first invoice may arrive before checkout.session.completed, in which case
// only the transaction is recorded (activation owns subscription creation).
func (l *StripePaymentLogic) handleInvoicePaid(ctx context.Context, invoice *stripe.Invoice) error {
	stripeSubscriptionId := invoiceSubscriptionId(invoice)
	if stripeSubscriptionId == "" {
		return nil
	}

	subscription, err := l.repos.OrganizationSubscriptionRepo.FindOneByStripeSubscriptionId(ctx, stripeSubscriptionId)
	if err != nil {
		return err
	}

	organizationId := metadataInt64(invoiceSubscriptionMetadata(invoice), stripeMetaOrganizationId)
	subscriptionId := int64(0)
	if subscription != nil {
		organizationId = subscription.OrganizationId
		subscriptionId = subscription.Id
	}
	transaction := &model.PaymentTransaction{
		OrganizationId:             organizationId,
		OrganizationSubscriptionId: subscriptionId,
		StripeInvoiceId:            invoice.ID,
		StripeSubscriptionId:       stripeSubscriptionId,
		Amount:                     centsToAmount(invoice.AmountPaid),
		Currency:                   string(invoice.Currency),
		Status:                     constants.PaymentTransactionStatusPaid,
		PeriodStart:                unixToNullDate(invoicePeriodStart(invoice)),
		PeriodEnd:                  unixToNullDate(invoicePeriodEnd(invoice)),
		PaidAt:                     &sql.NullTime{Time: time.Now(), Valid: true},
	}

	if subscription == nil {
		return l.repos.PaymentTransactionRepo.UpsertByInvoiceIdWithTx(ctx, l.repos.PaymentTransactionRepo.GetDB(), transaction)
	}

	oldEffective := subscription.Status.IsEffective()
	statusChanged := false
	// Suspended is an explicit admin action and is never auto-reactivated
	if subscription.Status != constants.OrganizationSubscriptionStatusSuspended &&
		subscription.Status != constants.OrganizationSubscriptionStatusActive {
		subscription.Status = constants.OrganizationSubscriptionStatusActive
		statusChanged = true
	}
	if periodEnd := invoicePeriodEnd(invoice); periodEnd > 0 {
		subscription.ExpireDate = unixToNullDate(periodEnd)
	}

	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, innerErr := l.repos.OrganizationSubscriptionRepo.UpdateOneByIdWithTx(ctx, tx, subscription.Id, subscription); innerErr != nil {
			return innerErr
		}
		return l.repos.PaymentTransactionRepo.UpsertByInvoiceIdWithTx(ctx, tx, transaction)
	})
	if err != nil {
		return err
	}

	if statusChanged && !oldEffective {
		// regained access (e.g. Expired -> Active after a late payment)
		l.publishRebuildForSubscription(ctx, subscription)
	}
	return nil
}

// handleInvoicePaymentFailed only records the failure. Stripe smart retries
// own the dunning process; the final state arrives via subscription events.
func (l *StripePaymentLogic) handleInvoicePaymentFailed(ctx context.Context, invoice *stripe.Invoice) error {
	stripeSubscriptionId := invoiceSubscriptionId(invoice)
	if stripeSubscriptionId == "" {
		return nil
	}
	subscription, err := l.repos.OrganizationSubscriptionRepo.FindOneByStripeSubscriptionId(ctx, stripeSubscriptionId)
	if err != nil {
		return err
	}
	organizationId := metadataInt64(invoiceSubscriptionMetadata(invoice), stripeMetaOrganizationId)
	subscriptionId := int64(0)
	if subscription != nil {
		organizationId = subscription.OrganizationId
		subscriptionId = subscription.Id
	}
	logx.WithContext(ctx).WithFields(
		logx.Field("stripe_invoice_id", invoice.ID),
		logx.Field("stripe_subscription_id", stripeSubscriptionId),
		logx.Field("organization_id", organizationId),
	).Error("stripe invoice payment failed")

	return l.repos.PaymentTransactionRepo.UpsertByInvoiceIdWithTx(ctx, l.repos.PaymentTransactionRepo.GetDB(), &model.PaymentTransaction{
		OrganizationId:             organizationId,
		OrganizationSubscriptionId: subscriptionId,
		StripeInvoiceId:            invoice.ID,
		StripeSubscriptionId:       stripeSubscriptionId,
		Amount:                     centsToAmount(invoice.AmountDue),
		Currency:                   string(invoice.Currency),
		Status:                     constants.PaymentTransactionStatusFailed,
		PeriodStart:                unixToNullDate(invoicePeriodStart(invoice)),
		PeriodEnd:                  unixToNullDate(invoicePeriodEnd(invoice)),
	})
}

// handleSubscriptionUpdated mirrors auto-renew and period changes. Stripe's
// terminal statuses are delegated to the deleted flow; past_due keeps access
// during the retry window (deny happens on subscription.deleted or cron).
func (l *StripePaymentLogic) handleSubscriptionUpdated(ctx context.Context, stripeSubscription *stripe.Subscription) error {
	subscription, err := l.repos.OrganizationSubscriptionRepo.FindOneByStripeSubscriptionId(ctx, stripeSubscription.ID)
	if err != nil || subscription == nil {
		// pre-activation updates are expected noise (activation owns creation)
		return err
	}

	switch stripeSubscription.Status {
	case stripe.SubscriptionStatusCanceled, stripe.SubscriptionStatusUnpaid, stripe.SubscriptionStatusIncompleteExpired:
		return l.handleSubscriptionDeleted(ctx, stripeSubscription)
	case stripe.SubscriptionStatusIncomplete:
		return nil
	}

	info := payment.NormalizeSubscription(stripeSubscription)
	subscription.IsAutoRenew = !info.CancelAtPeriodEnd
	if !info.CurrentPeriodEnd.IsZero() {
		subscription.ExpireDate = &sql.NullTime{Time: info.CurrentPeriodEnd, Valid: true}
	}
	if subscription.Status != constants.OrganizationSubscriptionStatusSuspended {
		switch stripeSubscription.Status {
		case stripe.SubscriptionStatusTrialing:
			subscription.Status = constants.OrganizationSubscriptionStatusTrialing
		default:
			// active and past_due (grace period) both keep access
			subscription.Status = constants.OrganizationSubscriptionStatusActive
		}
	}

	if _, err = l.repos.OrganizationSubscriptionRepo.UpdateOneById(ctx, subscription.Id, subscription); err != nil {
		return err
	}
	l.publishRebuildForSubscription(ctx, subscription)
	return nil
}

// handleSubscriptionDeleted expires the local subscription and falls back to
// the default free plan, reusing the cron's expired-subscription recovery.
func (l *StripePaymentLogic) handleSubscriptionDeleted(ctx context.Context, stripeSubscription *stripe.Subscription) error {
	subscription, err := l.repos.OrganizationSubscriptionRepo.FindOneByStripeSubscriptionId(ctx, stripeSubscription.ID)
	if err != nil || subscription == nil {
		return err
	}
	if subscription.Status == constants.OrganizationSubscriptionStatusExpired {
		return nil
	}

	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, subscription.PricingPlanId, planPreloadResources)
	if err != nil {
		return errorx.Wrap(err, "fail to find pricing plan")
	}

	recovery := expiredSubscriptionRecovery{Now: time.Now()}
	planFilter := []gormx.WhereRaw{{
		Raw:  "is_default = ? and type = ?",
		Args: []interface{}{true, plan.Type},
	}}
	recovery.DefaultPlan, _ = l.repos.PricingPlanRepo.FindOneByFilter(ctx, planFilter,
		gormx.Preload{Query: "Policies"}, gormx.Preload{Query: "Resources"})
	if recovery.DefaultPlan != nil {
		for _, defaultPolicy := range recovery.DefaultPlan.Policies {
			if defaultPolicy.ActualAmount == 0 {
				recovery.DefaultPolicy = defaultPolicy
				break
			}
		}
	}

	err = l.repos.OrganizationSubscriptionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		oldEffective := subscription.Status.IsEffective()
		if _, innerErr := l.subscriptionLogic.handleExpiredSubscription(ctx, tx, subscription, recovery); innerErr != nil {
			return innerErr
		}
		// remove-side unbind for the plan losing effectiveness
		_, innerErr := syncSubscriptionEffect(ctx, l.repos, tx, subscriptionEffectSyncReq{
			OrganizationId: subscription.OrganizationId,
			Plan:           plan,
			OldEffective:   oldEffective,
			NewEffective:   false,
		})
		return innerErr
	})
	if err != nil {
		return err
	}

	l.subscriptionLogic.publishOrgPolicyRebuild(ctx, []int64{subscription.OrganizationId}, plan.Type)
	return nil
}

// publishRebuildForSubscription publishes a policy rebuild scoped to the
// subscription's organization, resolving the user type from its plan.
func (l *StripePaymentLogic) publishRebuildForSubscription(ctx context.Context, subscription *model.OrganizationSubscription) {
	plan, err := l.repos.PricingPlanRepo.GetOne(ctx, subscription.PricingPlanId, nil)
	if err != nil || plan == nil {
		return
	}
	l.subscriptionLogic.publishOrgPolicyRebuild(ctx, []int64{subscription.OrganizationId}, plan.Type)
}

// stripePeriodEndOrFallback resolves the subscription's real period end from
// Stripe, falling back to the local billing-cycle calculation when the lookup
// fails (webhook processing must not hard-depend on a second API call).
func (l *StripePaymentLogic) stripePeriodEndOrFallback(ctx context.Context, stripeSubscriptionId string, policyType constants.PricingPolicyType) *sql.NullTime {
	if l.paymentGateway != nil {
		if info, err := l.paymentGateway.GetSubscription(ctx, stripeSubscriptionId); err == nil && !info.CurrentPeriodEnd.IsZero() {
			return &sql.NullTime{Time: info.CurrentPeriodEnd, Valid: true}
		}
	}
	return calculateExpireDate(policyType, time.Now())
}

// metadataInt64 reads an int64 metadata value, returning 0 when absent/invalid.
func metadataInt64(metadata map[string]string, key string) int64 {
	value, err := strconv.ParseInt(metadata[key], 10, 64)
	if err != nil {
		return 0
	}
	return value
}

// invoiceSubscriptionId extracts the subscription id an invoice belongs to.
func invoiceSubscriptionId(invoice *stripe.Invoice) string {
	if invoice.Parent == nil || invoice.Parent.SubscriptionDetails == nil ||
		invoice.Parent.SubscriptionDetails.Subscription == nil {
		return ""
	}
	return invoice.Parent.SubscriptionDetails.Subscription.ID
}

// invoiceSubscriptionMetadata returns the subscription metadata snapshot on the invoice.
func invoiceSubscriptionMetadata(invoice *stripe.Invoice) map[string]string {
	if invoice.Parent == nil || invoice.Parent.SubscriptionDetails == nil {
		return nil
	}
	return invoice.Parent.SubscriptionDetails.Metadata
}

// invoicePeriodEnd prefers the line item's service period over the invoice
// period (they differ for subscription invoices).
func invoicePeriodEnd(invoice *stripe.Invoice) int64 {
	if invoice.Lines != nil && len(invoice.Lines.Data) > 0 && invoice.Lines.Data[0].Period != nil {
		return invoice.Lines.Data[0].Period.End
	}
	return invoice.PeriodEnd
}

func invoicePeriodStart(invoice *stripe.Invoice) int64 {
	if invoice.Lines != nil && len(invoice.Lines.Data) > 0 && invoice.Lines.Data[0].Period != nil {
		return invoice.Lines.Data[0].Period.Start
	}
	return invoice.PeriodStart
}

// unixToNullDate converts a unix timestamp to a nullable date value.
func unixToNullDate(unix int64) *sql.NullTime {
	if unix <= 0 {
		return nil
	}
	return &sql.NullTime{Time: time.Unix(unix, 0).UTC(), Valid: true}
}

// centsToAmount converts integer cents to a decimal(10,2)-compatible amount.
func centsToAmount(cents int64) float64 {
	return float64(cents) / 100
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
