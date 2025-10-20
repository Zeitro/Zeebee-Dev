package pre_approve_engine

import (
	"github.com/pkg/errors"

	"zeebee-backend/common/constants"
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
	UnitCount               int                   `json:"unit_count"`
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
	County                  string                `json:"county"`
	HasAssignmentFee        bool                  `json:"has_assignment_fee"`
	AssignmentFeeAmount     float64               `json:"assignment_fee_amount"`
	InterestReserveMonths   int                   `json:"interest_reserve_months"`
}

// CalculatePreApprovedLoanAmountResult represents the result of calculating pre-approved loan amount
type CalculatePreApprovedLoanAmountResult struct {
	IsApproved             bool            `json:"is_approved"`
	TotalAmount            float64         `json:"total_amount"`
	InterestReservedAmount float64         `json:"interest_reserved_amount"`
	Rate                   float64         `json:"rate"`
	MonthlyPayment         float64         `json:"monthly_payment"`
	CashToClose            float64         `json:"cash_to_close"`
	InitialAmount          float64         `json:"initial_amount"`
	HoldbackAmount         float64         `json:"holdback_amount"`
	LoanTerm               int             `json:"loan_term"`
	Ltv                    float32         `json:"ltv"`
	Tlc                    float32         `json:"tlc"`
	Ltarv                  float32         `json:"ltarv"`
	ExperienceLevel        ExperienceLevel `json:"experience_level"`
}

// PreApproveEngine represents the pre-approval engine
type PreApproveEngine struct {
	PreApproveRules map[ExperienceLevel]*PreApproveRule `json:"pre_approve_rules"`
	// ExperienceLevelRule defines rules for determining experience level
	ExperienceLevelRule ExperienceLevelRule `json:"experience_level_rule"`
}

// CalPreApprovedLoanAmount calculates and returns the pre-approved loan amount
func (e *PreApproveEngine) CalPreApprovedLoanAmount(input InputLoanInfo) (result *CalculatePreApprovedLoanAmountResult, err error) {
	experienceLevel := e.ExperienceLevelRule.getExperienceLevel(input.TransactionHistoryCount)
	rule, ok := e.PreApproveRules[experienceLevel]
	if !ok {
		return nil, errors.WithStack(errors.New("no rule found for experience level"))
	}
	return rule.CalPreApprovedLoanAmount(input)
}
