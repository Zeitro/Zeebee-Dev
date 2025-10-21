package pre_approve_engine

import (
	"math"

	"github.com/duke-git/lancet/v2/condition"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/pkg/errors"

	"zeebee-backend/common/constants"
)

// HighValueState represents a high value state info
type HighValueState struct {
	State  string   `json:"state"`
	Cities []string `json:"cities"`
}

// LoanTypeRule represents the LTAIV/LTC/LTARV value rules for a specific loan type
type LoanTypeRule struct {
	LoanType       constants.LoanType `json:"loan_type"`        // loan type
	GeneralRules   Rules              `json:"general_rules"`    // general rules
	HighValueRules Rules              `json:"high_value_rules"` // high value rules
}

// PreApproveAmountRange represents the pre-approved loan amount range
type PreApproveAmountRange struct {
	StandardAmountRange  AmountRange `json:"standard_amount_range"`
	LargeUnitAmountRange AmountRange `json:"large_unit_amount_range"`
}

// AmountRange represents the amount range
type AmountRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// BaseRateRule represents a rule for determining base rates based on ranges
type BaseRateRule struct {
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
	BaseRate float64 `json:"base_rate"`
}

// AdjusterRule represents rules for adjusting rates based on various factors
type AdjusterRule struct {
	FicoScore                []RangeNumberRule           `json:"fico_score"`
	GUCExperienceLevel       map[ExperienceLevel]float64 `json:"guc_experience_level"`
	LtcAndLoanAmountRules    []LtcAndLoanAmountRule      `json:"ltc_and_loan_amount_rules"`
	LtcAndLoanPurposeRules   []LtcAndLoanPurposeRule     `json:"ltc_and_loan_purpose_rules"`
	LtarvAndLoanAmountRules  []LtarvAndLoanAmountRule    `json:"ltarv_and_loan_amount_rules"`
	LtarvAndLoanPurposeRules []LtarvAndLoanPurposeRule   `json:"ltarv_and_loan_purpose_rules"`
}

// RangeNumberRule represents a rule with min/max range and value
type RangeNumberRule struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Value float64 `json:"value"`
}

// RangeNumber represents a simple range with min and max values
type RangeNumber struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// LtcAndLoanAmountRule represents a rule based on LTC range and loan amount range
type LtcAndLoanAmountRule struct {
	LtcRange        RangeNumber `json:"ltc_range"`
	LoanAmountRange RangeNumber `json:"loan_amount_range"`
	Value           float64     `json:"value"`
}

// LtcAndLoanPurposeRule represents a rule based on LTC range and loan purpose
type LtcAndLoanPurposeRule struct {
	LtcRange    RangeNumber           `json:"ltc_range"`
	LoanPurpose constants.LoanPurpose `json:"loan_purpose"`
	Value       float64               `json:"value"`
}

// LtarvAndLoanAmountRule represents a rule based on LTARV range and loan amount range
type LtarvAndLoanAmountRule struct {
	LtarvRange      RangeNumber `json:"ltarv_range"`
	LoanAmountRange RangeNumber `json:"loan_amount_range"`
	Value           float64     `json:"value"`
}

// LtarvAndLoanPurposeRule represents a rule based on LTARV range and loan purpose
type LtarvAndLoanPurposeRule struct {
	LtarvRange  RangeNumber           `json:"ltarv_range"`
	LoanPurpose constants.LoanPurpose `json:"loan_purpose"`
	Value       float64               `json:"value"`
}

// ServiceFeeConfig represents configuration for service fees
type ServiceFeeConfig struct {
	OriginationFeePercent  float64 `json:"origination_fee_percent"`
	UnderwritingFeePercent float64 `json:"underwriting_fee_percent"`
}

