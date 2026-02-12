package types

type LoanBudget struct {
	OriginalBudget        float64 `json:"original_budget"`
	RemainingBudget       float64 `json:"remaining_budget"`
	CostSpentToDateAmount float64 `json:"cost_spent_to_date_amount"`
}
