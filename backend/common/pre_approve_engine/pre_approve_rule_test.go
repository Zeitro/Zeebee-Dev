package pre_approve_engine

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"zeebee-backend/common/constants"
)

func TestPreApproveRule_CalPreApprovedLoanAmount(t *testing.T) {
	// Create a test rule with sample data
	rule := &PreApproveRule{
		ExperienceLevel:      ExperienceLevelExperienced,
		ExtraFactor:          0.05,
		MinLargeUnitCount:    5,
		MaxUnitCount:         100,
		MinLoanTerm:          6,
		MaxLoanTerm:          24,
		MinHeavyRehabPercent: 0.3,
		BaseRate:             0.07,
		PreApproveAmountRange: PreApproveAmountRange{
			StandardAmountRange:  AmountRange{Min: 50000, Max: 2000000},
			LargeUnitAmountRange: AmountRange{Min: 100000, Max: 5000000},
		},
		// ExperienceLevelRule: ExperienceLevelRule{
		// 	MaxTransactionHistoryCountForNoExperience: 0,
		// 	MaxTransactionHistoryCountForExperienced:  5,
		// },
		LtcBaseRules: []BaseRateRule{
			{Min: 0, Max: 0.7, BaseRate: 0.06},
			{Min: 0.7, Max: 1.0, BaseRate: 0.07},
		},
		LtarvBaseRules: []BaseRateRule{
			{Min: 0, Max: 0.7, BaseRate: 0.06},
			{Min: 0.7, Max: 1.0, BaseRate: 0.07},
		},
		AdjusterRule: AdjusterRule{
			FicoScore: []RangeNumberRule{
				{Min: 0, Max: 600, Value: 0.02},
				{Min: 600, Max: 700, Value: 0.01},
				{Min: 700, Max: 850, Value: 0},
			},
			GUCExperienceLevel: map[ExperienceLevel]float64{
				ExperienceLevelNoExperienced: 0.02,
				ExperienceLevelExperienced:   0.01,
				ExperienceLevelInstitutional: 0,
			},
			LtcAndLoanAmountRules: []LtcAndLoanAmountRule{
				{LtcRange: RangeNumber{Min: 0, Max: 0.7}, LoanAmountRange: RangeNumber{Min: 0, Max: 500000}, Value: 0.005},
			},
			LtarvAndLoanAmountRules: []LtarvAndLoanAmountRule{
				{LtarvRange: RangeNumber{Min: 0, Max: 0.7}, LoanAmountRange: RangeNumber{Min: 0, Max: 500000}, Value: 0.005},
			},
			LtcAndLoanPurposeRules: []LtcAndLoanPurposeRule{
				{LtcRange: RangeNumber{Min: 0, Max: 0.7}, LoanPurpose: constants.LoanPurposePurchase, Value: 0.005},
			},
			LtarvAndLoanPurposeRules: []LtarvAndLoanPurposeRule{
				{LtarvRange: RangeNumber{Min: 0, Max: 0.7}, LoanPurpose: constants.LoanPurposePurchase, Value: 0.005},
			},
		},
		ServiceFeeConfig: ServiceFeeConfig{
			OriginationFeePercent:  0.01,
			UnderwritingFeePercent: 0.005,
		},
		PriceCapRuleMap: map[constants.LoanType]*LoanTypeRule{
			constants.LoanTypeRehab: {
				LoanType: constants.LoanTypeRehab,
				GeneralRules: &RehabLoanTypeRules{
					Heavy: RehabLoanTypeRule{
						GeneralRule: &RuleCap{
							Ltaiv: &Values{Standard: 0.8, LargeUnit: 0.75},
							Ltc:   &Values{Standard: 0.9, LargeUnit: 0.85},
							Ltarv: &Values{Standard: 0.75, LargeUnit: 0.7},
						},
						HighValueRule: &RuleCap{
							Ltaiv: &Values{Standard: 0.75, LargeUnit: 0.7},
							Ltc:   &Values{Standard: 0.85, LargeUnit: 0.8},
							Ltarv: &Values{Standard: 0.7, LargeUnit: 0.65},
						},
					},
					Light: RehabLoanTypeRule{
						GeneralRule: &RuleCap{
							Ltaiv: &Values{Standard: 0.85, LargeUnit: 0.8},
							Ltc:   &Values{Standard: 0.95, LargeUnit: 0.9},
							Ltarv: &Values{Standard: 0.8, LargeUnit: 0.75},
						},
						HighValueRule: &RuleCap{
							Ltaiv: &Values{Standard: 0.8, LargeUnit: 0.75},
							Ltc:   &Values{Standard: 0.9, LargeUnit: 0.85},
							Ltarv: &Values{Standard: 0.75, LargeUnit: 0.7},
						},
					},
				},
			},
			constants.LoanTypeConstruction: {
				LoanType: constants.LoanTypeConstruction,
				GeneralRules: &GeneralLoanTypeRules{
					GeneralRules: &RuleCap{
						Ltaiv: &Values{Standard: 0.8, LargeUnit: 0.75},
						Ltc:   &Values{Standard: 0.9, LargeUnit: 0.85},
						Ltarv: &Values{Standard: 0.75, LargeUnit: 0.7},
					},
					HighValueRules: &RuleCap{
						Ltaiv: &Values{Standard: 0.75, LargeUnit: 0.7},
						Ltc:   &Values{Standard: 0.85, LargeUnit: 0.8},
						Ltarv: &Values{Standard: 0.7, LargeUnit: 0.65},
					},
				},
			},
		},
		HighValueStates: []*HighValueState{},
	}

	t.Run("approved loan for rehab", func(t *testing.T) {
		input := InputLoanInfo{
			LoanType:                constants.LoanTypeRehab,
			LoanPurpose:             constants.LoanPurposePurchase,
			FicoScore:               720,
			TransactionHistoryCount: 3,
			UnitCount:               2,
			LoanTerm:                12,
			EstCurrentPropertyValue: 500000,
			EstAfterRepairValue:     600000,
			OriginalBudgetAmount:    100000,
			PurchaseAmount:          500000,
		}

		result, err := rule.CalPreApprovedLoanAmount(input)
		assert.NoError(t, err)
		assert.True(t, result.IsApproved)
		assert.Greater(t, result.TotalAmount, 0.0)
		assert.Equal(t, ExperienceLevelExperienced, result.ExperienceLevel)
	})

	t.Run("approved loan for construction", func(t *testing.T) {
		input := InputLoanInfo{
			LoanType:                constants.LoanTypeConstruction,
			LoanPurpose:             constants.LoanPurposePurchase,
			FicoScore:               720,
			TransactionHistoryCount: 3,
			UnitCount:               2,
			LoanTerm:                12,
			EstCurrentPropertyValue: 500000,
			EstAfterRepairValue:     600000,
			OriginalBudgetAmount:    100000,
			PurchaseAmount:          500000,
			InterestReserveMonths:   6,
		}

		result, err := rule.CalPreApprovedLoanAmount(input)
		assert.NoError(t, err)
		assert.True(t, result.IsApproved)
		assert.Greater(t, result.TotalAmount, 0.0)
		assert.GreaterOrEqual(t, result.InterestReservedAmount, 0.0)
	})

	t.Run("rejected loan due to high unit count", func(t *testing.T) {
		input := InputLoanInfo{
			LoanType:                constants.LoanTypeRehab,
			LoanPurpose:             constants.LoanPurposePurchase,
			FicoScore:               720,
			TransactionHistoryCount: 3,
			UnitCount:               150, // Exceeds MaxUnitCount
			LoanTerm:                12,
			EstCurrentPropertyValue: 500000,
			EstAfterRepairValue:     600000,
			OriginalBudgetAmount:    100000,
			PurchaseAmount:          500000,
		}

		result, err := rule.CalPreApprovedLoanAmount(input)
		assert.NoError(t, err)
		assert.False(t, result.IsApproved)
	})

	t.Run("rejected loan due to budget exceeding property value", func(t *testing.T) {
		input := InputLoanInfo{
			LoanType:                constants.LoanTypeRehab,
			LoanPurpose:             constants.LoanPurposePurchase,
			FicoScore:               720,
			TransactionHistoryCount: 3,
			UnitCount:               2,
			LoanTerm:                12,
			EstCurrentPropertyValue: 500000,
			EstAfterRepairValue:     600000,
			OriginalBudgetAmount:    600000, // Exceeds property value
			PurchaseAmount:          500000,
		}

		result, err := rule.CalPreApprovedLoanAmount(input)
		assert.NoError(t, err)
		assert.False(t, result.IsApproved)
	})

	t.Run("rejected loan due to invalid loan type", func(t *testing.T) {
		input := InputLoanInfo{
			LoanType:                "invalid_loan_type",
			LoanPurpose:             constants.LoanPurposePurchase,
			FicoScore:               720,
			TransactionHistoryCount: 3,
			UnitCount:               2,
			LoanTerm:                12,
			EstCurrentPropertyValue: 500000,
			EstAfterRepairValue:     600000,
			OriginalBudgetAmount:    100000,
			PurchaseAmount:          500000,
		}

		_, err := rule.CalPreApprovedLoanAmount(input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid loan type")
	})

	t.Run("rejected loan due to zero property value", func(t *testing.T) {
		input := InputLoanInfo{
			LoanType:                constants.LoanTypeRehab,
			LoanPurpose:             constants.LoanPurposePurchase,
			FicoScore:               720,
			TransactionHistoryCount: 3,
			UnitCount:               2,
			LoanTerm:                12,
			EstCurrentPropertyValue: 0, // Invalid value
			EstAfterRepairValue:     600000,
			OriginalBudgetAmount:    100000,
			PurchaseAmount:          500000,
		}

		result, err := rule.CalPreApprovedLoanAmount(input)
		assert.NoError(t, err)
		assert.False(t, result.IsApproved)
	})

	t.Run("large unit loan calculation", func(t *testing.T) {
		input := InputLoanInfo{
			LoanType:                constants.LoanTypeRehab,
			LoanPurpose:             constants.LoanPurposePurchase,
			FicoScore:               720,
			TransactionHistoryCount: 3,
			UnitCount:               10, // Large unit
			LoanTerm:                12,
			EstCurrentPropertyValue: 500000,
			EstAfterRepairValue:     600000,
			OriginalBudgetAmount:    100000,
			PurchaseAmount:          500000,
		}

		result, err := rule.CalPreApprovedLoanAmount(input)
		assert.NoError(t, err)
		assert.True(t, result.IsApproved)
	})
}
