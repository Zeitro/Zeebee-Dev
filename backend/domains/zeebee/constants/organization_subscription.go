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
