package constants

import (
	"testing"

	"zeebee-backend/common/enums"
)

func TestPricingPlanType_Validate(t *testing.T) {
	tests := []struct {
		name     string
		planType PricingPlanType
		wantErr  bool
	}{
		{"Standard", PricingPlanTypeStandard, false},
		{"Custom", PricingPlanTypeCustom, false},
		{"Invalid", PricingPlanType("Invalid"), true},
		{"Empty", PricingPlanType(""), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.planType.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PricingPlanType.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPricingPolicyType_Validate(t *testing.T) {
	tests := []struct {
		name       string
		policyType PricingPolicyType
		wantErr    bool
	}{
		{"Monthly", PricingPolicyTypeMonthly, false},
		{"Quarterly", PricingPolicyTypeQuarterly, false},
		{"Yearly", PricingPolicyTypeYearly, false},
		{"Invalid", PricingPolicyType("Invalid"), true},
		{"Empty", PricingPolicyType(""), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.policyType.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PricingPolicyType.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrganizationSubscriptionStatus_Validate(t *testing.T) {
	tests := []struct {
		name    string
		status  OrganizationSubscriptionStatus
		wantErr bool
	}{
		{"Trialing", OrganizationSubscriptionStatusTrialing, false},
		{"Active", OrganizationSubscriptionStatusActive, false},
		{"Expiring", OrganizationSubscriptionStatusExpiring, false},
		{"Expired", OrganizationSubscriptionStatusExpired, false},
		{"Suspended", OrganizationSubscriptionStatusSuspended, false},
		{"Invalid", OrganizationSubscriptionStatus("Invalid"), true},
		{"Empty", OrganizationSubscriptionStatus(""), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.status.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OrganizationSubscriptionStatus.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPricingPlanEnumsRegistered(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		count  int
		values []string
	}{
		{"pricing_plan_type", "pricing_plan_type", 2, []string{"Standard", "Custom"}},
		{"pricing_policy_type", "pricing_policy_type", 3, []string{"Monthly", "Quarterly", "Yearly"}},
		{"organization_subscription_status", "organization_subscription_status", 5, []string{"Trialing", "Active", "Expiring", "Expired", "Suspended"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := enums.Get(tt.key)
			if len(opts) != tt.count {
				t.Errorf("enums.Get(%q) returned %d opts, want %d", tt.key, len(opts), tt.count)
			}
			for i, v := range tt.values {
				if i < len(opts) && opts[i].Value != v {
					t.Errorf("enums.Get(%q)[%d].Value = %q, want %q", tt.key, i, opts[i].Value, v)
				}
			}
		})
	}
}
