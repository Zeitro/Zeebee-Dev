package model

import (
	"zeebee-backend/domains/zeebee/constants"
)

// StripeEvent is the idempotency ledger for received Stripe webhook events.
// The unique index on event_id guarantees each event is processed at most once.
type StripeEvent struct {
	Base
	EventId   string                      `gorm:"column:event_id;type:varchar(64);not null;default:'';uniqueIndex;<-:create" json:"event_id"`
	EventType string                      `gorm:"column:event_type;type:varchar(64);not null;default:''" json:"event_type"`
	Status    constants.StripeEventStatus `gorm:"column:status;type:varchar(20);not null;default:'Processing'" json:"status"`
	Payload   string                      `gorm:"column:payload;type:jsonb;not null;default:'{}'" json:"payload"`
	Error     string                      `gorm:"column:error;type:text;not null;default:''" json:"error"`
}

func (StripeEvent) TableName() string {
	return "stripe_events"
}