// PreApproveRule represents the pre-approve rule
type PreApproveRule struct {
	ExperienceLevel ExperienceLevel `json:"experience_level"`
	// ExtraFactor is an extra factor for calculating ltc/ltaiv/ltarv
	ExtraFactor float64 `json:"extra_factor"`
	// PriceCapRuleMap is a map of loan type rules
	PriceCapRuleMap map[constants.LoanType]*LoanTypeRule `json:"price_cap_rule_map"`
	// HighValueStates represents high value states
	HighValueStates []*HighValueState `json:"high_value_states"`
	// BaseRate represents the base rate
	BaseRate float64 `json:"base_rate"`
	// PreApproveAmountRange represents pre-approved loan amount range
	PreApproveAmountRange PreApproveAmountRange `json:"pre_approve_amount_range"`
	// MinLargeUnitCount is the minimum unit count to be considered large unit
	MinLargeUnitCount int `json:"min_large_unit_count"`
	// MaxUnitCount is the maximum allowed unit count
	MaxUnitCount int `json:"max_unit_count"`
	// MinLoanTerm is the minimum loan term
	MinLoanTerm int `json:"min_loan_term"`
	// MaxLoanTerm is the maximum loan term
	MaxLoanTerm int `json:"max_loan_term"`
	// MinHeavyRehabPercent is the minimum percent to be considered heavy rehab
	MinHeavyRehabPercent float64 `json:"min_heavy_rehab_percent"`

	// LtcBaseRateRules are base rate rules for LTC
	LtcBaseRateRules []BaseRateRule `json:"ltc_base_rate_rules"`
	// LtarvBaseRateRules are base rate rules for LTARV
	LtarvBaseRateRules []BaseRateRule `json:"ltarv_base_rate_rules"`
	// AdjusterRule defines rate adjustment rules
	AdjusterRule AdjusterRule `json:"adjuster_rule"`
	// ServiceFeeConfig defines service fee configuration
	ServiceFeeConfig ServiceFeeConfig `json:"service_fee_config"`
}

// CalPreApprovedLoanAmount calculates and returns the pre-approved loan amount
func (r *PreApproveRule) CalPreApprovedLoanAmount(input InputLoanInfo) (result *CalculatePreApprovedLoanAmountResult, err error) {
	result = &CalculatePreApprovedLoanAmountResult{
		IsApproved: false,
	}

	// Get experience level
	// result.ExperienceLevel = r.ExperienceLevelRule.getExperienceLevel(input.TransactionHistoryCount)
	result.ExperienceLevel = r.ExperienceLevel

	// Check if loan is rejected
	if r.isLoanRejected(input) {
		return result, nil
	}

	// Calculate LTAIV, LTC, LTARV
	ltaiv, ltc, ltarv, err := r.getTotalLoanCap(input)
	if err != nil {
		return result, err
	}
	// Check if loan is rejected based on calculated values
	if ltaiv == 0 || ltc == 0 || ltarv == 0 {
		return result, nil
	}

	// Calculate max loan amount based on different criteria
	maxByLTAIV := input.EstCurrentPropertyValue * float64(ltaiv)
	maxByLTC := (input.EstCurrentPropertyValue + input.OriginalBudgetAmount) * float64(ltc)
	maxByLTARV := input.EstAfterRepairValue * float64(ltarv)

	// Get final loan amount using the most conservative value
	totalLoanAmount := math.Min(math.Min(maxByLTAIV, maxByLTC), maxByLTARV)
	isLargeUnit := r.isLargeUnit(input.UnitCount)
	if isLargeUnit {
		totalLoanAmount = math.Min(totalLoanAmount, r.PreApproveAmountRange.LargeUnitAmountRange.Max)
		totalLoanAmount = math.Max(totalLoanAmount, r.PreApproveAmountRange.LargeUnitAmountRange.Min)
	} else {
		totalLoanAmount = math.Min(totalLoanAmount, r.PreApproveAmountRange.StandardAmountRange.Max)
		totalLoanAmount = math.Max(totalLoanAmount, r.PreApproveAmountRange.StandardAmountRange.Min)
	}

	// Calculate GUC Interest reserve amount for construction loans
	if input.LoanType == constants.LoanTypeConstruction {
		interestReserveAmount := r.calculateInterestReserveAmount(totalLoanAmount, result.ExperienceLevel, input.InterestReserveMonths)
		totalLoanAmount += interestReserveAmount
		result.InterestReservedAmount = interestReserveAmount
	}
	result.TotalAmount = totalLoanAmount
	result.LoanTerm = input.LoanTerm

	if input.EstCurrentPropertyValue <= 0 {
		return result, errors.WithStack(errors.New("Invalid Estimated Current Value"))
	}

	// Calculate output LTC and LTARV ratios
	result.Tlc = float32(totalLoanAmount / (input.EstCurrentPropertyValue + input.OriginalBudgetAmount))
	result.Ltarv = float32(totalLoanAmount / input.EstAfterRepairValue)

	// Calculate rate
	baseRate := r.calculateBaseRate(ltc, ltarv)
	adjuster := r.calculateAdjusters(input, result.TotalAmount)
	result.Rate = baseRate + adjuster

	// Calculate monthly payment
	result.MonthlyPayment = r.calculateMonthlyPayment(result.TotalAmount, result.Rate, input.LoanTerm)

	// Calculate cash to close
	result.CashToClose = r.calculateCashToClose(input, result)

	// Calculate initial amount and holdback amount
	result.InitialAmount, result.HoldbackAmount = r.calculateInitialAndHoldbackAmount(input, result)

	result.IsApproved = true
	return result, nil
}

