package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go/v82"
	"gorm.io/gorm"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/payment"
	"zeebee-backend/domains/zeebee/repo"
)

// fakeGateway implements payment.Gateway for tests: it records calls and
// parses ConstructEvent payloads as plain JSON (no signature check).
type fakeGateway struct {
	customers     int
	products      int
	prices        int
	archived      []string
	cancelled     []string
	cancelAtEnd   map[string]bool
	periodEnd     time.Time
	checkoutCalls []payment.CheckoutParams
}

func newFakeGateway() *fakeGateway {
	return &fakeGateway{
		cancelAtEnd: map[string]bool{},
		periodEnd:   time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC),
	}
}

func (g *fakeGateway) CreateCustomer(ctx context.Context, params payment.CustomerParams) (string, error) {
	g.customers++
	return fmt.Sprintf("cus_%d", g.customers), nil
}
func (g *fakeGateway) CreateProduct(ctx context.Context, params payment.ProductParams) (string, error) {
	g.products++
	return fmt.Sprintf("prod_%d", g.products), nil
}
func (g *fakeGateway) UpdateProductName(ctx context.Context, productId, name string) error {
	return nil
}
func (g *fakeGateway) CreatePrice(ctx context.Context, params payment.PriceParams) (string, error) {
	g.prices++
	return fmt.Sprintf("price_%d", g.prices), nil
}
func (g *fakeGateway) ArchivePrice(ctx context.Context, priceId string) error {
	g.archived = append(g.archived, priceId)
	return nil
}
func (g *fakeGateway) CreateCheckoutSession(ctx context.Context, params payment.CheckoutParams) (string, error) {
	g.checkoutCalls = append(g.checkoutCalls, params)
	return "https://checkout.stripe.com/test-session", nil
}
func (g *fakeGateway) SetCancelAtPeriodEnd(ctx context.Context, subscriptionId string, cancel bool) error {
	g.cancelAtEnd[subscriptionId] = cancel
	return nil
}
func (g *fakeGateway) CancelSubscription(ctx context.Context, subscriptionId string) error {
	g.cancelled = append(g.cancelled, subscriptionId)
	return nil
}
func (g *fakeGateway) GetSubscription(ctx context.Context, subscriptionId string) (*payment.SubscriptionInfo, error) {
	return &payment.SubscriptionInfo{
		Id:               subscriptionId,
		Status:           "active",
		CurrentPeriodEnd: g.periodEnd,
	}, nil
}
func (g *fakeGateway) ConstructEvent(payload []byte, sigHeader string) (stripe.Event, error) {
	var event stripe.Event
	err := json.Unmarshal(payload, &event)
	return event, err
}

func setupStripePaymentLogic(t *testing.T) (*StripePaymentLogic, *fakeGateway, *gorm.DB) {
	t.Helper()
	db := setupTestDB(t)
	repos := repo.NewAllRepos(db)
	gateway := newFakeGateway()
	baseLogic := NewBaseLogic(repos, nil, nil, gateway)
	authTopicMap := map[constants.UserType]string{
		constants.UserTypeBorrower: "test-borrower-topic",
		constants.UserTypeLender:   "test-lender-topic",
	}
	subscriptionLogic := NewOrganizationSubscriptionLogic(baseLogic, authTopicMap)
	return NewStripePaymentLogic(baseLogic, authTopicMap, subscriptionLogic), gateway, db
}

// seedPaidPlan creates an Active borrower plan with a paid monthly policy.
func seedPaidPlan(t *testing.T, db *gorm.DB) (*model.PricingPlan, *model.PricingPlanPolicy) {
	t.Helper()
	plan := &model.PricingPlan{
		Name:     "Pro",
		Type:     constants.UserTypeBorrower,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive,
	}
	require.NoError(t, db.Create(plan).Error)
	policy := &model.PricingPlanPolicy{
		PricingPlanId:  plan.Id,
		OriginalAmount: 50,
		ActualAmount:   19.99,
		Type:           constants.PricingPolicyTypeMonthly,
	}
	require.NoError(t, db.Create(policy).Error)
	return plan, policy
}

