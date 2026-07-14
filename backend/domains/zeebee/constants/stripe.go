package constants

// StripeEventStatus is the processing status of a received Stripe webhook event.
type StripeEventStatus string

const (
	StripeEventStatusProcessing StripeEventStatus = "Processing"
	StripeEventStatusProcessed  StripeEventStatus = "Processed"
	StripeEventStatusFailed     StripeEventStatus = "Failed"
	StripeEventStatusSkipped    StripeEventStatus = "Skipped"
)

// PaymentTransactionStatus is the status of a payment transaction synced from Stripe.
type PaymentTransactionStatus string

const (
	PaymentTransactionStatusPaid     PaymentTransactionStatus = "Paid"
	PaymentTransactionStatusFailed   PaymentTransactionStatus = "Failed"
	PaymentTransactionStatusRefunded PaymentTransactionStatus = "Refunded"
)