// isLoanRejected checks if a loan should be rejected based on various criteria
func (r *PreApproveRule) isLoanRejected(input InputLoanInfo) bool {
	// Check unit count - reject if exceeds maximum
	if input.UnitCount > r.MaxUnitCount {
		return true
	}

	// Budget amount check - reject if budget exceeds property value
	if input.OriginalBudgetAmount > input.EstCurrentPropertyValue {
		return true
	}

	// Loan term check - reject if outside allowed range
	if input.LoanTerm < r.MinLoanTerm || input.LoanTerm > r.MaxLoanTerm {
		return true
	}

	// Check borrower equity - reject if less than 10% of total cost
	borrowerEquity := r.calculateBorrowerEquity(input)
	totalCost := r.calculateTotalCost(input)
	if borrowerEquity < totalCost*0.1 {
		return true
	}

	// Loan type and purpose combinations that are not allowed
	if input.LoanType == constants.LoanTypeConstruction && input.LoanPurpose == constants.LoanPurposeCashOut {
		return true
	}

	if input.LoanType == constants.LoanTypeRehab && input.LoanPurpose == constants.LoanPurposeCashOut {
		return true
	}

	return false
}

// calculateBorrowerEquity calculates borrower equity based on loan purpose
func (r *PreApproveRule) calculateBorrowerEquity(input InputLoanInfo) (borrowerEquity float64) {
	switch input.LoanPurpose {
	case constants.LoanPurposePurchase:
		// For purchase loans, equity is purchase amount minus initial amount
		borrowerEquity = input.PurchaseAmount - input.InitialAmount
	case constants.LoanPurposeRateTerm, constants.LoanPurposeCashOut:
		// For rate term or cash out refinances, equity is property value minus total loan amount
		borrowerEquity = input.EstCurrentPropertyValue - input.TotalAmount
	default:
		// Default to purchase calculation
		borrowerEquity = input.PurchaseAmount - input.InitialAmount
	}
	return borrowerEquity
}

// calculateTotalCost calculates total cost amount based on loan type
func (r *PreApproveRule) calculateTotalCost(input InputLoanInfo) (totalCostAmount float64) {
	switch input.LoanType {
	case constants.LoanTypeRehab:
		// For rehab loans, total cost is property value plus budget amount
		totalCostAmount = input.EstCurrentPropertyValue + input.OriginalBudgetAmount
	case constants.LoanTypeConstruction, constants.LoanTypeBridge:
		// For construction and bridge loans, total cost is just property value
		totalCostAmount = input.EstCurrentPropertyValue
	default:
		// Default calculation includes property value and budget amount
		totalCostAmount = input.EstCurrentPropertyValue + input.OriginalBudgetAmount
	}
	return totalCostAmount
}

