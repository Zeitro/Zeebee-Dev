package model

import (
	"database/sql"

	"zeebee-backend/domains/zeebee/constants"
)

// PaymentTransaction is a payment record synced from Stripe invoices.
// stripe_invoice_id is unique so replayed webhook events stay idempotent.
type PaymentTransaction struct {
	Base
	OrganizationId             int64                              `gorm:"column:organization_id;type:int4;not null;default:0" json:"organization_id"`
	OrganizationSubscriptionId int64                              `gorm:"column:organization_subscription_id;type:int4;not null;default:0" json:"organization_subscription_id"`
	StripeInvoiceId            string                             `gorm:"column:stripe_invoice_id;type:varchar(64);not null;default:'';uniqueIndex" json:"stripe_invoice_id"`
	StripeSubscriptionId       string                             `gorm:"column:stripe_subscription_id;type:varchar(64);not null;default:''" json:"stripe_subscription_id"`
	Amount                     float64                            `gorm:"column:amount;type:decimal(10,2);not null;default:0" json:"amount"`
	Currency                   string                             `gorm:"column:currency;type:varchar(10);not null;default:'usd'" json:"currency"`
	Status                     constants.PaymentTransactionStatus `gorm:"column:status;type:varchar(20);not null;default:''" json:"status"`
	PeriodStart                *sql.NullTime                      `gorm:"column:period_start;type:date" json:"period_start"`
	PeriodEnd                  *sql.NullTime                      `gorm:"column:period_end;type:date" json:"period_end"`
	PaidAt                     *sql.NullTime                      `gorm:"column:paid_at;type:timestamptz" json:"paid_at"`
}

func (PaymentTransaction) TableName() string {
	return "payment_transactions"
}
