package payment

import (
	"context"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/client"
	"github.com/stripe/stripe-go/v82/webhook"

	"zeebee-backend/common/errorx"
)

// Client implements Gateway on top of the official stripe-go SDK.
type Client struct {
	config Config
	api    *client.API
}

func NewClient(config Config) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &Client{
		config: config,
		api:    client.New(config.SecretKey, nil),
	}, nil
}

func (c *Client) CreateCustomer(ctx context.Context, p CustomerParams) (string, error) {
	params := &stripe.CustomerParams{
		Params: stripe.Params{Context: ctx},
		Email:  stripe.String(p.Email),
		Name:   stripe.String(p.Name),
	}
	params.AddMetadata("organization_id", strconv.FormatInt(p.OrganizationId, 10))
	customer, err := c.api.Customers.New(params)
	if err != nil {
		return "", errorx.Wrap(err, "stripe create customer failed")
	}
	return customer.ID, nil
}

func (c *Client) CreateProduct(ctx context.Context, p ProductParams) (string, error) {
	params := &stripe.ProductParams{
		Params: stripe.Params{Context: ctx},
		Name:   stripe.String(p.Name),
	}
	params.AddMetadata("pricing_plan_id", strconv.FormatInt(p.PricingPlanId, 10))
	product, err := c.api.Products.New(params)
	if err != nil {
		return "", errorx.Wrap(err, "stripe create product failed")
	}
	return product.ID, nil
}

func (c *Client) UpdateProductName(ctx context.Context, productId, name string) error {
	_, err := c.api.Products.Update(productId, &stripe.ProductParams{
		Params: stripe.Params{Context: ctx},
		Name:   stripe.String(name),
	})
	return errorx.Wrap(err, "stripe update product failed")
}

func (c *Client) CreatePrice(ctx context.Context, p PriceParams) (string, error) {
	params := &stripe.PriceParams{
		Params:     stripe.Params{Context: ctx},
		Product:    stripe.String(p.ProductId),
		UnitAmount: stripe.Int64(p.UnitAmountCents),
		Currency:   stripe.String(c.config.Currency),
		Recurring: &stripe.PriceRecurringParams{
			Interval:      stripe.String(p.Interval),
			IntervalCount: stripe.Int64(p.IntervalCount),
		},
	}
	params.AddMetadata("pricing_plan_policy_id", strconv.FormatInt(p.PolicyId, 10))
	price, err := c.api.Prices.New(params)
	if err != nil {
		return "", errorx.Wrap(err, "stripe create price failed")
	}
	return price.ID, nil
}

func (c *Client) ArchivePrice(ctx context.Context, priceId string) error {
	_, err := c.api.Prices.Update(priceId, &stripe.PriceParams{
		Params: stripe.Params{Context: ctx},
		Active: stripe.Bool(false),
	})
	return errorx.Wrap(err, "stripe archive price failed")
}

func (c *Client) CreateCheckoutSession(ctx context.Context, p CheckoutParams) (string, error) {
	params := &stripe.CheckoutSessionParams{
		Params:   stripe.Params{Context: ctx},
		Mode:     stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		Customer: stripe.String(p.CustomerId),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(p.PriceId),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(c.config.SuccessUrl),
		CancelURL:  stripe.String(c.config.CancelUrl),
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: p.Metadata,
		},
	}
	for key, value := range p.Metadata {
		params.AddMetadata(key, value)
	}
	session, err := c.api.CheckoutSessions.New(params)
	if err != nil {
		return "", errorx.Wrap(err, "stripe create checkout session failed")
	}
	return session.URL, nil
}

func (c *Client) SetCancelAtPeriodEnd(ctx context.Context, subscriptionId string, cancel bool) error {
	_, err := c.api.Subscriptions.Update(subscriptionId, &stripe.SubscriptionParams{
		Params:            stripe.Params{Context: ctx},
		CancelAtPeriodEnd: stripe.Bool(cancel),
	})
	return errorx.Wrap(err, "stripe update subscription auto-renew failed")
}

func (c *Client) CancelSubscription(ctx context.Context, subscriptionId string) error {
	_, err := c.api.Subscriptions.Cancel(subscriptionId, &stripe.SubscriptionCancelParams{
		Params: stripe.Params{Context: ctx},
	})
	return errorx.Wrap(err, "stripe cancel subscription failed")
}

func (c *Client) GetSubscription(ctx context.Context, subscriptionId string) (*SubscriptionInfo, error) {
	subscription, err := c.api.Subscriptions.Get(subscriptionId, &stripe.SubscriptionParams{
		Params: stripe.Params{Context: ctx},
	})
	if err != nil {
		return nil, errorx.Wrap(err, "stripe get subscription failed")
	}
	return NormalizeSubscription(subscription), nil
}

func (c *Client) ConstructEvent(payload []byte, sigHeader string) (stripe.Event, error) {
	event, err := webhook.ConstructEventWithOptions(payload, sigHeader, c.config.WebhookSecret,
		webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true})
	if err != nil {
		return event, errorx.Wrap(err, "stripe webhook signature verification failed")
	}
	return event, nil
}

// NormalizeSubscription flattens a Stripe Subscription into SubscriptionInfo.
// Since API version 2025-03-31 current_period_end lives on subscription items.
func NormalizeSubscription(subscription *stripe.Subscription) *SubscriptionInfo {
	info := &SubscriptionInfo{
		Id:                subscription.ID,
		Status:            string(subscription.Status),
		CancelAtPeriodEnd: subscription.CancelAtPeriodEnd,
	}
	if subscription.Items != nil && len(subscription.Items.Data) > 0 {
		info.CurrentPeriodEnd = time.Unix(subscription.Items.Data[0].CurrentPeriodEnd, 0).UTC()
	}
	return info
}
