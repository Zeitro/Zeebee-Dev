package constants

import (
	"zeebee-backend/common/enums"
)

// OrganizationSubscriptionStatus organization subscription status
type OrganizationSubscriptionStatus string

const OrganizationSubscriptionStatusName = "organization_subscription_status"

const (
	OrganizationSubscriptionStatusTrialing  OrganizationSubscriptionStatus = "Trialing"
	OrganizationSubscriptionStatusActive    OrganizationSubscriptionStatus = "Active"
	OrganizationSubscriptionStatusExpiring  OrganizationSubscriptionStatus = "Expiring"
	OrganizationSubscriptionStatusExpired   OrganizationSubscriptionStatus = "Expired"
	OrganizationSubscriptionStatusSuspended OrganizationSubscriptionStatus = "Suspended"
)

// EffectiveStatuses are the subscription statuses that grant pricing-plan
// resource access. Only Active and Trialing are effective; Suspended loses
// access, Expiring is a reminder-only state, Expired is terminal.
var EffectiveStatuses = []OrganizationSubscriptionStatus{
	OrganizationSubscriptionStatusActive,
	OrganizationSubscriptionStatusTrialing,
	OrganizationSubscriptionStatusExpiring,
}

// IsEffective reports whether the status grants plan resource access.
// This is the single source of truth for "effective" semantics across logic,
// repo, and auth layers.
func (s OrganizationSubscriptionStatus) IsEffective() bool {
	return s == OrganizationSubscriptionStatusActive ||
		s == OrganizationSubscriptionStatusTrialing ||
		s == OrganizationSubscriptionStatusExpiring
}

func (s OrganizationSubscriptionStatus) Validate() error {
	ranges := []OrganizationSubscriptionStatus{
		OrganizationSubscriptionStatusTrialing,
		OrganizationSubscriptionStatusActive,
		OrganizationSubscriptionStatusExpiring,
		OrganizationSubscriptionStatusExpired,
		OrganizationSubscriptionStatusSuspended,
	}
	return validateConstant(s, ranges, OrganizationSubscriptionStatusName)
}

func init() {

	enums.Register(OrganizationSubscriptionStatusName, []enums.EnumOpt{
		{Label: string(OrganizationSubscriptionStatusTrialing), Value: string(OrganizationSubscriptionStatusTrialing)},
		{Label: string(OrganizationSubscriptionStatusActive), Value: string(OrganizationSubscriptionStatusActive)},
		{Label: string(OrganizationSubscriptionStatusExpiring), Value: string(OrganizationSubscriptionStatusExpiring)},
		{Label: string(OrganizationSubscriptionStatusExpired), Value: string(OrganizationSubscriptionStatusExpired)},
		{Label: string(OrganizationSubscriptionStatusSuspended), Value: string(OrganizationSubscriptionStatusSuspended)},
	})
}
