package data

import "database/sql"

type Address struct {
	State    string `json:"state" copier:"State"`
	City     string `json:"city"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	ZipCode  string `json:"zip_code"`
}

type LoanBudgetEntity struct {
	OriginalBudget        float64 `json:"original_budget"`
	RemainingBudget       float64 `json:"remaining_budget"`
	CostSpentToDateAmount float64 `json:"cost_spent_to_date_amount"`
}

type LoanExtendValueEntity struct {
	Pricing             string               `json:"pricing"`
	LoanStatusUpdates   string               `json:"loan_status_updates"`
	MeetingNotes        string               `json:"meeting_notes"`
	ActionList          string               `json:"action_list"`
	BusinessLine        string               `json:"business_line"`
	LoanPreapprovalInfo *LoanPreapprovalInfo `json:"loan_preapproval_info,omitempty"`
}

type LoanPreapprovalInfo struct {
	Input  *LoanPreapprovalInfoInput  `json:"input,omitempty"`
	Result *LoanPreapprovalInfoResult `json:"result,omitempty"`
}

type LoanPreapprovalInfoInput struct {
	Borrower *LoanPreapprovalInfoInputBorrower `json:"borrower,omitempty"`
	Property *LoanPreapprovalInfoInputProperty `json:"property,omitempty"`
	Loan     *LoanPreapprovalInfoInputLoan     `json:"loan,omitempty"`
}

type LoanPreapprovalInfoResult struct {
	LenderOrganizationName string  `json:"lender_organization_name,omitempty"`
	LenderOrganizationId   int64   `json:"lender_organization_id,omitempty"`
	IsApproved             bool    `json:"is_approved,omitempty"`
	RejectReason           string  `json:"reject_reason,omitempty"`
	TotalAmount            float64 `json:"total_amount,omitempty"`
	InterestReservedAmount float64 `json:"interest_reserved_amount,omitempty"`
	Rate                   float32 `json:"rate,omitempty"`
	MonthlyPayment         float64 `json:"monthly_payment,omitempty"`
	CashToClose            float64 `json:"cash_to_close,omitempty"`
	InitialAmount          float64 `json:"initial_amount,omitempty"`
	HoldbackAmount         float64 `json:"holdback_amount,omitempty"`
	LoanTerm               int32   `json:"loan_term,omitempty"`
	Ltv                    float32 `json:"ltv,omitempty"`
	Ltc                    float32 `json:"ltc,omitempty"`
	Ltarv                  float32 `json:"ltarv,omitempty"`
	ExperienceLevel        string  `json:"experience_level,omitempty"`
}

type LoanPreapprovalInfoInputBorrower struct {
	Name                string `json:"name,omitempty"`
	Email               string `json:"email,omitempty"`
	Phone               string `json:"phone,omitempty"`
	EstCreditScore      int32  `json:"est_credit_score,omitempty"`
	FlipPropertyCount   int32  `json:"flip_property_count,omitempty"`
	PurchaseRentalCount int32  `json:"purchase_rental_property_count,omitempty"`
}

type LoanPreapprovalInfoInputProperty struct {
	Address       *Address `json:"address,omitempty"`
	PropertyType  string   `json:"property_type,omitempty"`
	IsPermitReady bool     `json:"is_permit_ready,omitempty"`
	Unit          int32    `json:"unit,omitempty"`
}

type LoanPreapprovalInfoInputLoan struct {
	LoanTerm               int32         `json:"loan_term,omitempty"`
	EstPropertyValueAmount float64       `json:"est_property_value_amount,omitempty"`
	Arv                    float64       `json:"arv,omitempty"`
	TotalAmount            float64       `json:"total_amount,omitempty"`
	PrincipalAmount        float64       `json:"principal_amount,omitempty"`
	PurchaseDate           *sql.NullTime `json:"purchase_date,omitempty"`
	PurchaseAmount         float64       `json:"purchase_amount,omitempty"`
	CostSpentToDateAmount  float64       `json:"cost_spent_to_date_amount,omitempty"`
	BudgetAmount           float64       `json:"budget_amount,omitempty"`
	InitialAmount          float64       `json:"initial_amount,omitempty"`
	HoldbackAmount         float64       `json:"holdback_amount,omitempty"`
	InterestReserveMonths  int32         `json:"interest_reserve_months,omitempty"`
	LoanPurpose            string        `json:"loan_purpose,omitempty"`
	LoanType               string        `json:"loan_type,omitempty"`
	HasAssignmentFee       bool          `json:"has_assignment_fee,omitempty"`
	AssignmentFeeAmount    float64       `json:"assignment_fee_amount,omitempty"`
}
