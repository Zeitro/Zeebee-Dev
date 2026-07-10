package payment

import (
	"context"
	"time"

	"github.com/stripe/stripe-go/v82"
)

// CustomerParams are the inputs for creating a Stripe Customer.
type CustomerParams struct {
	Email          string
	Name           string
	OrganizationId int64
}

// ProductParams are the inputs for creating a Stripe Product.
type ProductParams struct {
	Name          string
	PricingPlanId int64
}

// PriceParams are the inputs for creating a recurring Stripe Price.
type PriceParams struct {
	ProductId       string
	UnitAmountCents int64
	Interval        string // month | year
	IntervalCount   int64
	PolicyId        int64
}

// CheckoutParams are the inputs for creating a Stripe Checkout Session in
// subscription mode. Metadata is attached to both the session and the
// subscription it creates, so webhook handlers can resolve the local entities.
type CheckoutParams struct {
	CustomerId string
	PriceId    string
	Metadata   map[string]string
}

// SubscriptionInfo is a normalized view of a Stripe Subscription.
type SubscriptionInfo struct {
	Id                string
	Status            string
	CancelAtPeriodEnd bool
	CurrentPeriodEnd  time.Time
}

// NewGateway builds a Gateway from config, or nil when Stripe is disabled
// (empty secret key). An invalid non-empty config returns an error.
func NewGateway(config Config) (Gateway, error) {
	if !config.Enabled() {
		return nil, nil
	}
	return NewClient(config)
}

// Gateway abstracts the Stripe API so domain logic stays testable with fakes.
// Ensure/lookup semantics (check DB, create, write back) live in the logic
// layer; the gateway only wraps remote calls.
type Gateway interface {
	CreateCustomer(ctx context.Context, params CustomerParams) (customerId string, err error)
	CreateProduct(ctx context.Context, params ProductParams) (productId string, err error)
	UpdateProductName(ctx context.Context, productId, name string) error
	CreatePrice(ctx context.Context, params PriceParams) (priceId string, err error)
	ArchivePrice(ctx context.Context, priceId string) error
	CreateCheckoutSession(ctx context.Context, params CheckoutParams) (checkoutUrl string, err error)
	SetCancelAtPeriodEnd(ctx context.Context, subscriptionId string, cancel bool) error
	CancelSubscription(ctx context.Context, subscriptionId string) error
	GetSubscription(ctx context.Context, subscriptionId string) (*SubscriptionInfo, error)
	// ConstructEvent verifies the webhook signature and parses the event.
	ConstructEvent(payload []byte, sigHeader string) (stripe.Event, error)
}
