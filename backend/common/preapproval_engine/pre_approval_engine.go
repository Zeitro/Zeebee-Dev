package preapproval_engine

import (
	"math"

	"zeebee-backend/common/errorx"

	"zeebee-backend/domains/zeebee/constants"
)

// ExperienceLevel represents the borrower's experience level
type ExperienceLevel int

const (
	// ExperienceLevelNoExperienced represents a borrower with no experience
	ExperienceLevelNoExperienced ExperienceLevel = iota
	// ExperienceLevelExperienced represents a borrower with some experience
	ExperienceLevelExperienced
	// ExperienceLevelInstitutional represents an institutional borrower
	ExperienceLevelInstitutional
)

// InputLoanInfo represents the input loan information for calculating pre-approved loan amount
type InputLoanInfo struct {
	LoanType                constants.LoanType    `json:"loan_type"`
	LoanPurpose             constants.LoanPurpose `json:"loan_purpose"`
	FicoScore               int                   `json:"fico_score"`
	TransactionHistoryCount int                   `json:"transaction_history_count"`
	Unit                    int                   `json:"unit"`
	LoanTerm                int                   `json:"loan_term"`
	EstCurrentPropertyValue float64               `json:"est_current_property_value"`
	EstAfterRepairValue     float64               `json:"est_after_repair_value"`
	TotalAmount             float64               `json:"total_amount"`
	OriginalBudgetAmount    float64               `json:"original_budget_amount"`
	PurchaseAmount          float64               `json:"purchase_amount"`
	InitialAmount           float64               `json:"initial_amount"`
	HoldbackAmount          float64               `json:"holdback_amount"`
	PrincipalAmount         float64               `json:"principal_amount"`
	State                   string                `json:"state"`
	City                    string                `json:"city"`
	HasAssignmentFee        bool                  `json:"has_assignment_fee"`
	AssignmentFeeAmount     float64               `json:"assignment_fee_amount"`
	InterestReserveMonths   int                   `json:"interest_reserve_months"`
	IsHighValueZone         bool                  `json:"is_high_value_zone"`
}

// CalculatePreApprovalLoanAmountResult represents the result of calculating pre-approved loan amount
type CalculatePreApprovalLoanAmountResult struct {
	RequestId              string          `json:"request_id"`
	LenderCompanyName      string          `json:"lender_company_name"`
	LenderCompanyId        int64           `json:"lender_company_id"`
	IsApproved             bool            `json:"is_approved"`
	RejectReason           string          `json:"reject_reason"`
	TotalAmount            float64         `json:"total_amount"`
	InterestReservedAmount float64         `json:"interest_reserved_amount"`
	Rate                   float64         `json:"rate"`
	MonthlyPayment         float64         `json:"monthly_payment"`
	CashToClose            float64         `json:"cash_to_close"`
	InitialAmount          float64         `json:"initial_amount"`
	HoldbackAmount         float64         `json:"holdback_amount"`
	LoanTerm               int             `json:"loan_term"`
	Ltv                    float32         `json:"ltv"`
	Ltc                    float32         `json:"ltc"`
	Ltarv                  float32         `json:"ltarv"`
	ExperienceLevel        ExperienceLevel `json:"experience_level"`
}

// PreApprovalEngine represents the pre-approval engine
type PreApprovalEngine struct {
	LenderCompanyName string                               `json:"lender_company_name"`
	LenderCompanyId   int64                                `json:"lender_company_id"`
	PreApprovalRules  map[ExperienceLevel]*PreApprovalRule `json:"pre_approve_rules"`
	// ExperienceLevelRule defines rules for determining experience level
	ExperienceLevelRule ExperienceLevelRule `json:"experience_level_rule"`
}

// CalPreApprovaldLoanAmount calculates and returns the pre-approved loan amount
func (e *PreApprovalEngine) CalPreApprovalLoanAmount(input InputLoanInfo) (result *CalculatePreApprovalLoanAmountResult, err error) {
	experienceLevel := e.ExperienceLevelRule.getExperienceLevel(input.TransactionHistoryCount)
	rule, ok := e.PreApprovalRules[experienceLevel]
	if !ok {
		return nil, errorx.WithStack(errorx.New("No Rule Found For Experience Level"))
	}
	result, err = rule.CalPreApprovalLoanAmount(input)
	if err != nil {
		return nil, err
	}
	result.LenderCompanyId = e.LenderCompanyId
	result.LenderCompanyName = e.LenderCompanyName
	return result, nil
}

