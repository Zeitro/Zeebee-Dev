package types

import "zeebee-backend/common/preapproval_engine"

type LoanBudget struct {
	OriginalBudget        float64 `json:"original_budget"`
	RemainingBudget       float64 `json:"remaining_budget"`
	CostSpentToDateAmount float64 `json:"cost_spent_to_date_amount"`
}

type LoanPreApprovalResult struct {
	RequestId              string                             `json:"request_id"`
	LenderCompanyName      string                             `json:"lender_company_name"`
	LenderCompanyId        int64                              `json:"lender_company_id"`
	IsApproved             bool                               `json:"is_approved"`
	RejectReason           string                             `json:"reject_reason"`
	TotalAmount            float64                            `json:"total_amount"`
	InterestReservedAmount float64                            `json:"interest_reserved_amount"`
	Rate                   float64                            `json:"rate"`
	MonthlyPayment         float64                            `json:"monthly_payment"`
	CashToClose            float64                            `json:"cash_to_close"`
	InitialAmount          float64                            `json:"initial_amount"`
	HoldbackAmount         float64                            `json:"holdback_amount"`
	LoanTerm               int                                `json:"loan_term"`
	Ltv                    float32                            `json:"ltv"`
	Ltc                    float32                            `json:"ltc"`
	Ltarv                  float32                            `json:"ltarv"`
	ExperienceLevel        preapproval_engine.ExperienceLevel `json:"experience_level"`
}

type LoanPreApprovalRequest struct {
	Borrower struct {
		Name                        string `json:"name,omitempty"`
		Email                       string `json:"email,omitempty"`
		Phone                       string `json:"phone,omitempty"`
		FicoScore                   int32  `json:"fico_score,omitempty"`
		FlipPropertyCount           int32  `json:"flip_property_count,omitempty"`
		PurchaseRentalPropertyCount int32  `json:"purchase_rental_property_count,omitempty"`
	}
	Property struct {
		Address struct {
			State    string `json:"state,omitempty"`
			City     string `json:"city,omitempty"`
			Address1 string `json:"address1,omitempty"`
			Address2 string `json:"address2,omitempty"`
			ZipCode  string `json:"zip_code,omitempty"`
		} `json:"address,omitempty"`
		PropertyType  string `json:"property_type,omitempty"`
		IsPermitReady bool   `json:"is_permit_ready,omitempty"`
		Unit          int    `json:"unit,omitempty"`
	}
	Loan struct {
		LoanTerm               int32   `json:"loan_term,omitempty"`
		EstPropertyValueAmount float64 ` json:"est_current_property_value_amount,omitempty"`
		Arv                    float64 `json:"arv,omitempty"`
		TotalAmount            float64 `json:"total_amount,omitempty"`
		PrincipalAmount        float64 `json:"principal_amount,omitempty"`
		PurchaseDate           string  `json:"purchase_date,omitempty"`
		PurchaseAmount         float64 `json:"purchase_amount,omitempty"`
		CostSpentToDateAmount  float64 `json:"cost_spent_to_date_amount,omitempty"`
		BudgetAmount           float64 `son:"budget_amount,omitempty"`
		InitialAmount          float64 `json:"initial_amount,omitempty"`
		HoldbackAmount         float64 `json:"holdback_amount,omitempty"`
		InterestReserveMonths  int32   `json:"interest_reserve_months,omitempty"`
		LoanPurpose            string  `json:"loan_purpose,omitempty"`
		LoanType               string  `json:"loan_type,omitempty"`
		HasAssignmentFee       bool    `json:"has_assignment_fee,omitempty"`
		AssignmentFeeAmount    float64 `json:"assignment_fee_amount,omitempty"`
	}
}

type PreApproveLoanCacheData struct {
	RequestData LoanPreApprovalRequest
	ResultData  *LoanPreApprovalResult
}