// seedOwnedOrg creates a borrower org owned by the given user id.
func seedOwnedOrg(t *testing.T, db *gorm.DB, ownerId int64) *model.Organization {
	t.Helper()
	org := &model.Organization{
		Name:         "Checkout Org",
		Status:       constants.OrganizationStatusActive,
		Type:         constants.UserTypeBorrower,
		OwnerId:      ownerId,
		ContactEmail: "owner@example.com",
	}
	require.NoError(t, db.Create(org).Error)
	return org
}

// --- checkout ---

func TestCheckout_LazyProvisioningAndUrl(t *testing.T) {
	l, gateway, db := setupStripePaymentLogic(t)
	ctx := context.Background()
	plan, policy := seedPaidPlan(t, db)
	org := seedOwnedOrg(t, db, 42)

	url, err := l.CreateCheckoutSession(ctx, dto.StripeCheckoutRequestDto{
		OrganizationId:      org.Id,
		UserId:              42,
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: policy.Id,
		UserType:            constants.UserTypeBorrower,
	})
	require.NoError(t, err)
	assert.Equal(t, "https://checkout.stripe.com/test-session", url)

	// customer/product/price ids are written back for reuse
	var storedOrg model.Organization
	require.NoError(t, db.First(&storedOrg, org.Id).Error)
	assert.Equal(t, "cus_1", storedOrg.StripeCustomerId)
	var storedPolicy model.PricingPlanPolicy
	require.NoError(t, db.First(&storedPolicy, policy.Id).Error)
	assert.Equal(t, "price_1", storedPolicy.StripePriceId)

	// second checkout reuses everything
	_, err = l.CreateCheckoutSession(ctx, dto.StripeCheckoutRequestDto{
		OrganizationId:      org.Id,
		UserId:              42,
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: policy.Id,
		UserType:            constants.UserTypeBorrower,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, gateway.customers)
	assert.Equal(t, 1, gateway.products)
	assert.Equal(t, 1, gateway.prices)

	// metadata carries the local entity ids for the webhook
	require.Len(t, gateway.checkoutCalls, 2)
	metadata := gateway.checkoutCalls[0].Metadata
	assert.Equal(t, fmt.Sprintf("%d", org.Id), metadata[stripeMetaOrganizationId])
	assert.Equal(t, fmt.Sprintf("%d", plan.Id), metadata[stripeMetaPlanId])
	assert.Equal(t, fmt.Sprintf("%d", policy.Id), metadata[stripeMetaPolicyId])
}

func TestCheckout_RejectsNonOwner(t *testing.T) {
	l, _, db := setupStripePaymentLogic(t)
	plan, policy := seedPaidPlan(t, db)
	org := seedOwnedOrg(t, db, 42)

	_, err := l.CreateCheckoutSession(context.Background(), dto.StripeCheckoutRequestDto{
		OrganizationId:      org.Id,
		UserId:              99, // not the owner
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: policy.Id,
		UserType:            constants.UserTypeBorrower,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "owner")
}

func TestCheckout_RejectsFreePolicyAndLiveStripeSubscription(t *testing.T) {
	l, _, db := setupStripePaymentLogic(t)
	ctx := context.Background()
	plan, policy := seedPaidPlan(t, db)
	org := seedOwnedOrg(t, db, 42)

	// free policy is not checkout-able
	freePolicy := &model.PricingPlanPolicy{
		PricingPlanId: plan.Id,
		ActualAmount:  0,
		Type:          constants.PricingPolicyTypeYearly,
	}
	require.NoError(t, db.Create(freePolicy).Error)
	_, err := l.CreateCheckoutSession(ctx, dto.StripeCheckoutRequestDto{
		OrganizationId: org.Id, UserId: 42,
		PricingPlanId: plan.Id, PricingPlanPolicyId: freePolicy.Id,
		UserType: constants.UserTypeBorrower,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "free")

	// an effective Stripe subscription blocks a second checkout (v1 policy)
	require.NoError(t, db.Create(&model.OrganizationSubscription{
		Code: "SUB1", OrganizationId: org.Id, PricingPlanId: plan.Id,
		PricingPlanPolicyId:  policy.Id,
		Status:               constants.OrganizationSubscriptionStatusActive,
		ExpireDate:           calculateExpireDate(constants.PricingPolicyTypeMonthly, time.Now()),
		Source:               constants.OrganizationSubscriptionSourceStripe,
		StripeSubscriptionId: "sub_live",
	}).Error)
	_, err = l.CreateCheckoutSession(ctx, dto.StripeCheckoutRequestDto{
		OrganizationId: org.Id, UserId: 42,
		PricingPlanId: plan.Id, PricingPlanPolicyId: policy.Id,
		UserType: constants.UserTypeBorrower,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "active paid subscription")
}

// --- webhook ---

func webhookEvent(t *testing.T, id, eventType string, object interface{}) []byte {
	t.Helper()
	raw, err := json.Marshal(object)
	require.NoError(t, err)
	payload, err := json.Marshal(map[string]interface{}{
		"id":   id,
		"type": eventType,
		"data": map[string]interface{}{"object": json.RawMessage(raw)},
	})
	require.NoError(t, err)
	return payload
}

func checkoutCompletedEvent(t *testing.T, eventId string, org *model.Organization, plan *model.PricingPlan, policy *model.PricingPlanPolicy, stripeSubId string) []byte {
	t.Helper()
	return webhookEvent(t, eventId, "checkout.session.completed", map[string]interface{}{
		"id":           "cs_test",
		"subscription": map[string]interface{}{"id": stripeSubId},
		"customer":     map[string]interface{}{"id": "cus_hook"},
		"metadata": map[string]string{
			stripeMetaOrganizationId: fmt.Sprintf("%d", org.Id),
			stripeMetaPlanId:         fmt.Sprintf("%d", plan.Id),
			stripeMetaPolicyId:       fmt.Sprintf("%d", policy.Id),
			stripeMetaUserType:       string(plan.Type),
		},
	})
}

func TestWebhook_CheckoutCompletedActivatesSubscription(t *testing.T) {
	l, _, db := setupStripePaymentLogic(t)
	ctx := context.Background()
	plan, policy := seedPaidPlan(t, db)
	org := seedOwnedOrg(t, db, 42)

	// existing trial subscription should be replaced
	require.NoError(t, db.Create(&model.OrganizationSubscription{
		Code: "TRIAL1", OrganizationId: org.Id, PricingPlanId: plan.Id,
		Status:     constants.OrganizationSubscriptionStatusTrialing,
		ExpireDate: calculateExpireDate(constants.PricingPolicyTypeMonthly, time.Now()),
	}).Error)

	payload := checkoutCompletedEvent(t, "evt_1", org, plan, policy, "sub_abc")
	require.NoError(t, l.HandleWebhookEvent(ctx, payload, "sig"))

	var subs []*model.OrganizationSubscription
	require.NoError(t, db.Where("organization_id = ?", org.Id).Order("id").Find(&subs).Error)
	require.Len(t, subs, 2)
	assert.Equal(t, constants.OrganizationSubscriptionStatusExpired, subs[0].Status)
	newSub := subs[1]
	assert.Equal(t, constants.OrganizationSubscriptionStatusActive, newSub.Status)
	assert.Equal(t, constants.OrganizationSubscriptionSourceStripe, newSub.Source)
	assert.Equal(t, "sub_abc", newSub.StripeSubscriptionId)
	assert.True(t, newSub.IsAutoRenew)
	assert.Equal(t, policy.ActualAmount, newSub.ActualAmount)
	// expire date comes from the fake gateway's period end
	require.NotNil(t, newSub.ExpireDate)
	assert.Equal(t, "2026-08-10", newSub.ExpireDate.Time.Format("2006-01-02"))

	// customer id backfilled from the session
	var storedOrg model.Organization
	require.NoError(t, db.First(&storedOrg, org.Id).Error)
	assert.Equal(t, "cus_hook", storedOrg.StripeCustomerId)

	// event recorded as processed
	var event model.StripeEvent
	require.NoError(t, db.Where("event_id = ?", "evt_1").First(&event).Error)
	assert.Equal(t, constants.StripeEventStatusProcessed, event.Status)
}

func TestWebhook_ReplayedEventIsIdempotent(t *testing.T) {
	l, _, db := setupStripePaymentLogic(t)
	ctx := context.Background()
	plan, policy := seedPaidPlan(t, db)
	org := seedOwnedOrg(t, db, 42)

	payload := checkoutCompletedEvent(t, "evt_dup", org, plan, policy, "sub_dup")
	require.NoError(t, l.HandleWebhookEvent(ctx, payload, "sig"))
	require.NoError(t, l.HandleWebhookEvent(ctx, payload, "sig"))

	var count int64
	require.NoError(t, db.Model(&model.OrganizationSubscription{}).
		Where("stripe_subscription_id = ?", "sub_dup").Count(&count).Error)
	assert.Equal(t, int64(1), count)
}

func TestWebhook_UnknownEventIsSkipped(t *testing.T) {
	l, _, db := setupStripePaymentLogic(t)
	payload := webhookEvent(t, "evt_unknown", "customer.created", map[string]interface{}{"id": "cus_x"})
	require.NoError(t, l.HandleWebhookEvent(context.Background(), payload, "sig"))

	var event model.StripeEvent
	require.NoError(t, db.Where("event_id = ?", "evt_unknown").First(&event).Error)
	assert.Equal(t, constants.StripeEventStatusSkipped, event.Status)
}

func TestWebhook_InvoicePaidBeforeActivationOnlyRecordsTransaction(t *testing.T) {
	l, _, db := setupStripePaymentLogic(t)
	ctx := context.Background()

	payload := webhookEvent(t, "evt_inv_early", "invoice.paid", map[string]interface{}{
		"id":          "in_early",
		"amount_paid": 1999,
		"currency":    "usd",
		"parent": map[string]interface{}{
			"subscription_details": map[string]interface{}{
				"subscription": map[string]interface{}{"id": "sub_not_yet"},
				"metadata":     map[string]string{stripeMetaOrganizationId: "7"},
			},
		},
	})
	require.NoError(t, l.HandleWebhookEvent(ctx, payload, "sig"))

	var transaction model.PaymentTransaction
	require.NoError(t, db.Where("stripe_invoice_id = ?", "in_early").First(&transaction).Error)
	assert.Equal(t, int64(7), transaction.OrganizationId)
	assert.Equal(t, int64(0), transaction.OrganizationSubscriptionId)
	assert.InDelta(t, 19.99, transaction.Amount, 0.001)
	assert.Equal(t, constants.PaymentTransactionStatusPaid, transaction.Status)

	// no subscription was created by the invoice path
	var count int64
	require.NoError(t, db.Model(&model.OrganizationSubscription{}).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestWebhook_InvoicePaidExtendsExpireDate(t *testing.T) {
	l, _, db := setupStripePaymentLogic(t)
	ctx := context.Background()
	plan, policy := seedPaidPlan(t, db)
	org := seedOwnedOrg(t, db, 42)

	payload := checkoutCompletedEvent(t, "evt_act", org, plan, policy, "sub_renew")
	require.NoError(t, l.HandleWebhookEvent(ctx, payload, "sig"))

	nextPeriodEnd := time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)
	renewalPayload := webhookEvent(t, "evt_inv_renew", "invoice.paid", map[string]interface{}{
		"id":          "in_renew",
		"amount_paid": 1999,
		"currency":    "usd",
		"lines": map[string]interface{}{
			"data": []map[string]interface{}{
				{"period": map[string]int64{"start": nextPeriodEnd.AddDate(0, -1, 0).Unix(), "end": nextPeriodEnd.Unix()}},
			},
		},
		"parent": map[string]interface{}{
			"subscription_details": map[string]interface{}{
				"subscription": map[string]interface{}{"id": "sub_renew"},
			},
		},
	})
	require.NoError(t, l.HandleWebhookEvent(ctx, renewalPayload, "sig"))

	var subscription model.OrganizationSubscription
	require.NoError(t, db.Where("stripe_subscription_id = ?", "sub_renew").First(&subscription).Error)
	require.NotNil(t, subscription.ExpireDate)
	assert.Equal(t, "2026-09-10", subscription.ExpireDate.Time.Format("2006-01-02"))
	assert.Equal(t, constants.OrganizationSubscriptionStatusActive, subscription.Status)
}

func TestWebhook_SubscriptionDeletedExpiresAndFallsBackToDefaultPlan(t *testing.T) {
	l, _, db := setupStripePaymentLogic(t)
	ctx := context.Background()
	plan, policy := seedPaidPlan(t, db)
	org := seedOwnedOrg(t, db, 42)

	// default free plan available for recovery
	defaultPlan := &model.PricingPlan{
		Name: "Basic", Type: constants.UserTypeBorrower,
		PlanType: constants.PricingPlanTypeStandard,
		Status:   constants.CommonStatusActive, IsDefault: true,
	}
	require.NoError(t, db.Create(defaultPlan).Error)
	require.NoError(t, db.Create(&model.PricingPlanPolicy{
		PricingPlanId: defaultPlan.Id, ActualAmount: 0,
		Type: constants.PricingPolicyTypeMonthly,
	}).Error)

	payload := checkoutCompletedEvent(t, "evt_act2", org, plan, policy, "sub_gone")
	require.NoError(t, l.HandleWebhookEvent(ctx, payload, "sig"))

	deletedPayload := webhookEvent(t, "evt_del", "customer.subscription.deleted", map[string]interface{}{
		"id":     "sub_gone",
		"status": "canceled",
	})
	require.NoError(t, l.HandleWebhookEvent(ctx, deletedPayload, "sig"))

	var stripeSub model.OrganizationSubscription
	require.NoError(t, db.Where("stripe_subscription_id = ?", "sub_gone").First(&stripeSub).Error)
	assert.Equal(t, constants.OrganizationSubscriptionStatusExpired, stripeSub.Status)

	// recovered onto the default free plan
	var recovered []*model.OrganizationSubscription
	require.NoError(t, db.Where("organization_id = ? AND pricing_plan_id = ? AND status = ?",
		org.Id, defaultPlan.Id, constants.OrganizationSubscriptionStatusActive).Find(&recovered).Error)
	assert.Len(t, recovered, 1)
}

func TestWebhook_SubscriptionUpdatedSyncsAutoRenew(t *testing.T) {
	l, _, db := setupStripePaymentLogic(t)
	ctx := context.Background()
	plan, policy := seedPaidPlan(t, db)
	org := seedOwnedOrg(t, db, 42)

	payload := checkoutCompletedEvent(t, "evt_act3", org, plan, policy, "sub_toggle")
	require.NoError(t, l.HandleWebhookEvent(ctx, payload, "sig"))

	updatedPayload := webhookEvent(t, "evt_upd", "customer.subscription.updated", map[string]interface{}{
		"id":                   "sub_toggle",
		"status":               "active",
		"cancel_at_period_end": true,
	})
	require.NoError(t, l.HandleWebhookEvent(ctx, updatedPayload, "sig"))

	var subscription model.OrganizationSubscription
	require.NoError(t, db.Where("stripe_subscription_id = ?", "sub_toggle").First(&subscription).Error)
	assert.False(t, subscription.IsAutoRenew)
	assert.Equal(t, constants.OrganizationSubscriptionStatusActive, subscription.Status)
}

func TestSetAutoRenew_TogglesGatewayAndLocalFlag(t *testing.T) {
	l, gateway, db := setupStripePaymentLogic(t)
	ctx := context.Background()
	plan, policy := seedPaidPlan(t, db)
	org := seedOwnedOrg(t, db, 42)

	require.NoError(t, db.Create(&model.OrganizationSubscription{
		Code: "SUBAR", OrganizationId: org.Id, PricingPlanId: plan.Id,
		PricingPlanPolicyId:  policy.Id,
		Status:               constants.OrganizationSubscriptionStatusActive,
		ExpireDate:           calculateExpireDate(constants.PricingPolicyTypeMonthly, time.Now()),
		IsAutoRenew:          true,
		Source:               constants.OrganizationSubscriptionSourceStripe,
		StripeSubscriptionId: "sub_ar",
	}).Error)

	require.NoError(t, l.SetAutoRenew(ctx, dto.StripeAutoRenewRequestDto{
		OrganizationId: org.Id, UserId: 42, Enabled: false,
	}))
	assert.True(t, gateway.cancelAtEnd["sub_ar"])

	var subscription model.OrganizationSubscription
	require.NoError(t, db.Where("code = ?", "SUBAR").First(&subscription).Error)
	assert.False(t, subscription.IsAutoRenew)
}

// --- cron / manual-flow coordination ---

func TestUpdateStatus_StripeAutoRenewGetsGraceWindow(t *testing.T) {
	l, _, db := setupStripePaymentLogic(t)
	ctx := context.Background()
	plan, policy := seedPaidPlan(t, db)
	org := seedOwnedOrg(t, db, 42)

	// expired yesterday, but auto-renewing via Stripe: within the 3-day grace
	// window the cron must neither expire nor mark it Expiring
	require.NoError(t, db.Create(&model.OrganizationSubscription{
		Code: "SUBGRACE", OrganizationId: org.Id, PricingPlanId: plan.Id,
		PricingPlanPolicyId:  policy.Id,
		Status:               constants.OrganizationSubscriptionStatusActive,
		ExpireDate:           &sql.NullTime{Time: time.Now().AddDate(0, 0, -1), Valid: true},
		IsAutoRenew:          true,
		Source:               constants.OrganizationSubscriptionSourceStripe,
		StripeSubscriptionId: "sub_grace",
	}).Error)

	require.NoError(t, l.subscriptionLogic.UpdateStatus(ctx, constants.UserTypeBorrower))

	var subscription model.OrganizationSubscription
	require.NoError(t, db.Where("code = ?", "SUBGRACE").First(&subscription).Error)
	assert.Equal(t, constants.OrganizationSubscriptionStatusActive, subscription.Status)
}

func TestUpdateStatus_StripeAutoRenewExpiresAndCancelsAfterGrace(t *testing.T) {
	l, gateway, db := setupStripePaymentLogic(t)
	ctx := context.Background()
	plan, policy := seedPaidPlan(t, db)
	org := seedOwnedOrg(t, db, 42)

	// expired beyond the grace window: cron force-expires and cancels remotely
	require.NoError(t, db.Create(&model.OrganizationSubscription{
		Code: "SUBLAPSE", OrganizationId: org.Id, PricingPlanId: plan.Id,
		PricingPlanPolicyId:  policy.Id,
		Status:               constants.OrganizationSubscriptionStatusActive,
		ExpireDate:           &sql.NullTime{Time: time.Now().AddDate(0, 0, -(stripeAutoRenewGraceDays + 1)), Valid: true},
		IsAutoRenew:          true,
		Source:               constants.OrganizationSubscriptionSourceStripe,
		StripeSubscriptionId: "sub_lapse",
	}).Error)

	require.NoError(t, l.subscriptionLogic.UpdateStatus(ctx, constants.UserTypeBorrower))

	var subscription model.OrganizationSubscription
	require.NoError(t, db.Where("code = ?", "SUBLAPSE").First(&subscription).Error)
	assert.Equal(t, constants.OrganizationSubscriptionStatusExpired, subscription.Status)
	assert.Contains(t, gateway.cancelled, "sub_lapse")
}

func TestAdminCreate_CancelsReplacedStripeSubscription(t *testing.T) {
	l, gateway, db := setupStripePaymentLogic(t)
	ctx := context.Background()
	plan, policy := seedPaidPlan(t, db)
	org := seedOwnedOrg(t, db, 42)

	// live Stripe subscription that the admin is about to replace manually
	require.NoError(t, db.Create(&model.OrganizationSubscription{
		Code: "SUBOLD", OrganizationId: org.Id, PricingPlanId: plan.Id,
		PricingPlanPolicyId:  policy.Id,
		Status:               constants.OrganizationSubscriptionStatusActive,
		ExpireDate:           calculateExpireDate(constants.PricingPolicyTypeMonthly, time.Now()),
		IsAutoRenew:          true,
		Source:               constants.OrganizationSubscriptionSourceStripe,
		StripeSubscriptionId: "sub_replaced",
	}).Error)

	_, err := l.subscriptionLogic.Create(ctx, dto.OrganizationSubscriptionCreateDto{
		OrganizationId:      org.Id,
		PricingPlanId:       plan.Id,
		PricingPlanPolicyId: policy.Id,
		Type:                constants.UserTypeBorrower,
	})
	require.NoError(t, err)
	assert.Contains(t, gateway.cancelled, "sub_replaced")

	var old model.OrganizationSubscription
	require.NoError(t, db.Where("code = ?", "SUBOLD").First(&old).Error)
	assert.Equal(t, constants.OrganizationSubscriptionStatusExpired, old.Status)
}