func GetDefaultPreApprovalEngines() map[string]*PreApprovalEngine {
	zeebeeMinLargeUnitCount := 5
	zeebeeMaxUnitCount := 30
	zeebeeMinLoanTerm := 12
	zeebeeMaxLoanTerm := 24
	zeebeeMinHeavyRehabPercent := 0.3
	zeebeeExtraFactor := 0.05
	zeebeePreApprovalAmountRange := PreApprovalAmountRange{
		StandardAmountRange: AmountRange{
			Min: 75_000,
			Max: 3_000_000,
		},
		LargeUnitAmountRange: AmountRange{
			Min: 75_000,
			Max: 5_000_000,
		},
	}
	zeebeeAdjusterRule := AdjusterRule{
		FicoScore: []RangeNumberRule{
			{
				Min:   300,
				Max:   701,
				Value: 0.0025,
			},
			{
				Min:   741,
				Max:   851,
				Value: -0.0025,
			},
		},
		GUCExperienceLevel: map[ExperienceLevel]float64{
			ExperienceLevelInstitutional: 0.005,
			ExperienceLevelExperienced:   0.0075,
		},
		LtcAndLoanAmountRules: []LtcAndLoanAmountRule{
			{
				LtcRange: RangeNumber{
					Min: 0.75,
					Max: 0.86,
				},
				LoanAmountRange: RangeNumber{
					Min: 3_000_000,
					Max: math.MaxFloat64,
				},
				Value: 0.0025,
			},
		},
		LtcAndLoanPurposeRules: []LtcAndLoanPurposeRule{
			{
				LtcRange: RangeNumber{
					Min: 0.70,
					Max: 0.81,
				},
				LoanPurpose: constants.LoanPurposeRateTerm,
				Value:       0.0025,
			},
			{
				LtcRange: RangeNumber{
					Min: 0.70,
					Max: 0.81,
				},
				LoanPurpose: constants.LoanPurposeCashOut,
				Value:       0.0025,
			},
		},
		LtarvAndLoanAmountRules: []LtarvAndLoanAmountRule{
			{
				LtarvRange: RangeNumber{
					Min: 0.66,
					Max: 0.76,
				},
				LoanAmountRange: RangeNumber{
					Min: 3_000_000,
					Max: math.MaxFloat64,
				},
				Value: 0.0025,
			},
		},
		LtarvAndLoanPurposeRules: []LtarvAndLoanPurposeRule{
			{
				LtarvRange: RangeNumber{
					Min: 0.6,
					Max: 0.71,
				},
				LoanPurpose: constants.LoanPurposeRateTerm,
				Value:       0.0025,
			},
			{
				LtarvRange: RangeNumber{
					Min: 0.6,
					Max: 0.71,
				},
				LoanPurpose: constants.LoanPurposeCashOut,
				Value:       0.0025,
			},
		},
	}
	zeebeeServiceFeeConfig := ServiceFeeConfig{
		OriginationFeePercent:  0.01,
		UnderwritingFeePercent: 0.0025,
	}
	zeebeeHighValueStates := []*HighValueState{}

	engines := make(map[string]*PreApprovalEngine)
	engines["zeebee"] = &PreApprovalEngine{
		LenderCompanyName: "ZeeBee",
		LenderCompanyId:   1,
		ExperienceLevelRule: ExperienceLevelRule{
			MaxTransactionHistoryCountForExperienced:  9,
			MaxTransactionHistoryCountForNoExperience: 2,
		},
		PreApprovalRules: map[ExperienceLevel]*PreApprovalRule{
			// No Experienced
			ExperienceLevelNoExperienced: &PreApprovalRule{
				ExperienceLevel: ExperienceLevelNoExperienced,
				ExtraFactor:     zeebeeExtraFactor,
				PriceCapRuleMap: map[constants.LoanType]*LoanTypeRule{
					constants.LoanTypeBridge: {
						LoanType: constants.LoanTypeBridge,
						GeneralRules: &GeneralLoanTypeRules{
							GeneralRules: &RuleCap{
								Ltaiv: &Values{
									Standard:  0.7,
									LargeUnit: 0.7,
								},
								Ltc: &Values{
									Standard:  0.7,
									LargeUnit: 0.7,
								},
							},
						},
					},
					constants.LoanTypeRehab: {
						LoanType: constants.LoanTypeRehab,
						GeneralRules: &RehabLoanTypeRules{
							Light: RehabLoanTypeRule{
								GeneralRule: &RuleCap{
									Ltaiv: &Values{
										Standard:  0.75,
										LargeUnit: 0.75,
									},
									Ltc: &Values{
										Standard:  0.8,
										LargeUnit: 0.8,
									},
									Ltarv: &Values{
										Standard:  0.7,
										LargeUnit: 0.7,
									},
								},
							},
						},
					},
				},
				HighValueStates:        zeebeeHighValueStates,
				BaseRate:               0.1,
				PreApprovalAmountRange: zeebeePreApprovalAmountRange,
				MinLargeUnitCount:      zeebeeMinLargeUnitCount,
				MaxUnitCount:           zeebeeMaxUnitCount,
				MinLoanTerm:            zeebeeMinLoanTerm,
				MaxLoanTerm:            zeebeeMaxLoanTerm,
				MinHeavyRehabPercent:   zeebeeMinHeavyRehabPercent,
				LtcBaseRateRules: []BaseRateRule{
					{
						Min:      0,
						Max:      0.71,
						BaseRate: 0.0975,
					},
					{
						Min:      0.71,
						Max:      0.76,
						BaseRate: 0.1025,
					},
					{
						Min:      0.76,
						Max:      0.81,
						BaseRate: 0.1075,
					},
				},
				LtarvBaseRateRules: []BaseRateRule{
					{
						Min:      0,
						Max:      0.61,
						BaseRate: 0.0975,
					},
					{
						Min:      0.61,
						Max:      0.66,
						BaseRate: 0.1025,
					},
					{
						Min:      0.66,
						Max:      0.71,
						BaseRate: 0.1075,
					},
				},
				AdjusterRule:     zeebeeAdjusterRule,
				ServiceFeeConfig: zeebeeServiceFeeConfig,
			},
			// Experienced
			ExperienceLevelExperienced: &PreApprovalRule{
				ExperienceLevel: ExperienceLevelExperienced,
				ExtraFactor:     zeebeeExtraFactor,
				PriceCapRuleMap: map[constants.LoanType]*LoanTypeRule{
					constants.LoanTypeBridge: {
						LoanType: constants.LoanTypeBridge,
						GeneralRules: &GeneralLoanTypeRules{
							GeneralRules: &RuleCap{
								Ltaiv: &Values{
									Standard:  0.75,
									LargeUnit: 0.75,
								},
								Ltc: &Values{
									Standard:  0.75,
									LargeUnit: 0.75,
								},
							},
						},
					},
					constants.LoanTypeConstruction: {
						LoanType: constants.LoanTypeConstruction,
						GeneralRules: &GeneralLoanTypeRules{
							GeneralRules: &RuleCap{
								Ltaiv: &Values{
									Standard:  0.85,
									LargeUnit: 0.7,
								},
								Ltc: &Values{
									Standard:  0.85,
									LargeUnit: 0.7,
								},
								Ltarv: &Values{
									Standard:  0.7,
									LargeUnit: 0.7,
								},
							},
						},
					},
					constants.LoanTypeRehab: {
						LoanType: constants.LoanTypeRehab,
						GeneralRules: &RehabLoanTypeRules{
							Light: RehabLoanTypeRule{
								GeneralRule: &RuleCap{
									Ltaiv: &Values{
										Standard:  0.85,
										LargeUnit: 0.85,
									},
									Ltc: &Values{
										Standard:  0.85,
										LargeUnit: 0.85,
									},
									Ltarv: &Values{
										Standard:  0.75,
										LargeUnit: 0.75,
									},
								},
							},
							Heavy: RehabLoanTypeRule{
								GeneralRule: &RuleCap{
									Ltaiv: &Values{
										Standard:  0.85,
										LargeUnit: 0.85,
									},
									Ltc: &Values{
										Standard:  0.85,
										LargeUnit: 0.85,
									},
									Ltarv: &Values{
										Standard:  0.7,
										LargeUnit: 0.7,
									},
								},
							},
						},
					},
				},
				HighValueStates:        zeebeeHighValueStates,
				BaseRate:               0.09,
				PreApprovalAmountRange: zeebeePreApprovalAmountRange,
				MinLargeUnitCount:      zeebeeMinLargeUnitCount,
				MaxUnitCount:           zeebeeMaxUnitCount,
				MinLoanTerm:            zeebeeMinLoanTerm,
				MaxLoanTerm:            zeebeeMaxLoanTerm,
				MinHeavyRehabPercent:   zeebeeMinHeavyRehabPercent,
				LtcBaseRateRules: []BaseRateRule{
					{
						Min:      0,
						Max:      0.71,
						BaseRate: 0.09,
					},
					{
						Min:      0.71,
						Max:      0.81,
						BaseRate: 0.0925,
					},
					{
						Min:      0.81,
						Max:      0.86,
						BaseRate: 0.095,
					},
				},
				LtarvBaseRateRules: []BaseRateRule{
					{
						Min:      0,
						Max:      0.61,
						BaseRate: 0.09,
					},
					{
						Min:      0.61,
						Max:      0.71,
						BaseRate: 0.0925,
					},
					{
						Min:      0.71,
						Max:      0.76,
						BaseRate: 0.095,
					},
				},
				AdjusterRule:     zeebeeAdjusterRule,
				ServiceFeeConfig: zeebeeServiceFeeConfig,
			},
			// institutional
			ExperienceLevelInstitutional: &PreApprovalRule{
				ExperienceLevel: ExperienceLevelInstitutional,
				ExtraFactor:     zeebeeExtraFactor,
				PriceCapRuleMap: map[constants.LoanType]*LoanTypeRule{
					constants.LoanTypeBridge: {
						LoanType: constants.LoanTypeBridge,
						GeneralRules: &GeneralLoanTypeRules{
							GeneralRules: &RuleCap{
								Ltaiv: &Values{
									Standard:  0.8,
									LargeUnit: 0.8,
								},
								Ltc: &Values{
									Standard:  0.8,
									LargeUnit: 0.8,
								},
							},
						},
						HighValueRules: &GeneralLoanTypeRules{
							GeneralRules: &RuleCap{
								Ltaiv: &Values{
									Standard:  0.8,
									LargeUnit: 0.8,
								},
								Ltc: &Values{
									Standard:  0.65,
									LargeUnit: 0.65,
								},
							},
						},
					},
					constants.LoanTypeConstruction: {
						LoanType: constants.LoanTypeConstruction,
						GeneralRules: &GeneralLoanTypeRules{
							GeneralRules: &RuleCap{
								Ltaiv: &Values{
									Standard:  0.85,
									LargeUnit: 0.7,
								},
								Ltc: &Values{
									Standard:  0.85,
									LargeUnit: 0.7,
								},
								Ltarv: &Values{
									Standard:  0.75,
									LargeUnit: 0.75,
								},
							},
						},
						HighValueRules: &GeneralLoanTypeRules{
							GeneralRules: &RuleCap{
								Ltaiv: &Values{
									Standard:  0.8,
									LargeUnit: 0.8,
								},
								Ltc: &Values{
									Standard:  0.65,
									LargeUnit: 0.65,
								},
							},
						},
					},
					constants.LoanTypeRehab: {
						LoanType: constants.LoanTypeRehab,
						GeneralRules: &RehabLoanTypeRules{
							Light: RehabLoanTypeRule{
								GeneralRule: &RuleCap{
									Ltaiv: &Values{
										Standard:  0.85,
										LargeUnit: 0.8,
									},
									Ltc: &Values{
										Standard:  0.9,
										LargeUnit: 0.8,
									},
									Ltarv: &Values{
										Standard:  0.75,
										LargeUnit: 0.75,
									},
								},
							},
							Heavy: RehabLoanTypeRule{
								GeneralRule: &RuleCap{
									Ltaiv: &Values{
										Standard:  0.85,
										LargeUnit: 0.8,
									},
									Ltc: &Values{
										Standard:  0.9,
										LargeUnit: 0.8,
									},
									Ltarv: &Values{
										Standard:  0.75,
										LargeUnit: 0.75,
									},
								},
							},
						},
						HighValueRules: &RehabLoanTypeRules{
							Light: RehabLoanTypeRule{
								GeneralRule: &RuleCap{
									Ltaiv: &Values{
										Standard:  0.85,
										LargeUnit: 0.8,
									},
									Ltc: &Values{
										Standard:  0.65,
										LargeUnit: 0.65,
									},
									Ltarv: &Values{
										Standard:  0.75,
										LargeUnit: 0.75,
									},
								},
							},
							Heavy: RehabLoanTypeRule{
								GeneralRule: &RuleCap{
									Ltaiv: &Values{
										Standard:  0.85,
										LargeUnit: 0.8,
									},
									Ltc: &Values{
										Standard:  0.8,
										LargeUnit: 0.8,
									},
									Ltarv: &Values{
										Standard:  0.65,
										LargeUnit: 0.65,
									},
								},
								HighValueRule: &RuleCap{
									Ltaiv: &Values{
										Standard:  0.85,
										LargeUnit: 0.8,
									},
									Ltc: &Values{
										Standard:  0.8,
										LargeUnit: 0.8,
									},
									Ltarv: &Values{
										Standard:  0.65,
										LargeUnit: 0.65,
									},
								},
							},
						},
						// GeneralRules: &GeneralLoanTypeRules{
						// 	GeneralRules: &RuleCap{
						// 		Ltaiv: &Values{
						// 			Standard:  0.85,
						// 			LargeUnit: 0.80,
						// 		},
						// 		Ltc: &Values{
						// 			Standard:  0.90,
						// 			LargeUnit: 0.8,
						// 		},
						// 		Ltarv: &Values{
						// 			Standard:  0.75,
						// 			LargeUnit: 0.75,
						// 		},
						// 	},
						// },
						// HighValueRules: &GeneralLoanTypeRules{
						// 	GeneralRules: &RuleCap{
						// 		Ltaiv: &Values{
						// 			Standard:  0.8,
						// 			LargeUnit: 0.8,
						// 		},
						// 		Ltc: &Values{
						// 			Standard:  0.65,
						// 			LargeUnit: 0.65,
						// 		},
						// 	},
						// },
					},
				},
				HighValueStates:        zeebeeHighValueStates,
				BaseRate:               0.0875,
				PreApprovalAmountRange: zeebeePreApprovalAmountRange,
				MinLargeUnitCount:      zeebeeMinLargeUnitCount,
				MaxUnitCount:           zeebeeMaxUnitCount,
				MinLoanTerm:            zeebeeMinLoanTerm,
				MaxLoanTerm:            zeebeeMaxLoanTerm,
				MinHeavyRehabPercent:   zeebeeMinHeavyRehabPercent,
				LtcBaseRateRules: []BaseRateRule{
					{
						Min:      0,
						Max:      0.71,
						BaseRate: 0.0825,
					},
					{
						Min:      0.71,
						Max:      0.81,
						BaseRate: 0.085,
					},
					{
						Min:      0.81,
						Max:      0.86,
						BaseRate: 0.0875,
					},
					{
						Min:      0.86,
						Max:      0.91,
						BaseRate: 0.09,
					},
				},
				LtarvBaseRateRules: []BaseRateRule{
					{
						Min:      0,
						Max:      0.61,
						BaseRate: 0.0825,
					},
					{
						Min:      0.61,
						Max:      0.71,
						BaseRate: 0.085,
					},
					{
						Min:      0.71,
						Max:      0.76,
						BaseRate: 0.0875,
					},
				},
				AdjusterRule:     zeebeeAdjusterRule,
				ServiceFeeConfig: zeebeeServiceFeeConfig,
			},
		},
	}

	return engines
}