// getTotalLoanCap calculates the loan caps for LTAIV, LTC, and LTARV
func (r *PreApproveRule) getTotalLoanCap(input InputLoanInfo) (ltaiv float64, ltc float64, ltarv float64, err error) {
	// Determine various conditions
	isHighValueZone := r.isHighValueZone(input.State, input.County)
	isHeavyRehab := r.isHeavyRehab(input.LoanType, input.OriginalBudgetAmount, input.EstCurrentPropertyValue)
	isLargeUnit := r.isLargeUnit(input.UnitCount)

	// Get price cap rule for the loan type
	priceCapRule, ok := r.PriceCapRuleMap[input.LoanType]
	if !ok {
		return 0, 0, 0, errors.WithStack(errors.New("invalid loan type"))
	}

	var generalPriceCap *RuleCap
	var highValueZonePriceCap *RuleCap
	var priceRule Rules
	var rehabPriceRule RehabRules

	// Select rules based on high value zone
	if isHighValueZone && priceCapRule.HighValueRules != nil {
		priceRule = priceCapRule.HighValueRules
	} else {
		priceRule = priceCapRule.GeneralRules
	}

	// Handle rehab loan rules differently
	if input.LoanType == constants.LoanTypeRehab {
		rehabPriceRule, ok = priceRule.(RehabRules)
		if !ok {
			return 0, 0, 0, errors.WithStack(errors.New("invalid rehab rules"))
		}
		if isHeavyRehab {
			generalPriceCap, highValueZonePriceCap = rehabPriceRule.GetHeavyRules()
		} else {
			generalPriceCap, highValueZonePriceCap = rehabPriceRule.GetLightRules()
		}
	} else {
		generalPriceCap = priceRule.GetGeneralRules()
		highValueZonePriceCap = priceRule.GetHighValueRules()
	}

	// Calculate LTAIV based on conditions
	if isHighValueZone {
		if highValueZonePriceCap.Ltaiv == nil {
			ltaiv = 0
		} else {
			ltaiv = condition.Ternary(isLargeUnit, highValueZonePriceCap.Ltaiv.LargeUnit, highValueZonePriceCap.Ltaiv.Standard)
		}
	} else {
		if generalPriceCap.Ltaiv == nil {
			ltaiv = 0
		} else {
			ltaiv = condition.Ternary(isLargeUnit, generalPriceCap.Ltaiv.LargeUnit, generalPriceCap.Ltaiv.Standard)
		}
	}

	// Calculate LTC based on conditions
	if isHighValueZone {
		if highValueZonePriceCap.Ltc == nil {
			ltc = 0
		} else {
			ltc = condition.Ternary(isLargeUnit, highValueZonePriceCap.Ltc.LargeUnit, highValueZonePriceCap.Ltc.Standard)
		}
	} else {
		if generalPriceCap.Ltc == nil {
			ltc = 0
		} else {
			ltc = condition.Ternary(isLargeUnit, generalPriceCap.Ltc.LargeUnit, generalPriceCap.Ltc.Standard)
		}
	}

	// Calculate LTARV based on conditions
	if isHighValueZone {
		if highValueZonePriceCap.Ltarv == nil {
			ltarv = 0
		} else {
			ltarv = condition.Ternary(isLargeUnit, highValueZonePriceCap.Ltarv.LargeUnit, highValueZonePriceCap.Ltarv.Standard)
		}
	} else {
		if generalPriceCap.Ltarv == nil {
			ltarv = 0
		} else {
			ltarv = condition.Ternary(isLargeUnit, generalPriceCap.Ltarv.LargeUnit, generalPriceCap.Ltarv.Standard)
		}
	}

	// Apply extra factor adjustments
	if ltarv > 0 {
		ltarv += r.ExtraFactor
	}
	if ltc > 0 {
		ltc += r.ExtraFactor
	}
	if ltaiv > 0 {
		ltaiv += r.ExtraFactor
	}

	return ltaiv, ltc, ltarv, nil
}

// isHighValueZone checks if the given state and city are in the high-value zone
func (r *PreApproveRule) isHighValueZone(state string, city string) (isHighValueZone bool) {
	// Check if state is in high value states
	stateInfo, ok := slice.FindBy(r.HighValueStates, func(index int, value *HighValueState) bool {
		return state == value.State
	})
	if !ok {
		return false
	}
	// Check if city is in the high value state's cities
	_, isHighValueZone = slice.FindBy(stateInfo.Cities, func(index int, value string) bool {
		return city == value
	})

	return isHighValueZone
}

// isHeavyRehab checks if the rehabilitation is considered heavy based on the budget amount and current property value
func (r *PreApproveRule) isHeavyRehab(loanType constants.LoanType, budgetAmount float64, currentPropertyValue float64) (isHeavyRehab bool) {
	// Only rehab loans can be heavy rehab
	if loanType != constants.LoanTypeRehab {
		return false
	}
	// Compare the ratio of budget to property value with the minimum threshold
	return budgetAmount/currentPropertyValue > r.MinHeavyRehabPercent
}

// isLargeUnit returns true if the unit count qualifies as a large unit
func (r *PreApproveRule) isLargeUnit(unit int) (isLargeUnit bool) {
	return r.MinLargeUnitCount > 0 && unit >= r.MinLargeUnitCount
}

