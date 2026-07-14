package payment

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

func validConfig() Config {
	return Config{
		SecretKey:     "sk_test_xxx",
		WebhookSecret: "whsec_test_secret",
		SuccessUrl:    "https://example.com/billing/success",
		CancelUrl:     "https://example.com/billing/cancel",
		Currency:      "usd",
	}
}

func TestConfigValidate(t *testing.T) {
	cfg := validConfig()
	require.NoError(t, cfg.Validate())

	missingKey := validConfig()
	missingKey.SecretKey = ""
	assert.Error(t, missingKey.Validate())
	assert.False(t, missingKey.Enabled())

	missingSuccess := validConfig()
	missingSuccess.SuccessUrl = ""
	assert.Error(t, missingSuccess.Validate())

	missingCancel := validConfig()
	missingCancel.CancelUrl = ""
	assert.Error(t, missingCancel.Validate())

	defaultCurrency := validConfig()
	defaultCurrency.Currency = ""
	require.NoError(t, defaultCurrency.Validate())
	assert.Equal(t, "usd", defaultCurrency.Currency)
}

func TestConstructEventVerifiesSignature(t *testing.T) {
	cfg := validConfig()
	stripeClient, err := NewClient(cfg)
	require.NoError(t, err)

	payload := []byte(`{"id":"evt_test_1","object":"event","type":"invoice.paid","api_version":"2025-04-30.basil"}`)
	signedPayload := webhook.GenerateTestSignedPayload(&webhook.UnsignedPayload{
		Payload:   payload,
		Secret:    cfg.WebhookSecret,
		Timestamp: time.Now(),
	})

	event, err := stripeClient.ConstructEvent(signedPayload.Payload, signedPayload.Header)
	require.NoError(t, err)
	assert.Equal(t, "evt_test_1", event.ID)
	assert.Equal(t, stripe.EventType("invoice.paid"), event.Type)

	// tampered payload must be rejected
	_, err = stripeClient.ConstructEvent([]byte(`{"id":"evt_forged"}`), signedPayload.Header)
	assert.Error(t, err)
}

func TestNormalizeSubscription(t *testing.T) {
	periodEnd := time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)
	info := NormalizeSubscription(&stripe.Subscription{
		ID:                "sub_123",
		Status:            stripe.SubscriptionStatusActive,
		CancelAtPeriodEnd: true,
		Items: &stripe.SubscriptionItemList{
			Data: []*stripe.SubscriptionItem{
				{CurrentPeriodEnd: periodEnd.Unix()},
			},
		},
	})
	assert.Equal(t, "sub_123", info.Id)
	assert.Equal(t, "active", info.Status)
	assert.True(t, info.CancelAtPeriodEnd)
	assert.Equal(t, periodEnd, info.CurrentPeriodEnd)

	// no items -> zero time, no panic
	empty := NormalizeSubscription(&stripe.Subscription{ID: "sub_empty"})
	assert.True(t, empty.CurrentPeriodEnd.IsZero())
}
