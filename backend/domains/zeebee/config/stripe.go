package config

// StripeConf is the shared Stripe configuration shape for service YAML files.
// All fields are optional: an empty SecretKey means payment is disabled for
// the service and the payment gateway is not constructed.
type StripeConf struct {
	SecretKey     string `json:",optional"`
	WebhookSecret string `json:",optional"`
	SuccessUrl    string `json:",optional"`
	CancelUrl     string `json:",optional"`
	Currency      string `json:",default=usd"`
}