// calculateInterestReserveAmount calculates the interest reserve amount
func (r *PreApproveRule) calculateInterestReserveAmount(totalLoanAmount float64, experienceLevel ExperienceLevel, months int) float64 {
	defaultRate := r.getDefaultRate()
	interestReserveAmount := float64(months) * (0.7 * defaultRate * totalLoanAmount)
	return interestReserveAmount
}

// getDefaultRate returns the default rate
func (r *PreApproveRule) getDefaultRate() float64 {
	return r.BaseRate
}

// calculateBaseRate calculates and returns the base rate
func (r *PreApproveRule) calculateBaseRate(ltc float64, ltarv float64) (baseRate float64) {
	// Get base rates for LTC and LTARV
	ltcRate := r.getLtcBaseRate(ltc)
	ltarvRate := r.getLtarvBaseRate(ltarv)
	// Use the minimum of the two rates
	baseRate = math.Min(ltcRate, ltarvRate)
	return baseRate
}

// getLtcBaseRate returns the LTC base rate based on the LTC value
func (r *PreApproveRule) getLtcBaseRate(ltc float64) (baseRate float64) {
	// Find the matching base rate rule
	baseRateRule, ok := slice.FindBy(r.LtcBaseRateRules, func(index int, value BaseRateRule) bool {
		return ltc >= value.Min && ltc <= value.Max
	})
	if !ok {
		// If no rule matches, use the last rule's base rate or 0 if no rules exist
		l := len(r.LtcBaseRateRules)
		if l > 0 {
			baseRate = r.LtcBaseRateRules[l-1].BaseRate
		} else {
			baseRate = 0
		}
		return baseRate
	}
	baseRate = baseRateRule.BaseRate
	return baseRate
}

// getLtarvBaseRate returns the LTARV base rate based on the LTARV value
func (r *PreApproveRule) getLtarvBaseRate(ltarv float64) (baseRate float64) {
	// Find the matching base rate rule
	baseRateRule, ok := slice.FindBy(r.LtarvBaseRateRules, func(index int, value BaseRateRule) bool {
		return ltarv >= value.Min && ltarv <= value.Max
	})
	if !ok {
		// If no rule matches, use the last rule's base rate or 0 if no rules exist
		l := len(r.LtarvBaseRateRules)
		if l > 0 {
			baseRate = r.LtarvBaseRateRules[l-1].BaseRate
		} else {
			baseRate = 0
		}
		return baseRate
	}
	baseRate = baseRateRule.BaseRate
	return baseRate
}

// calculateAdjusters calculates rate adjusters based on various factors
func (r *PreApproveRule) calculateAdjusters(input InputLoanInfo, totalAmount float64) (adjusters float64) {
	rule := r.AdjusterRule

	// FICO score adjustments
	slice.ForEach(rule.FicoScore, func(index int, item RangeNumberRule) {
		if item.Max > float64(input.FicoScore) && item.Min <= float64(input.FicoScore) {
			adjusters += item.Value
		}
	})

	// GUC (Guaranteed Underwriting Certificate) adjustments for construction loans
	if input.LoanType == constants.LoanTypeConstruction {
		adjusters += rule.GUCExperienceLevel[r.ExperienceLevel]
	}

	// LTC/LTARV and loan amount adjustments
	var ltcAndLoanAmountAdjuster float64
	var ltarvAndLoanAmountAdjuster float64

	slice.ForEach(rule.LtcAndLoanAmountRules, func(index int, item LtcAndLoanAmountRule) {
		if float64(input.FicoScore) < item.LtcRange.Min || float64(input.FicoScore) > item.LtcRange.Max {
			return
		}
		if item.LoanAmountRange.Max > totalAmount && item.LoanAmountRange.Min <= totalAmount {
			ltcAndLoanAmountAdjuster = item.Value
		}
	})

	slice.ForEach(rule.LtarvAndLoanAmountRules, func(index int, item LtarvAndLoanAmountRule) {
		if float64(input.FicoScore) < item.LtarvRange.Min || float64(input.FicoScore) > item.LtarvRange.Max {
			return
		}
		if item.LoanAmountRange.Max > totalAmount && item.LoanAmountRange.Min <= totalAmount {
			ltarvAndLoanAmountAdjuster = item.Value
		}
	})
	adjusters += math.Min(ltcAndLoanAmountAdjuster, ltarvAndLoanAmountAdjuster)

	// LTC/LTARV and loan purpose adjustments
	var ltcAndLoanPurposeAdjuster float64
	var ltarvAndLoanPurposeAdjuster float64

	slice.ForEach(rule.LtcAndLoanPurposeRules, func(index int, item LtcAndLoanPurposeRule) {
		if float64(input.FicoScore) < item.LtcRange.Min || float64(input.FicoScore) > item.LtcRange.Max {
			return
		}
		if item.LoanPurpose == input.LoanPurpose {
			ltcAndLoanPurposeAdjuster = item.Value
		}
	})

	slice.ForEach(rule.LtarvAndLoanPurposeRules, func(index int, item LtarvAndLoanPurposeRule) {
		if float64(input.FicoScore) < item.LtarvRange.Min || float64(input.FicoScore) > item.LtarvRange.Max {
			return
		}
		if item.LoanPurpose == input.LoanPurpose {
			ltarvAndLoanPurposeAdjuster = item.Value
		}
	})
	adjusters += math.Min(ltcAndLoanPurposeAdjuster, ltarvAndLoanPurposeAdjuster)
	return adjusters
}

// calculateMonthlyPayment calculates the monthly payment
func (r *PreApproveRule) calculateMonthlyPayment(totalLoanAmount float64, rate float64, loanTerm int) float64 {
	if loanTerm <= 0 {
		return 0
	}
	monthlyRate := rate / 12
	return monthlyRate * totalLoanAmount / float64(loanTerm)
}

// calculateCashToClose calculates the cash to close (fees)
func (r *PreApproveRule) calculateCashToClose(input InputLoanInfo, result *CalculatePreApprovedLoanAmountResult) float64 {
	originationFee := r.ServiceFeeConfig.OriginationFeePercent * result.TotalAmount
	underwritingFee := r.ServiceFeeConfig.UnderwritingFeePercent * result.TotalAmount
	outOfPocket := r.calculateOutOfPocket(input, result)
	return originationFee + underwritingFee + outOfPocket
}

// calculateOutOfPocket calculates out of pocket expenses
func (r *PreApproveRule) calculateOutOfPocket(input InputLoanInfo, result *CalculatePreApprovedLoanAmountResult) float64 {
	switch input.LoanType {
	case constants.LoanTypeRehab:
		if input.LoanPurpose == constants.LoanPurposePurchase {
			assignmentFee := float64(0)
			if input.HasAssignmentFee {
				assignmentFee = input.AssignmentFeeAmount
			}
			return input.PurchaseAmount + assignmentFee + input.OriginalBudgetAmount - result.TotalAmount
		} else { // Refinance
			return input.PrincipalAmount + input.OriginalBudgetAmount - result.TotalAmount
		}
	case constants.LoanTypeConstruction:
		assignmentFee := float64(0)
		if input.HasAssignmentFee {
			assignmentFee = input.AssignmentFeeAmount
		}
		if input.LoanPurpose == constants.LoanPurposePurchase {
			return input.PurchaseAmount + assignmentFee + input.OriginalBudgetAmount + result.InterestReservedAmount - result.TotalAmount
		} else { // Refinance
			return input.PrincipalAmount + input.OriginalBudgetAmount + result.InterestReservedAmount - result.TotalAmount
		}
	case constants.LoanTypeBridge:
		assignmentFee := float64(0)
		if input.HasAssignmentFee {
			assignmentFee = input.AssignmentFeeAmount
		}
		if input.LoanPurpose == constants.LoanPurposePurchase {
			return input.PurchaseAmount + assignmentFee - result.TotalAmount
		} else {
			return input.PrincipalAmount - result.TotalAmount
		}
	default:
		return 0
	}
}

// calculateInitialAndHoldbackAmount calculates initial and holdback amounts
func (r *PreApproveRule) calculateInitialAndHoldbackAmount(input InputLoanInfo, result *CalculatePreApprovedLoanAmountResult) (initialAmount, holdbackAmount float64) {
	holdbackAmount = input.OriginalBudgetAmount

	// Initial loan = total loan amount - original budget - interest reserve amount
	interestReserveAmount := float64(0)
	if input.LoanType == constants.LoanTypeConstruction {
		interestReserveAmount = r.calculateInterestReserveAmount(result.TotalAmount, result.ExperienceLevel, input.InterestReserveMonths)
	}
	initialAmount = result.TotalAmount - input.OriginalBudgetAmount - interestReserveAmount

	return initialAmount, holdbackAmount
}
